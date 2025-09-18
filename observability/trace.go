package observability

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

func NewOTLPTraceExporter(ctx context.Context, otlpEndpoint string) *otlptrace.Exporter {
	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otlpEndpoint))
	traceExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		log.Fatalf("Failed to create the collector trace exporter: %v", err)
	}
	return traceExp
}
