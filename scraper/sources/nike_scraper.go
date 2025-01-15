package sources

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"

    "github.com/PuerkitoBio/goquery"
)

func ScrapeNike(db *sql.DB) error {
    url := "https://www.nike.com/w/new-3n82y"
    resp, err := http.Get(url)
    if err != nil {
        return fmt.Errorf("failed to get Nike page: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
    }

    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to parse Nike page: %w", err)
    }

    doc.Find(".product-card").Each(func(i int, s *goquery.Selection) {
        itemName := s.Find(".product-card__title").Text()
        price := s.Find(".product-price").Text()
        size := "Various" 
        onSale := false

        if err := insertClothes(db, "Nike", itemName, size, price, onSale); err != nil {
            log.Printf("Error inserting Nike item: %v", err)
        }
    })

    log.Println("Nike scraping complete.")
    return nil
}
