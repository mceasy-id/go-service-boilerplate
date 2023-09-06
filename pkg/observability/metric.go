package observability

import (
	"context"
	"mceasy/service-demo/config"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func InitMeterProvider(config *config.Config) (*sdkmetric.MeterProvider, error) {
	exporter, err := otlpmetricgrpc.New(
		context.Background(),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(config.Observability.OtlpEndpoint),
	)
	if err != nil {
		return nil, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(2*time.Second))),
		sdkmetric.WithResource(initResource(config.App.Name, config.App.Version, config.App.Env)),
	)

	otel.SetMeterProvider(mp)

	return mp, nil
}
