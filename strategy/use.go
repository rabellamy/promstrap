package strategy

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rabellamy/promstrap"
)

// USE describes a set of metrics that are useful for measuring the performance
// and usage of things that behave like resources that can be either used or
// unused (queues, CPUs, memory, disk space, etc.).
// https://www.brendangregg.com/usemethod.html
type USE struct {
	// Utilization is the average time that the resource was busy servicing work,
	// a percent over a time interval. The average time that the resource was
	// busy (e.g. one disk at 90% I/O utilization).
	Utilization *prometheus.GaugeVec
	// Saturation is How "full" your service is. A measure of your system fraction,
	// emphasizing the resources that are most constrained. The degree to which extra
	// work is queued (or denied) that can't be serviced (e.g., in a memory-constrained system,
	// show memory; in an I/O-constrained system, show I/O or another example could be
	// scheduler run queue length).
	Saturation *prometheus.GaugeVec
	// Errors is the rate of requests that fail, either explicitly (e.g., HTTP 500s),
	// implicitly (for example, an HTTP 200 success response, but coupled with the wrong content).
	Errors *prometheus.CounterVec
}

// USEOpts is the options to create a USE strategy.
type USEOpts struct {
	Namespace         string   `validate:"required"`
	SaturationName    string   `validate:"required"`
	SaturationHelp    string   `validate:"required"`
	SaturationLabels  []string `validate:"required"`
	UtilizationName   string   `validate:"required"`
	UtilizationHelp   string   `validate:"required"`
	UtilizationLabels []string `validate:"required"`
}

// NEWUSE creates a USE strategy.
func NewUSE(opts USEOpts) (*USE, error) {
	if err := promstrap.Validate(opts); err != nil {
		return nil, err
	}

	utilizationGauge, err := promstrap.NewGaugeWithLabels(promstrap.MetricsOpts{
		Namespace: opts.Namespace,
		Name:      opts.UtilizationName,
		Help:      opts.UtilizationHelp,
		Labels:    opts.UtilizationLabels,
	})
	if err != nil {
		return nil, err
	}

	saturationGauge, err := promstrap.NewGaugeWithLabels(promstrap.MetricsOpts{
		Namespace: opts.Namespace,
		Name:      opts.SaturationName,
		Help:      opts.SaturationHelp,
		Labels:    opts.SaturationLabels,
	})
	if err != nil {
		return nil, err
	}

	errorsCounter, err := promstrap.NewCounterWithLabels(promstrap.MetricsOpts{
		Namespace: opts.Namespace,
		Name:      "errors_total",
		Help:      "Number of errors",
		Labels:    []string{"error"},
	})
	if err != nil {
		return nil, err
	}

	return &USE{
		Utilization: utilizationGauge,
		Saturation:  saturationGauge,
		Errors:      errorsCounter,
	}, nil
}

// Register registers the USE strategy with the Prometheus DefaultRegisterer.
func (u USE) Register() error {
	err := RegisterStrategyFields(u)
	if err != nil {
		return err
	}

	return nil
}
