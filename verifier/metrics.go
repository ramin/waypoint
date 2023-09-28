package verifier

import (
	"log"

	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

type Metrics struct {
	// counters
	TotalBlocks api.Int64Counter
	Reads       api.Int64Counter
	Misses      api.Int64Counter
	Failures    api.Int64Counter
	Errors      api.Int64Counter
	Writes      api.Int64Counter
	// timers
}

func NewMetrics() (*Metrics, error) {
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter("celestia/waypoint")

	m := &Metrics{}
	err = m.injectCounters(meter)
	if err != nil {
		return nil, err
	}
	err = m.injectTimers(meter)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Metrics) injectCounters(meter api.Meter) error {
	total, err := meter.Int64Counter("totalBlocks", api.WithDescription("total blocks tested"))
	if err != nil {
		return err
	}

	m.TotalBlocks = total

	writes, err := meter.Int64Counter("writes", api.WithDescription("total blocks written to"))
	if err != nil {
		return err
	}
	m.Writes = writes

	reads, err := meter.Int64Counter("reads", api.WithDescription("total blocks tested"))
	if err != nil {
		return err
	}

	m.Reads = reads

	misses, err := meter.Int64Counter("misses", api.WithDescription("total read misses"))
	if err != nil {
		return err
	}

	m.Misses = misses

	failures, err := meter.Int64Counter("failures", api.WithDescription("total write failures"))
	if err != nil {
		return err
	}

	m.Failures = failures

	errors, err := meter.Int64Counter("errors", api.WithDescription("total misc failures"))
	if err != nil {
		return err
	}

	m.Failures = errors

	return nil
}

func (m *Metrics) injectTimers(meter api.Meter) error {
	return nil
}
