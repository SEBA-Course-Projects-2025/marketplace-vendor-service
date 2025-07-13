package tracer

import (
	"context"
	"encoding/base64"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"log"
	"os"
)

var Tracer = otel.Tracer("vendor-service")

func InitTracer() func(ctx context.Context) error {

	ctx := context.Background()

	stackId := os.Getenv("GRAFANA_TEMPO_STACK_ID")
	apiKey := os.Getenv("GRAFANA_TEMPO_API_KEY")
	endpoint := os.Getenv("GRAFANA_TEMPO_ENDPOINT")

	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(stackId+":"+apiKey))

	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(endpoint), otlptracehttp.WithHeaders(map[string]string{
		"Authorization": auth,
	}), otlptracehttp.WithURLPath("otlp/v1/traces"))

	if err != nil {
		log.Fatalf("Error creating OpenTelemetry exporter: %v", err)
	}

	otelRes, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceName("vendor-service")))

	if err != nil {
		log.Fatalf("Error creating OpenTelemetry resource: %v", err)
	}

	provider := trace.NewTracerProvider(trace.WithBatcher(exporter), trace.WithResource(otelRes))

	otel.SetTracerProvider(provider)

	Tracer = otel.Tracer("vendor-service")

	return provider.Shutdown

}
