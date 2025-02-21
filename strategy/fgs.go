package strategy

import (
	"github.com/go-playground/validator"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rabellamy/promstrap/metrics"
)

// FourGoldenSignals of monitoring are latency, traffic, errors, and saturation.
// If you can only measure four metrics of your user-facing system, focus on these four.
// https://sre.google/sre-book/monitoring-distributed-systems/
type FourGoldenSignals struct {
	Latency    *Distribution
	Traffic    *prometheus.CounterVec
	Errors     *prometheus.CounterVec
	Saturation *prometheus.GaugeVec

	fgsOpts FourGoldenSignalsOpts
}

type FGSLatencyOpt struct {
	// LatencyName must be explicitly set.
	// Should follow the pattern {system}_{type}_latency_seconds
	// Examples:
	// - http_request_latency_seconds
	// - grpc_method_latency_seconds
	// - database_query_latency_seconds
	LatencyName string `validate:"required"`
	// LatencyType is the type of latency being measured (e.g., "http", "grpc", "database")
	LatencyType string `validate:"required"`
	// LatencyHelp is the help text for the latency metric
	LatencyHelp string `validate:"required"`
	// LatencyLabels are the labels to attach to the latency metric
	LatencyLabels []string `validate:"required"`
	// Buckets defines the histogram buckets into which observations are counted
	Buckets []float64
	// Objectives defines the quantile rank estimates with their respective absolute error
	Objectives map[float64]float64
}

type FGSTrafficOpt struct {
	// TrafficName must be explicitly set.
	// Should follow the pattern {system}_{type}_requests_total
	// Examples:
	// - http_server_requests_total
	// - grpc_client_requests_total
	// - database_operations_total
	TrafficName string `validate:"required"`
	// TrafficType is the type of traffic being measured (e.g., "http", "grpc")
	TrafficType string `validate:"required"`
	// TrafficHelp is the help text for the traffic metric
	TrafficHelp string `validate:"required"`
	// TrafficLabels are the labels to attach to the traffic metric
	TrafficLabels []string `validate:"required"`
}

type FGSErrorsOpt struct {
	// ErrorName is the name of the errors metric. If not specified, defaults to "errors_total"
	ErrorName string
	// ErrorHelp is the help text for the errors metric
	ErrorHelp string `validate:"required"`
	// ErrorLabels are the labels to attach to the errors metric
	ErrorLabels []string `validate:"required"`
}

type FGSSaturationOpt struct {
	// SaturationName must be explicitly set.
	// Should follow the pattern {resource}_{type}_saturation_{unit}
	// Examples:
	// - memory_heap_saturation_bytes
	// - threadpool_worker_saturation_ratio
	// - connection_pool_saturation_percent
	SaturationName string `validate:"required"`
	// SaturationHelp is the help text for the saturation metric
	SaturationHelp string `validate:"required"`
	// SaturationLabels are the labels to attach to the saturation metric
	SaturationLabels []string `validate:"required"`
}

type FourGoldenSignalsOpts struct {
	Namespace     string           `validate:"required"`
	LatencyOpt    FGSLatencyOpt    `validate:"required"`
	TrafficOpt    FGSTrafficOpt    `validate:"required"`
	ErrorsOpt     FGSErrorsOpt     `validate:"required"`
	SaturationOpt FGSSaturationOpt `validate:"required"`
}

func NewFourGoldenSignals(opts FourGoldenSignalsOpts) (*FourGoldenSignals, error) {
	validate := validator.New()
	if err := validate.Struct(opts); err != nil {
		return nil, err
	}

	latencyName := getFGSLatencyMetricName(opts)
	latency, err := NewDistribution(DistributionOpts{
		Namespace:  opts.Namespace,
		Name:       latencyName,
		Help:       opts.LatencyOpt.LatencyHelp,
		Labels:     opts.LatencyOpt.LatencyLabels,
		Buckets:    opts.LatencyOpt.Buckets,
		Objectives: opts.LatencyOpt.Objectives,
	})
	if err != nil {
		return nil, err
	}

	trafficName := getFGSTrafficMetricName(opts)
	traffic, err := metrics.NewCounterWithLabels(metrics.CounterOpts{
		Namespace: opts.Namespace,
		Name:      trafficName,
		Help:      opts.TrafficOpt.TrafficHelp,
		Labels:    opts.TrafficOpt.TrafficLabels,
	})
	if err != nil {
		return nil, err
	}

	errorsName := getFGSErrorMetricName(opts)
	errors, err := metrics.NewCounterWithLabels(metrics.CounterOpts{
		Namespace: opts.Namespace,
		Name:      errorsName,
		Help:      opts.ErrorsOpt.ErrorHelp,
		Labels:    opts.ErrorsOpt.ErrorLabels,
	})
	if err != nil {
		return nil, err
	}

	saturationName := getFGSSaturationMetricName(opts)
	saturation, err := metrics.NewGaugeWithLabels(metrics.GaugeOpts{
		Namespace: opts.Namespace,
		Name:      saturationName,
		Help:      opts.SaturationOpt.SaturationHelp,
		Labels:    opts.SaturationOpt.SaturationLabels,
	})
	if err != nil {
		return nil, err
	}

	return &FourGoldenSignals{
		Latency:    latency,
		Traffic:    traffic,
		Errors:     errors,
		Saturation: saturation,
		fgsOpts:    opts,
	}, nil
}

func (f FourGoldenSignals) Register() error {
	err := RegisterStrategyFields(f)
	if err != nil {
		return err
	}

	return nil
}

func (f FourGoldenSignals) LatencyMetricName() string {
	return getFGSLatencyMetricName(f.fgsOpts)
}

func (f FourGoldenSignals) TrafficMetricName() string {
	return getFGSTrafficMetricName(f.fgsOpts)
}

func (f FourGoldenSignals) ErrorMetricName() string {
	return getFGSErrorMetricName(f.fgsOpts)
}

func (f FourGoldenSignals) SaturationMetricName() string {
	return getFGSSaturationMetricName(f.fgsOpts)
}

func getFGSLatencyMetricName(opts FourGoldenSignalsOpts) string {
	return opts.LatencyOpt.LatencyName
}

func getFGSTrafficMetricName(opts FourGoldenSignalsOpts) string {
	return opts.TrafficOpt.TrafficName
}

func getFGSErrorMetricName(opts FourGoldenSignalsOpts) string {
	if opts.ErrorsOpt.ErrorName != "" {
		return opts.ErrorsOpt.ErrorName
	}
	return "errors_total"
}

func getFGSSaturationMetricName(opts FourGoldenSignalsOpts) string {
	return opts.SaturationOpt.SaturationName
}
