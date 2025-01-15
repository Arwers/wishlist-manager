package sources

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "regexp"
    "strings"

    "github.com/PuerkitoBio/goquery"
)

// ScrapeZalando scrapes Zalando's website for clothes data.
func ScrapeZalando(db *sql.DB) error {
    url := "https://www.zalando.com/men-home/" // Example mens page
    resp, err := http.Get(url)
    if err != nil {
        return fmt.Errorf("failed to GET Zalando page: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
    }

    // Parse HTML
    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to parse Zalando page: %w", err)
    }

    doc.Find(".z-grid-item").Each(func(i int, s *goquery.Selection) {
        itemName := s.Find(".z-article-card-title").Text()
        priceText := s.Find(".z-article-card-price").Text()

        itemName = strings.TrimSpace(itemName)
        priceText = strings.TrimSpace(priceText)

        saleText := s.Find(".z-article-card-sale").Text()
        onSale := strings.TrimSpace(saleText) != ""

        size := "Various"

        if err := insertClothes(db, "Zalando", itemName, size, priceText, onSale); err != nil {
            log.Printf("Error inserting Zalando item: %v", err)
        }
    })

    log.Println("Zalando scraping complete.")
    return nil
}

func insertClothes(db *sql.DB, brand, itemName, size, priceStr string, onSale bool) error {
    query := `
        INSERT INTO clothes (brand, item_name, size, price, on_sale)
        VALUES ($1, $2, $3, $4, $5)
    `
    var price float64
    priceStr = strings.ReplaceAll(priceStr, ",", ".")

    re := regexp.MustCompile(`[^\d\.]`)
    cleanPrice := re.ReplaceAllString(priceStr, "")
    fmt.Sscanf(cleanPrice, "%f", &price)

    _, err := db.Exec(query, brand, itemName, size, price, onSale)
    if err != nil {
        return fmt.Errorf("failed to insert clothes item: %w", err)
    }
    return nil
}
