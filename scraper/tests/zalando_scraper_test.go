// tests/zalando_scraper_test.go
package tests

import (
    "database/sql"
    "fmt"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/sirupsen/logrus"
    "github.com/stretchr/testify/assert"

    "example.com/my-scraper/sources"
)

func TestScrapeZalando(t *testing.T) {
    // Mock HTTP Server
    mockHTML := `
    <html>
        <body>
            <div class="z-grid-item">
                <div class="z-article-card-title">Zalando Jacket</div>
                <div class="z-article-card-price">€79,99</div>
                <div class="z-article-card-sale">Sale</div>
            </div>
            <div class="z-grid-item">
                <div class="z-article-card-title">Zalando Sneakers</div>
                <div class="z-article-card-price">€59,99</div>
            </div>
        </body>
    </html>
    `
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, mockHTML)
    }))
    defer server.Close()

    // Mock Database
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    // Expect INSERT for Zalando Jacket
    mock.ExpectExec("INSERT INTO clothes").
        WithArgs("Zalando", "Zalando Jacket", "Various", 79.99, true).
        WillReturnResult(sqlmock.NewResult(1, 1))

    // Expect INSERT for Zalando Sneakers
    mock.ExpectExec("INSERT INTO clothes").
        WithArgs("Zalando", "Zalando Sneakers", "Various", 59.99, false).
        WillReturnResult(sqlmock.NewResult(2, 1))

    // Logger
    log := logrus.New()
    log.Out = nil // Disable logging during tests

    // HTTP Client pointing to the mock server
    client := server.Client()

    // Call ScrapeZalando with mock server URL
    err = sources.ScrapeZalando(db, log, client, server.URL)
    assert.NoError(t, err)

    // Ensure all expectations were met
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectations: %s", err)
    }
}
