package observability

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// InitTelemetry initializes OpenTelemetry using environment variables
// This approach uses the standard OpenTelemetry auto-instrumentation pattern
func InitTelemetry(serviceName string) (func(), error) {
	var exporter sdktrace.SpanExporter
	var err error

	// Use environment variables for OTLP configuration
	// This avoids the URL encoding issues with manual configuration

	// Get service name from environment or use provided name
	serviceNameFromEnv := os.Getenv("OTEL_SERVICE_NAME")
	if serviceNameFromEnv != "" {
		serviceName = serviceNameFromEnv
	}

	// Create resource with service metadata
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("1.0.0"),
			semconv.ServiceInstanceIDKey.String(fmt.Sprintf("%s-%d", serviceName, os.Getpid())),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Try to use gRPC OTLP exporter to avoid HTTP URL encoding issues
	otlpEndpoint := normalizeOTLPEndpoint(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if otlpEndpoint != "" {
		fmt.Printf("Initializing OTLP gRPC exporter for endpoint: %s\n", otlpEndpoint)
		exporter, err = otlptracegrpc.New(context.Background(),
			otlptracegrpc.WithEndpoint(otlpEndpoint),
			otlptracegrpc.WithInsecure(), // Use insecure connection for local development
		)
		if err != nil {
			fmt.Printf("Failed to create OTLP gRPC exporter, falling back to stdout: %v\n", err)
			exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
			if err != nil {
				return nil, fmt.Errorf("failed to create stdout exporter: %w", err)
			}
		}
	} else {
		fmt.Println("No OTLP endpoint configured, using stdout exporter")
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, fmt.Errorf("failed to create stdout exporter: %w", err)
		}
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // Sample all traces for development
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			fmt.Printf("Error shutting down tracer provider: %v\n", err)
		}
	}, nil
}

func normalizeOTLPEndpoint(endpoint string) string {
	endpoint = strings.TrimSpace(endpoint)
	endpoint = strings.TrimPrefix(endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "http://")
	return endpoint
}
