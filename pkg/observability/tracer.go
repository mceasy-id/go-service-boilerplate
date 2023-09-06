package observability

import (
	"context"
	"mceasy/service-demo/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func InitTracerProvider(config *config.Config) (*sdktrace.TracerProvider, error) {
	var (
		exporter sdktrace.SpanExporter
		err      error
	)

	if config.Observability.Mode == "console" {
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	} else if config.Observability.Mode == "otlp" {
		exporter, err = otlptracegrpc.New(
			context.Background(),
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(config.Observability.OtlpEndpoint),
		)
	}

	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(initResource(config.App.Name, config.App.Version, config.App.Env)),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tracerProvider, nil
}
