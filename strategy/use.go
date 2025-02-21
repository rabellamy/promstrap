package strategy

import (
	"github.com/go-playground/validator"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rabellamy/promstrap/metrics"
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

	useOpts USEOpts
}

type USEUtilizationOpt struct {
	// UtilizationName is the name of the utilization metric. If not specified, defaults to "utilization".
	UtilizationName string `validate:"required"`
	// UtilizationHelp is the help text for the utilization metric.
	UtilizationHelp string `validate:"required"`
	// UtilizationLabels are the labels to attach to the utilization metric.
	UtilizationLabels []string `validate:"required"`
}

type USESaturationOpt struct {
	// SaturationName is the name of the saturation metric. If not specified, defaults to "saturation".
	SaturationName string `validate:"required"`
	// SaturationHelp is the help text for the saturation metric.
	SaturationHelp string `validate:"required"`
	// SaturationLabels are the labels to attach to the saturation metric.
	SaturationLabels []string `validate:"required"`
}

type USEErrorsOpt struct {
	// ErrorName is the name of the errors metric. If not specified, defaults to "errors_total".
	ErrorName string
	// ErrorLabels are the labels to attach to the errors metric.
	ErrorLabels []string `validate:"required"`
}

// USEOpts is the options to create a USE strategy.
type USEOpts struct {
	Namespace      string            `validate:"required"`
	UtilizationOpt USEUtilizationOpt `validate:"required"`
	SaturationOpt  USESaturationOpt  `validate:"required"`
	ErrorsOpt      USEErrorsOpt      `validate:"required"`
}

// NewUSE creates a USE strategy.
func NewUSE(opts USEOpts) (*USE, error) {
	validate := validator.New()
	if err := validate.Struct(opts); err != nil {
		return nil, err
	}

	utilizationName := getUSEUtilizationMetricName(opts)
	utilizationGauge, err := metrics.NewGaugeWithLabels(metrics.GaugeOpts{
		Namespace: opts.Namespace,
		Name:      utilizationName,
		Help:      opts.UtilizationOpt.UtilizationHelp,
		Labels:    opts.UtilizationOpt.UtilizationLabels,
	})
	if err != nil {
		return nil, err
	}

	saturationName := getUSESaturationMetricName(opts)
	saturationGauge, err := metrics.NewGaugeWithLabels(metrics.GaugeOpts{
		Namespace: opts.Namespace,
		Name:      saturationName,
		Help:      opts.SaturationOpt.SaturationHelp,
		Labels:    opts.SaturationOpt.SaturationLabels,
	})
	if err != nil {
		return nil, err
	}

	errorsName := getUSEErrorsMetricName(opts)
	errorsCounter, err := metrics.NewCounterWithLabels(metrics.CounterOpts{
		Namespace: opts.Namespace,
		Name:      errorsName,
		Help:      "Number of errors",
		Labels:    opts.ErrorsOpt.ErrorLabels,
	})
	if err != nil {
		return nil, err
	}

	return &USE{
		Utilization: utilizationGauge,
		Saturation:  saturationGauge,
		Errors:      errorsCounter,
		useOpts:     opts,
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

// UtilizationMetricName returns the name of the utilization metric.
func (u USE) UtilizationMetricName() string {
	return getUSEUtilizationMetricName(u.useOpts)
}

// SaturationMetricName returns the name of the saturation metric.
func (u USE) SaturationMetricName() string {
	return getUSESaturationMetricName(u.useOpts)
}

// ErrorMetricName returns the name of the errors metric.
func (u USE) ErrorMetricName() string {
	return getUSEErrorsMetricName(u.useOpts)
}

func getUSEUtilizationMetricName(opts USEOpts) string {
	if opts.UtilizationOpt.UtilizationName != "" {
		return opts.UtilizationOpt.UtilizationName
	}
	return "utilization"
}

func getUSESaturationMetricName(opts USEOpts) string {
	if opts.SaturationOpt.SaturationName != "" {
		return opts.SaturationOpt.SaturationName
	}
	return "saturation"
}

func getUSEErrorsMetricName(opts USEOpts) string {
	if opts.ErrorsOpt.ErrorName != "" {
		return opts.ErrorsOpt.ErrorName
	}
	return "errors_total"
}
