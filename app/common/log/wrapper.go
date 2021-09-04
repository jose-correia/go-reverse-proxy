package log

import (
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
)

// NewLogger creates a Logger with UTC time
func NewLogger() log.Logger {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "timestamp", log.DefaultTimestampUTC)
	return logger
}

// NewNopLogger is go-kit/log.NewNopLogger
func NewNopLogger() log.Logger {
	return log.NewNopLogger()
}

// NewSyncLogger is usefull for debugging, use sparingly
func NewSyncLogger() log.Logger {
	return log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
}

func Str(i interface{}) string { return fmt.Sprintf("%+v", i) }
