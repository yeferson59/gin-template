// Package logger provides structured logging utilities.
package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

// Init initializes the global logger with configuration.
func Init() *logrus.Logger {
	Log = logrus.New()

	// Set log level
	level := strings.ToUpper(os.Getenv("LOG_LEVEL"))
	switch level {
	case "TRACE":
		Log.SetLevel(logrus.TraceLevel)
	case "DEBUG":
		Log.SetLevel(logrus.DebugLevel)
	case "INFO":
		Log.SetLevel(logrus.InfoLevel)
	case "WARN", "WARNING":
		Log.SetLevel(logrus.WarnLevel)
	case "ERROR":
		Log.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		Log.SetLevel(logrus.FatalLevel)
	case "PANIC":
		Log.SetLevel(logrus.PanicLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}

	// Set log format
	logFormat := strings.ToUpper(os.Getenv("LOG_FORMAT"))
	if logFormat == "JSON" {
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// Set output
	Log.SetOutput(os.Stdout)

	return Log
}

// WithFields creates a new logger entry with the specified fields.
func WithFields(fields logrus.Fields) *logrus.Entry {
	if Log == nil {
		Init()
	}
	return Log.WithFields(fields)
}

// WithField creates a new logger entry with a single field.
func WithField(key string, value interface{}) *logrus.Entry {
	if Log == nil {
		Init()
	}
	return Log.WithField(key, value)
}

// Info logs an info message.
func Info(msg string) {
	if Log == nil {
		Init()
	}
	Log.Info(msg)
}

// Debug logs a debug message.
func Debug(msg string) {
	if Log == nil {
		Init()
	}
	Log.Debug(msg)
}

// Warn logs a warning message.
func Warn(msg string) {
	if Log == nil {
		Init()
	}
	Log.Warn(msg)
}

// Error logs an error message.
func Error(msg string) {
	if Log == nil {
		Init()
	}
	Log.Error(msg)
}

// Fatal logs a fatal message and exits.
func Fatal(msg string) {
	if Log == nil {
		Init()
	}
	Log.Fatal(msg)
}

// Panic logs a panic message and panics.
func Panic(msg string) {
	if Log == nil {
		Init()
	}
	Log.Panic(msg)
}
