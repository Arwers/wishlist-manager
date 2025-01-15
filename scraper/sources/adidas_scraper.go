package sources

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"

    "github.com/PuerkitoBio/goquery"
)

func ScrapeAdidas(db *sql.DB) error {
    url := "https://www.adidas.com/us/men-new_arrivals"
    resp, err := http.Get(url)
    if err != nil {
        return fmt.Errorf("failed to get Adidas page: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
    }

    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to parse Adidas page: %w", err)
    }

    doc.Find(".product-card").Each(func(i int, s *goquery.Selection) {
        itemName := s.Find(".product-card__title").Text()
        price := s.Find(".product-price").Text()
        size := "Various"
        onSale := false

        // Insert into DB
        if err := insertClothes(db, "Adidas", itemName, size, price, onSale); err != nil {
            log.Printf("Error inserting Adidas item: %v", err)
        }
    })

    log.Println("Adidas scraping complete.")
    return nil
}

func insertClothes(db *sql.DB, brand, itemName, size, priceStr string, onSale bool) error {
    query := `
        INSERT INTO clothes (brand, item_name, size, price, on_sale)
        VALUES ($1, $2, $3, $4, $5)
    `
    // Convert price string to float
    var price float64
    fmt.Sscanf(priceStr, "$%f", &price) // e.g. from "$59.99" to 59.99

    _, err := db.Exec(query, brand, itemName, size, price, onSale)
    if err != nil {
        return fmt.Errorf("failed to insert clothes item: %w", err)
    }
    return nil
}
