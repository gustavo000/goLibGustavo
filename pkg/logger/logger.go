package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Global logger instance
	logger zerolog.Logger
)

// Initialize the logger with configuration
func init() {
	// Set time format to RFC3339
	zerolog.TimeFieldFormat = zerolog.TimeFormatRFC3339
	
	// Configure console writer with pretty formatting for development
	consoleWriter := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = "2006-01-02 15:04:05"
	})
	
	// Create logger with console output
	logger = zerolog.New(consoleWriter).
		With().
		Timestamp().
		Caller().
		Logger()
	
	// Set global log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

// SetLogLevel sets the global log level
func SetLogLevel(level string) {
	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

// SetProductionMode configures logger for production (JSON output)
func SetProductionMode() {
	logger = zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Logger()
}

// Info logs an info message
func Info(msg string, fields ...map[string]interface{}) {
	event := logger.Info()
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Infof logs a formatted info message
func Infof(format string, v ...interface{}) {
	logger.Info().Msgf(format, v...)
}

// Warning logs a warning message
func Warning(msg string, fields ...map[string]interface{}) {
	event := logger.Warn()
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Warningf logs a formatted warning message
func Warningf(format string, v ...interface{}) {
	logger.Warn().Msgf(format, v...)
}

// Danger logs an error message (renamed from Error to Danger to match original API)
func Danger(msg string, fields ...map[string]interface{}) {
	event := logger.Error()
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Dangerf logs a formatted error message
func Dangerf(format string, v ...interface{}) {
	logger.Error().Msgf(format, v...)
}

// Debug logs a debug message
func Debug(msg string, fields ...map[string]interface{}) {
	event := logger.Debug()
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Debugf logs a formatted debug message
func Debugf(format string, v ...interface{}) {
	logger.Debug().Msgf(format, v...)
}

// With returns a logger with additional fields
func With(fields map[string]interface{}) zerolog.Logger {
	event := logger.With()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	return event.Logger()
}

// GetLogger returns the underlying zerolog.Logger instance
func GetLogger() zerolog.Logger {
	return logger
}

// Legacy compatibility functions that accept any type like the original
func InfoAny(v ...any) {
	logger.Info().Msgf("%v", v...)
}

func WarningAny(v ...any) {
	logger.Warn().Msgf("%v", v...)
}

func DangerAny(v ...any) {
	logger.Error().Msgf("%v", v...)
}

