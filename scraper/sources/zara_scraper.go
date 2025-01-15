package sources

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "strings"

    "github.com/PuerkitoBio/goquery"
)

func ScrapeZara(db *sql.DB) error {
    url := "https://www.zara.com/us/en/man-new-in-l711.html"
    resp, err := http.Get(url)
    if err != nil {
        return fmt.Errorf("failed to GET Zara page: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
    }

    // Parse HTML
    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to parse Zara page: %w", err)
    }

    doc.Find(".product-grid-product").Each(func(i int, s *goquery.Selection) {
        itemName := s.Find(".product-name").Text()
        priceText := s.Find(".price-amount").Text()
        saleLabel := s.Find(".sale-label").Text()
        onSale := strings.TrimSpace(saleLabel) != ""
        size := "Various"

        // Insert into DB
        if err := insertClothes(db, "Zara", itemName, size, priceText, onSale); err != nil {
            log.Printf("Error inserting Zara item: %v", err)
        }
    })

    log.Println("Zara scraping complete.")
    return nil
}

// insertClothes inserts data into the clothes table (same helper used by other scrapers).
func insertClothes(db *sql.DB, brand, itemName, size, priceStr string, onSale bool) error {
    query := `
        INSERT INTO clothes (brand, item_name, size, price, on_sale)
        VALUES ($1, $2, $3, $4, $5)
    `

    var price float64
    fmt.Sscanf(priceStr, "$%f", &price)

    _, err := db.Exec(query, brand, itemName, size, price, onSale)
    if err != nil {
        return fmt.Errorf("failed to insert clothes item: %w", err)
    }
    return nil
}
