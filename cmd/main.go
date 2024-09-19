package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"go.uber.org/automaxprocs/maxprocs"

	"github.com/chrisbrown1111/price-service/pkg/jwt"
	"github.com/chrisbrown1111/price-service/pkg/postgres"
)

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

func priceHandler(w http.ResponseWriter, r *http.Request) {

	jwtToken := r.URL.Query().Get("jwt")

	secretKey, err := jwt.GetJwtKey()
	if err != nil {
		log.Fatalf("priceHandler getJwtKey err: %v", err)
	}

	claims, err := jwt.ParseJWT(jwtToken, secretKey)
	if err != nil {
		log.Printf("Error parsing JWT: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Fetch the price from Postgres
	price, err := postgres.FetchPrice(claims.Item)
	if err != nil {
		log.Printf("fetchPrice %v", err)
		return
	}

	if price >= 0 {

		discount, err := fetchDiscount(claims.Quantity)
		if err != nil {
			log.Printf("Fetch Discount %v", err)
			return
		}

		var totalPrice float32
		var discountedPrice float32 = float32(price) - float32(price)*(discount.Discount)/float32(100)
		totalPrice = float32(discountedPrice * float32(claims.Quantity))

		if claims.Vatincl {
			totalPrice = totalPrice * 1.2
		}

		// Set the content type to application/json
		w.Header().Set("Content-Type", "application/json")
		// Write the JSON response
		fmt.Fprintf(w, "{\"quality\":%d,\"tot_price\":%.2f }", claims.Quantity, totalPrice)
		fmt.Println("Request done")
	} else {
		log.Println("Price < 0!")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func fetchDiscount(quantity int) (DiscountResponse, error) {
	// The URL of the HTTP server you want to call
	baseUrl := os.Getenv("REMOTE_SERVER")
	if len(baseUrl) == 0 {
		baseUrl = "http://localhost:7070"
	}
	baseUrl = baseUrl + "/discount"

	// Create a URL object
	u, err := url.Parse(baseUrl)
	if err != nil {
		log.Fatalf("Error parsing URL: %v\n", err)
	}

	query := u.Query()
	query.Set("quantity", strconv.Itoa(quantity))
	u.RawQuery = query.Encode()

	// Make the GET request
	resp, err := http.Get(u.String())
	if err != nil {
		log.Printf("Error making GET request: %v", err)
		return DiscountResponse{}, err
	}
	defer resp.Body.Close()

	// Check if the response status is 200 OK
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: status code is %d\n", resp.StatusCode)
	}

	// Decode the JSON response directly into the struct
	var discountResp DiscountResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&discountResp)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v\n", err)
	}

	// Print the unmarshalled data
	// log.Printf("Item: %s Discount: %.2f Quantity: %d\n", discountResp.Item, discountResp.Discount, discountResp.Quantity)
	return discountResp, err
}

func run() error {
	var err error
	log.SetOutput(os.Stdout)

	_, err = maxprocs.Set(maxprocs.Logger(log.Printf))
	if err != nil {
		return err
	}

	// Postgres connection string
	connStr := os.Getenv("DB_IP")
	if len(connStr) == 0 {
		connStr = "postgres://postgres:mysecretpassword@localhost:5432/test_db"
	}
	// give postgres the time to start when using docker compose
	time.Sleep(5 * time.Second)
	err = postgres.Init(connStr)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer postgres.CloseDB()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/price", priceHandler)

	log.Println("Starting the HTTP server on :8080")

	// Start the HTTP server
	if err := http.ListenAndServe("0.0.0.0:8080", mux); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on 8080: %v\n", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
