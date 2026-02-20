package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type level string

const (
	InfoLevel    level = "INFO"
	WarningLevel level = "WARNING"
	DangerLevel  level = "DANGER"
)

var mu sync.Mutex

func write(l level, w *os.File, v ...any) {
	mu.Lock()
	ts := time.Now().Format(time.RFC3339)
	fmt.Fprintln(w, ts, "["+string(l)+"]", fmt.Sprint(v...))
	mu.Unlock()
}

func Info(v ...any) {
	write(InfoLevel, os.Stdout, v...)
}

func Warning(v ...any) {
	write(WarningLevel, os.Stdout, v...)
}

func Danger(v ...any) {
	write(DangerLevel, os.Stderr, v...)
}

