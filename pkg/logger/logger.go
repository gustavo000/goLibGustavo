package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var (
	// Global logger instance
	logger zerolog.Logger
)

// Initialize the logger with configuration
func init() {
	// Set time format to RFC3339
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

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

func Error(msg string, err error, fields ...map[string]interface{}) {
	event := logger.Error()
	if err != nil {
		event = event.Err(err)
	}
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

func Fatal(msg string, err error, fields ...map[string]interface{}) {
	event := logger.Fatal()
	if err != nil {
		event = event.Err(err)
	}
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Legacy compatibility methods (maintaining your original API)
func InfoLegacy(v ...interface{}) {
	logger.Info().Msgf("%v", v...)
}

func WarningLegacy(v ...interface{}) {
	logger.Warn().Msgf("%v", v...)
}

func DangerLegacy(v ...interface{}) {
	logger.Error().Msgf("%v", v...)
}

// GetLogger returns the configured zerolog instance
func GetLogger() zerolog.Logger {
	return logger
}
