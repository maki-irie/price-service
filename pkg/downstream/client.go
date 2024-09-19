package downstream

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Client struct {
	baseURL string
	hClient *http.Client
}

func New(baseURL string) *Client {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	return &Client{
		baseURL: baseURL,
		hClient: &http.Client{
			Transport: &http.Transport{
				Proxy:                 http.ProxyFromEnvironment,
				DialContext:           dialer.DialContext,
				ForceAttemptHTTP2:     true,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,

				// This defaults to 2 connections per host, which will have a massive impact given
				// we're running a load test on a single node. This could be much higher,
				// or we could even disable HTTP keepalives, which it looks like the Quarkus
				// behaviour.
				MaxIdleConnsPerHost: 100,
			},
		},
	}
}

type MyClaims struct {
	jwt.RegisteredClaims
	Item     string `json:"item"`
	Vatincl  bool   `json:"vat_incl"`
	Quantity int    `json:"quantity"`
}

type ApiResponse struct {
	// Define fields based on the expected JSON response
	Pkey string `json:"pkey"`
}

func (c *Client) getJwtKey() ([]byte, error) {
	// Define the URL of the API endpoint
	u := fmt.Sprintf("%s/pkey", c.baseURL)

	response, err := c.hClient.Get(u)
	if err != nil {
		return nil, fmt.Errorf("making GET request: %v", err)
	}
	defer response.Body.Close()

	// Check if the response status is OK
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status code %d", response.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	// Unmarshal JSON response into ApiResponse struct
	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling JSON: %v", err)
	}
	return []byte(apiResponse.Pkey), nil
}

func (c *Client) ParseJWT(tokenString string) (*MyClaims, error) {
	secretKey, err := c.getJwtKey()
	if err != nil {
		return nil, fmt.Errorf("priceHandler getJwtKey err: %v", err)
	}

	claims := &MyClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Since you use the same signing key to verify the token both in the
		// signing process and the verification process, you can use the same key.
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

// Define a struct that matches the JSON structure
type DiscountResponse struct {
	Discount           float32  `json:"discount"`
	Item               string   `json:"item"`
	Quantity           int      `json:"quantity"`
	Applicable_in_eu   bool     `json:"applicable_in_eu"`
	Shipping_cost      float32  `json:"shipping_cost"`
	Shipping_time_days int      `json:"shipping_time_days"`
	Related_items      []string `json:"related_items"`
}

func (c *Client) FetchDiscount(quantity int) (DiscountResponse, error) {
	u := fmt.Sprintf("%s/discount?quantity=%d", c.baseURL, quantity)

	// Make the GET request
	resp, err := c.hClient.Get(u)
	if err != nil {
		log.Printf("Error making GET request: %v", err)
		return DiscountResponse{}, err
	}
	defer resp.Body.Close()

	// Check if the response status is 200 OK
	if resp.StatusCode != http.StatusOK {
		return DiscountResponse{}, fmt.Errorf("status code is %d", resp.StatusCode)
	}

	// Decode the JSON response directly into the struct
	var discountResp DiscountResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&discountResp)
	if err != nil {
		return DiscountResponse{}, fmt.Errorf("decoding JSON: %v", err)
	}

	return discountResp, err
}
