package jwt

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

type MyClaims struct {
	jwt.StandardClaims
	Item     string `json:"item"`
	Vatincl  bool   `json:"vat_incl"`
	Quantity int    `json:"quantity"`
}

type ApiResponse struct {
	// Define fields based on the expected JSON response
	Pkey string `json:"pkey"`
}

func GetJwtKey() ([]byte, error) {
	// Define the URL of the API endpoint
	url := os.Getenv("REMOTE_SERVER")
	if len(url) == 0 {
		url = "http://localhost:7070"
	}
	url = url + "/pkey"

	// Perform the GET request
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error making GET request: %v", err)
	}
	defer response.Body.Close()

	// Check if the response status is OK
	if response.StatusCode != http.StatusOK {
		log.Fatalf("Error: Received non-OK status code %d", response.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Unmarshal JSON response into ApiResponse struct
	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}
	return []byte(apiResponse.Pkey), err
}

func ParseJWT(tokenString string, secretKey []byte) (*MyClaims, error) {
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
