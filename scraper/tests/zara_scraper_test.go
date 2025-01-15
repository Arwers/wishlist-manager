// tests/zara_scraper_test.go
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

func TestScrapeZara(t *testing.T) {
    // Mock HTTP Server
    mockHTML := `
    <html>
        <body>
            <div class="product-grid-product">
                <div class="product-name">Zara Shirt</div>
                <div class="price-amount">$29.99</div>
                <div class="sale-label">Sale</div>
            </div>
            <div class="product-grid-product">
                <div class="product-name">Zara Pants</div>
                <div class="price-amount">$49.99</div>
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

    // Expect INSERT for Zara Shirt
    mock.ExpectExec("INSERT INTO clothes").
        WithArgs("Zara", "Zara Shirt", "Various", 29.99, true).
        WillReturnResult(sqlmock.NewResult(1, 1))

    // Expect INSERT for Zara Pants
    mock.ExpectExec("INSERT INTO clothes").
        WithArgs("Zara", "Zara Pants", "Various", 49.99, false).
        WillReturnResult(sqlmock.NewResult(2, 1))

    // Logger
    log := logrus.New()
    log.Out = nil // Disable logging during tests

    // HTTP Client pointing to the mock server
    client := server.Client()

    // Replace the URL in ScrapeZara to point to the mock server
    // Assuming ScrapeZara accepts the URL as a parameter or is configurable
    // For this example, we'll temporarily modify the function to accept a URL

    // To avoid modifying the original function, we'll use URL Rewriting via Transport
    originalTransport := client.Transport
    client.Transport = rewriteURLTransport(server.URL, originalTransport)

    defer func() {
        client.Transport = originalTransport
    }()

    // Call ScrapeZara
    err = sources.ScrapeZara(db, log, client)
    assert.NoError(t, err)

    // Ensure all expectations were met
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectations: %s", err)
    }
}

// rewriteURLTransport rewrites the request URL to point to the mock server.
func rewriteURLTransport(mockURL string, originalTransport http.RoundTripper) http.RoundTripper {
    return http.HandlerFunc(func(req *http.Request) (*http.Response, error) {
        // Rewrite the request URL to the mock server
        req.URL.Scheme = "http"
        req.URL.Host = strings.TrimPrefix(mockURL, "http://")
        return originalTransport.RoundTrip(req)
    })
}

// Implement http.RoundTripper as a function
type http.HandlerFunc func(*http.Request) (*http.Response, error)

func (f http.HandlerFunc) RoundTrip(req *http.Request) (*http.Response, error) {
    return f(req)
}
