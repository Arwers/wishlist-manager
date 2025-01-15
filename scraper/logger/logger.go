package logger

import (
    "os"

    "github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
    Log = logrus.New()

    // Set the output to stdout
    Log.Out = os.Stdout

    // Set the log format to JSON for structured logging
    Log.SetFormatter(&logrus.JSONFormatter{})

    // Set log level (can be parameterized via environment variables)
    Log.SetLevel(logrus.InfoLevel)
}
