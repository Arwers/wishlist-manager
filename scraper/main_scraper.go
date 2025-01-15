package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "log"
    "os"
    "example.com/my-scraper/sources"
)

func main() {
    connStr := os.Getenv("DB_CONN_STRING")
    if connStr == "" {
        connStr = "user=postgres password=secret dbname=clothes_scraper sslmode=disable"
    }

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Failed to open database connection: %v", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }

    fmt.Println("Connected to the database successfully!")

    if err := createTables(db); err != nil {
        log.Fatalf("Error creating tables: %v", err)
    }

    if err := sources.ScrapeAdidas(db); err != nil {
        log.Printf("Error scraping Adidas: %v", err)
    }
    if err := sources.ScrapeNike(db); err != nil {
        log.Printf("Error scraping Nike: %v", err)
    }
	
    fmt.Println("Scraping finished.")
}

func createTables(db *sql.DB) error {
    queries := []string{
        `CREATE TABLE IF NOT EXISTS clothes (
            id SERIAL PRIMARY KEY,
            brand VARCHAR(50),
            item_name VARCHAR(100),
            size VARCHAR(50),
            price NUMERIC(10, 2),
            on_sale BOOLEAN,
            scraped_date TIMESTAMP DEFAULT NOW()
         );`,
        `CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name VARCHAR(50),
            email VARCHAR(100)
         );`,
        `CREATE TABLE IF NOT EXISTS user_subscriptions (
            user_id INT REFERENCES users(id),
            clothes_id INT REFERENCES clothes(id),
            PRIMARY KEY (user_id, clothes_id)
         );`,
    }

    for _, q := range queries {
        if _, err := db.Exec(q); err != nil {
            return fmt.Errorf("failed to execute query (%s): %w", q, err)
        }
    }

    return nil
}
