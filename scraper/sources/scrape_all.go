package sources

import (
    "database/sql"

    "github.com/sirupsen/logrus"
)

// ScrapeAll triggers all known scrapers one after another.
func ScrapeAll(db *sql.DB, log *logrus.Logger) {
    if err := ScrapeAdidas(db, log); err != nil {
        log.WithFields(logrus.Fields{
            "error": err,
        }).Error("Error scraping Adidas")
    }
    if err := ScrapeNike(db, log); err != nil {
        log.WithFields(logrus.Fields{
            "error": err,
        }).Error("Error scraping Nike")
    }
    if err := ScrapeZara(db, log); err != nil {
        log.WithFields(logrus.Fields{
            "error": err,
        }).Error("Error scraping Zara")
    }
    if err := ScrapeZalando(db, log); err != nil {
        log.WithFields(logrus.Fields{
            "error": err,
        }).Error("Error scraping Zalando")
    }
}
