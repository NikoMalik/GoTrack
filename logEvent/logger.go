package logEvent

import (
	"io"
	"log"
	"os"

	kitlog "github.com/go-kit/log"
)

var logger kitlog.Logger

// Init initializes the logger with the provided environment variable name
func Init(envVar string) {
	var (
		logout io.Writer
		err    error
	)

	// Check if the environment variable is set
	logpath := os.Getenv(envVar)
	if logpath == "" {
		// If not set, default to standard error
		logout = os.Stderr
	} else {
		// If set, ensure the log file exists and open it
		if _, err := os.Stat(logpath); os.IsNotExist(err) {
			_, err := os.Create(logpath)
			if err != nil {
				log.Fatal(err)
			}
		}
		logout, err = os.OpenFile(logpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Initialize the kitlog logger
	logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(logout))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.Caller(4))
}

// Log logs a message with the initialized logger
func Log(args ...any) error {
	return logger.Log(args...)
}
