package postgres

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Declare a global variable for the connection pool
var pool *pgxpool.Pool

// Init - initialises the DB connection pool
func Init(connStr string) error {
	var err error

	// Create a connection pool
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Printf("Unable to parse connection string: %v\n", err)
		return err
	}
	config.MaxConns = 20
	pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return err
	}

	log.Println("Connected to the database")
	return nil
}

// CloseDB - closes the database connection
func CloseDB() {
	log.Println("Closing the database")
	if pool != nil {
		pool.Close()
	}
}

// FetchPrice - fetches the item price from database
func FetchPrice(name string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Prepare the query
	query := "SELECT price FROM items_table WHERE name = $1"

	// Execute the query and scan the result into the Item struct
	var price int
	err := pool.QueryRow(ctx, query, name).Scan(&price)
	if err != nil {
		if err == pgx.ErrNoRows {
			// No results found
			log.Println("Error no results found")
			return 0, err
		}
		return 0, err
	}

	return price, nil
}
