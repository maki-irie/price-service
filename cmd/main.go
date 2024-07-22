package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/chrisbrown1111/price-service/pkg/jwt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Define a struct that matches the JSON structure
type DiscountResponse struct {
	Discount           float32  `json:"discount"`
	Item               string   `json:"item"`
	Quantity           int      `json:"quantity"`
	Applicable_in_eu   bool     `json:"applicable_in_eu"`
	Shipping_cost      float32  `json:"shipping_cost"`
	Shipping_time_days int      `json:"shipping_time_days:`
	Related_items      []string `json:"related_items"`
}

// Declare a global variable for the connection pool
var pool *pgxpool.Pool

func priceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("priceHandler - starts")

	jwtToken := r.URL.Query().Get("jwt")

	//secretKey := []byte("UuMGMJMcRInbynNiyGX3zoz0YipAQyLmvn5efGo3wo4i9hf335Xh2TEbRpArlVxhRRly5G3Cgi1mLmQtVyogWuqy7xJfgB47iO2nXZbco0KDXX6SDDjSrpaWFFjmln7a")
	secretKey, err := jwt.GetJwtKey()
	if err != nil {
		log.Fatal("priceHandler getJwtKey err: ", err)
	}

	//jwtToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3MjEzMDk5OTIsImV4cCI6MTc1Mjg0NTk5MiwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIml0ZW0iOiJhcHBsZXMiLCJ2YXQtaW5jbCI6InRydWUiLCJxdWFudGl0eSI6IjEyNCJ9.RutoxPIxpGFUf62bU4rjY7Haq08Lydft1q0Vv58bOrU"
	fmt.Println("jwt: ", jwtToken)
	claims, err := jwt.ParseJWT(jwtToken, secretKey)
	if err != nil {
		fmt.Println("Error parsing JWT:", err)
		return
	}

	fmt.Printf("Claims: %#v\n", claims)

	// Fetch the price from Postgres
	price, err := fetchPrice(claims.Item)
	if err != nil {
		log.Fatal("fetchPrice", err)
	}

	if price != 0 {
		fmt.Println("Fetched price:", price)
		discount, err := fetchDiscount(claims.Quantity)
		if err != nil {
			log.Fatal("Fetch Discount", err)
		}

		var totalPrice float32
		quantity, err := strconv.Atoi(claims.Quantity)
		if err != nil {
			fmt.Printf("Error converting '%s' to bool: %v\n", claims.Quantity, err)
		}
		totalPrice = float32(price * quantity)
		VatInclBool, err := strconv.ParseBool(claims.Vatincl)
		if err != nil {
			fmt.Printf("Error converting '%s' to bool: %v\n", claims.Vatincl, err)
		}

		if VatInclBool {
			totalPrice = totalPrice * 1.2 * float32(quantity)
		}
		if discount.Discount > 0.0 {
			totalPrice = totalPrice - totalPrice*discount.Discount/100
		}

		fmt.Printf("totalPrice: %.2f\n", totalPrice)

		// Set the content type to application/json
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"quality\":%s,\"price\":%.2f }", claims.Quantity, totalPrice)
		// Write the JSON response
		//w.WriteHeader(http.StatusOK)
	} else {
		fmt.Println("Item not found")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func fetchPrice(name string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Prepare the query
	query := "SELECT price FROM item_table WHERE name = $1"

	// Execute the query and scan the result into the Item struct
	var price int
	err := pool.QueryRow(ctx, query, name).Scan(&price)
	if err != nil {
		if err == pgx.ErrNoRows {
			// No result found
			return 0, nil
		}
		return 0, err
	}

	return price, nil
}

func fetchDiscount(quantity string) (DiscountResponse, error) {
	fmt.Println("fetchDiscount() - quantity:", quantity)
	// The URL of the HTTP server you want to call
	baseUrl := "http://localhost:7070/discount"

	// Create a URL object
	u, err := url.Parse(baseUrl)
	if err != nil {
		log.Fatalf("Error parsing URL: %v\n", err)
	}

	query := u.Query()
	query.Set("quantity", quantity)
	u.RawQuery = query.Encode()

	// Make the GET request
	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatalf("Error making GET request: %v\n", err)
	}
	defer resp.Body.Close()

	// Check if the response status is 200 OK
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: status code is %d\n", resp.StatusCode)
	}

	// Decode the JSON response directly into the struct
	var discountResponse DiscountResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&discountResponse)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v\n", err)
	}

	// Print the unmarshalled data
	fmt.Printf("Discount: %.2f\n", discountResponse.Discount)
	fmt.Printf("Item: %s\n", discountResponse.Item)
	fmt.Printf("Quantity: %d\n", discountResponse.Quantity)

	return discountResponse, err
}

func main() {
	var err error
	// Connection string
	connStr := "postgres://postgres:mysecretpassword@localhost:5432/postgres"

	// Create a connection pool
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v\n", err)
	}
	config.MaxConns = 10
	pool, err = pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	fmt.Println("Connected to the database!")

	mux := http.NewServeMux()
	mux.HandleFunc("/price", priceHandler)

	fmt.Println("Starting the HTTP server on :8081")

	// Start the HTTP server
	if err := http.ListenAndServe("localhost:8081", mux); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Could not listen on 8081: %v\n", err)
	}

}
