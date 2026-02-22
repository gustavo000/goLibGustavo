package logger

import (
	"bytes"
	"testing"

	"github.com/rs/zerolog"
)

func TestLogger_BasicLevels(t *testing.T) {
	Info("info message")
	Warning("warning message")
	Danger("danger message")
}

func TestLogger_WithFields(t *testing.T) {
	fields := map[string]interface{}{
		"user_id": 12345,
		"action":  "login",
	}
	
	Info("User action", fields)
	Debug("Debug info", fields)
	Warning("Warning with context", fields)
}

func TestLogger_FormattedMessages(t *testing.T) {
	Infof("User %s logged in at %s", "john.doe", "2024-01-01")
	Warningf("Failed attempt: %d", 3)
	Dangerf("Error: %v", "connection timeout")
}

func TestLogger_ErrorWithError(t *testing.T) {
	err := &testError{"test error"}
	Error("Something went wrong", err, map[string]interface{}{
		"component": "database",
		"retry":     true,
	})
}

func TestLogger_SetLogLevel(t *testing.T) {
	SetLogLevel("debug")
	Debug("This should be visible")
	
	SetLogLevel("error")
	Debug("This should not be visible")
	Error("This should be visible", nil)
}

func TestLogger_ProductionMode(t *testing.T) {
	SetProductionMode()
	Info("This will be in JSON format")
}

func TestLogger_LegacyCompatibility(t *testing.T) {
	InfoLegacy("legacy info message")
	WarningLegacy("legacy warning message")
	DangerLegacy("legacy danger message")
}

// Test error type
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

