package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var logger *logrus.Logger

// Logger is a logger that supports log levels, context and structured logging.
type Logger interface {
	// Debug uses fmt.Sprint to construct and log a message at DEBUG level
	Debug(args ...interface{})
	// Info uses fmt.Sprint to construct and log a message at INFO level
	Info(args ...interface{})
	// Warning uses fmt.Sprint to construct and log a message at WARNING level
	Warning(args ...interface{})
	// Error uses fmt.Sprint to construct and log a message at ERROR level
	Error(args ...interface{})

	// Debugf uses fmt.Sprintf to construct and log a message at DEBUG level
	Debugf(format string, args ...interface{})
	// Infof uses fmt.Sprintf to construct and log a message at INFO level
	Infof(format string, args ...interface{})
	// Warningf uses fmt.Sprint to construct and log a message at WARNING level
	Warningf(format string, args ...interface{})
	// Errorf uses fmt.Sprintf to construct and log a message at ERROR level
	Errorf(format string, args ...interface{})
}

// New create new logger
func New() Logger {
	if logger != nil {
		return logger
	}

	logger = logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetReportCaller(true)

	// TODO configure logger for dev and production
	// Formatter := new(logrus.JSONFormatter)
	// Formatter.PrettyPrint = true

	Formatter := new(logrus.TextFormatter)
	Formatter.TimestampFormat = time.RFC3339
	logger.SetFormatter(Formatter)

	return logger
}
