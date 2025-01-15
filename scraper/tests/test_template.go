// tests/<scraper_name>_test.go
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

func TestScrape<Brand>(t *testing.T) {
    // Define mock HTML based on <Brand>'s website structure
    mockHTML := `
    <html>
        <body>
            <!-- Mock product entries -->
        </body>
    </html>
    `

    // Initialize mock HTTP server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, mockHTML)
    }))
    defer server.Close()

    // Initialize mock database
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    // Define expected database interactions (INSERT statements)
    mock.ExpectExec("INSERT INTO clothes").
        WithArgs("<Brand>", "<Item Name>", "Various", <Price>, <OnSale>).
        WillReturnResult(sqlmock.NewResult(1, 1))

    // Initialize logger
    log := logrus.New()
    log.Out = nil // Disable logging during tests

    // Initialize HTTP client
    client := server.Client()

    // Call the scraper function with mock server URL
    err = sources.Scrape<Brand>(db, log, client, server.URL)
    assert.NoError(t, err)

    // Ensure all expectations were met
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectations: %s", err)
    }
}
