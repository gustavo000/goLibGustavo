package observability

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type LokiLogger struct {
	logger zerolog.Logger
}

type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Service   string                 `json:"service"`
	Operation string                 `json:"operation"`
	ProductID *int                   `json:"product_id,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

type lokiPushRequest struct {
	Streams []lokiStream `json:"streams"`
}

type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

func NewLokiLogger() *LokiLogger {
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "bff-service").
		Str("component", "products").
		Str("version", "1.0.0").
		Logger()

	return &LokiLogger{
		logger: logger,
	}
}

func (l *LokiLogger) LogProductOperation(ctx context.Context, operation string, productID *int, details map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   fmt.Sprintf("Product %s operation", operation),
		Service:   "bff-service",
		Operation: operation,
		ProductID: productID,
		Details:   details,
	}

	event := l.logger.Info().
		Str("operation", operation).
		Str("service", "bff-service").
		Str("component", "products")

	if productID != nil {
		event = event.Int("product_id", *productID)
	}

	for key, value := range details {
		event = event.Interface(key, value)
	}

	event.Msg(fmt.Sprintf("Product %s completed", operation))
	l.sendToLoki(entry)
}

func (l *LokiLogger) LogProductError(ctx context.Context, operation string, productID *int, err error, details map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     "error",
		Message:   fmt.Sprintf("Product %s operation failed", operation),
		Service:   "bff-service",
		Operation: operation,
		ProductID: productID,
		Details:   details,
		Error:     err.Error(),
	}

	event := l.logger.Error().
		Str("operation", operation).
		Str("service", "bff-service").
		Str("component", "products").
		Err(err)
/
	if productID != nil {
		event = event.Int("product_id", *productID)
	}

	for key, value := range details {
		event = event.Interface(key, value)
	}

	event.Msg(fmt.Sprintf("Product %s failed", operation))
	l.sendToLoki(entry)
}

func (l *LokiLogger) sendToLoki(entry LogEntry) {
	lokiURL := os.Getenv("LOKI_URL")
	if lokiURL == "" {
		lokiURL = "http://localhost:3100"
	}

	line, err := json.Marshal(entry)
	if err != nil {
		l.logger.Error().Err(err).Msg("Failed to marshal log entry for Loki")
		return
	}

	labels := map[string]string{
		"service":   entry.Service,
		"component": "products",
		"operation": entry.Operation,
		"level":     entry.Level,
	}

	payload := lokiPushRequest{
		Streams: []lokiStream{
			{
				Stream: labels,
				Values: [][]string{
					{strconv.FormatInt(entry.Timestamp.UnixNano(), 10), string(line)},
				},
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		l.logger.Error().Err(err).Msg("Failed to marshal Loki push payload")
		return
	}

	go func() {
		req, err := http.NewRequest(http.MethodPost, lokiURL+"/loki/api/v1/push", bytes.NewReader(jsonData))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()
	}()
}
