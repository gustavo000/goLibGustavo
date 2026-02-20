package logger

import "testing"

func TestLogger_BasicLevels(t *testing.T) {
	Info("info message")
	Warning("warning message")
	Danger("danger message")
}

