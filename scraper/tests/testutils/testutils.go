// tests/testutils/testutils.go
package testutils

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/sirupsen/logrus"
)

// SetupMockServer initializes a mock HTTP server with given response.
func SetupMockServer(t *testing.T, response string) *httptest.Server {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, response)
    }))
    t.Cleanup(func() {
        server.Close()
    })
    return server
}

// SetupMockDB initializes a mock database and returns the db and sqlmock instances.
func SetupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
    }
    t.Cleanup(func() {
        db.Close()
    })
    return db, mock
}

// SetupLogger initializes a logger with output disabled.
func SetupLogger() *logrus.Logger {
    log := logrus.New()
    log.Out = nil
    return log
}
