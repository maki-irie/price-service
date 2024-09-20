package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"go.uber.org/automaxprocs/maxprocs"

	"github.com/chrisbrown1111/price-service/pkg/downstream"
	"github.com/chrisbrown1111/price-service/pkg/postgres"
)

func newPriceHandler() (http.HandlerFunc, error) {
	remoteServer := ""
	if s := os.Getenv("REMOTE_SERVER"); s == "" {
		return nil, fmt.Errorf("env var REMOTE_SERVER not found")
	} else {
		remoteServer = s
	}
	dsClient := downstream.New(remoteServer)
	return func(w http.ResponseWriter, r *http.Request) {
		// We should always read the body and close it, even if we don't need it. If we don't
		// we'll leak goroutines and http.Requests....
		_, err := io.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
		}
		defer r.Body.Close()

		jwtToken := r.URL.Query().Get("jwt")

		claims, err := dsClient.ParseJWT(jwtToken)
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
			discount, err := dsClient.FetchDiscount(claims.Quantity)
			if err != nil {
				log.Printf("Fetch Discount %v", err)
				return
			}

			floatPrice := float32(price)
			totalPrice := (floatPrice - ((floatPrice * discount.Discount) / float32(100))) * float32(claims.Quantity)

			if claims.Vatincl {
				totalPrice = totalPrice * 1.2
			}

			// Set the content type to application/json
			w.Header().Set("Content-Type", "application/json")

			// Write the JSON response
			_, err = w.Write(
				[]byte(
					fmt.Sprintf(`{\"quality\":%d,\"tot_price\":%.2f }`, claims.Quantity, totalPrice),
				),
			)
			if err != nil {
				log.Printf("Error writing response: %s", err)
			}

		} else {
			log.Println("Price < 0!")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}, nil
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

	err = postgres.Init(connStr)
	if err != nil {
		return fmt.Errorf("error initializing database: %v", err)
	}
	defer postgres.CloseDB()

	mux := http.NewServeMux()
	handler, err := newPriceHandler()
	if err != nil {
		return fmt.Errorf("handler config error: %w", err)
	}
	mux.HandleFunc("/api/price", handler)

	s := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		MaxHeaderBytes:    http.DefaultMaxHeaderBytes,
	}

	log.Println("Starting the HTTP server on :8080")

	// Start the HTTP server
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("could not listen on 8080: %v", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
