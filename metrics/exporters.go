package metrics

import (
	"context"
	"fmt"

	"github.com/ramin/waypoint/config"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

type Exporter int

const (
	Prometheus Exporter = iota
	OTLPHTTP
)

func NewExporter() (interface{}, error) {
	switch config.Read().MetricsExporter {
	case "PROMETHEUS", "prometheus", "Prometheus", "PROM", "prom", "PROMETHEUS_EXPORTER", "prometheus_exporter", "PrometheusExporter", "PROM_EXPORTER", "prom_exporter", "PromExporter":
		// Return a Prometheus exporter instance
		return prometheusExporter()
	case "OTLPHTTP", "otlphttp", "OTLP_HTTP", "otlp_http", "OTLP-HTTP":
		// Return an OTLP HTTP exporter instance
		return oltphttpExporter()
	default:
		return nil, fmt.Errorf("Unknown exporter: %s", config.Read().MetricsExporter)
	}
}

func prometheusExporter() (metric.Option, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	option := metric.WithReader(exporter)
	return option, nil
}

func oltphttpExporter() (metric.Option, error) {
	exporter, err := otlpmetrichttp.New(
		context.Background(),
		otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression),
		otlpmetrichttp.WithEndpoint(config.Read().MetricsEndpoint),
	)
	if err != nil {
		return nil, err
	}

	option := metric.WithReader(
		metric.NewPeriodicReader(
			exporter,
		))

	return option, nil
}
