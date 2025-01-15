package main

import (
    "database/sql"
    "fmt"
    "os"
    "time"

    "github.com/robfig/cron/v3"
    _ "github.com/lib/pq"

    "example.com/my-scraper/logger"
    "example.com/my-scraper/sources"
)

func main() {
    // Initialize Logger
    logger.InitLogger()
    log := logger.Log

    // Connect to DB
    connStr := os.Getenv("DB_CONN_STRING")
    if connStr == "" {
        connStr = "user=postgres password=secret dbname=clothes_scraper sslmode=disable"
    }

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Failed to open database connection: %v", err)
    }
    defer db.Close()

    // Ensure connection is valid
    if err := db.Ping(); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }
    log.Info("Connected to the database successfully!")

    // Create or verify tables (if you don't use migrations).
    if err := createTables(db, log); err != nil {
        log.Fatalf("Error creating tables: %v", err)
    }

    // Initialize cron scheduler
    c := cron.New(cron.WithSeconds()) // with seconds precision

    // Schedule job to run every hour, on the hour.
    _, err = c.AddFunc("0 0 * * * *", func() {
        log.WithFields(logrus.Fields{
            "time": time.Now().Format(time.RFC1123),
        }).Info("Running scheduled scraping job")
        sources.ScrapeAll(db, log) // pass the logger
    })
    if err != nil {
        log.Fatalf("Failed to add cron function: %v", err)
    }

    // Start scheduler (non-blocking in the background)
    c.Start()
    log.Info("Scheduler started. Scraping will run every hour.")

    // Keep the main function alive.
    select {}
}

func createTables(db *sql.DB, log *logrus.Logger) error {
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
            log.WithFields(logrus.Fields{
                "query": q,
                "error": err,
            }).Error("Failed to execute query")
            return fmt.Errorf("failed to execute query (%s): %w", q, err)
        }
    }

    log.Info("Database tables are set up successfully.")
    return nil
}
