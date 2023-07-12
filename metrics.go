package promstrap

import (
	"errors"
	"reflect"

	"github.com/go-playground/validator"
	"github.com/prometheus/client_golang/prometheus"
)

// MetricsOpts is the options for any of the basic Prometheus metric types.
type MetricsOpts struct {
	Namespace string   `validate:"required"`
	Name      string   `validate:"required"`
	Help      string   `validate:"required"`
	Labels    []string `validate:"required"`
}

// NewCounterWithLabels creates a Prometheus counter with labels based on the
// provided MetricsOpts.
// A counter is a cumulative metric that represents a single monotonically
// increasing counter whose value can only increase or be reset to zero on restart.
// Counters are for tracking cumulative totals over time, like the total number
// of HTTP requests or the number of errors.
func NewCounterWithLabels(opts MetricsOpts) (*prometheus.CounterVec, error) {
	if err := Validate(opts); err != nil {
		return nil, err
	}

	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: opts.Namespace,
		Name:      opts.Name,
		Help:      opts.Help,
	}, opts.Labels), nil
}

// NewGaugeWithLabels creates a Prometheus Gauge with labels based on the
// provided MetricsOpts.
// A gauge is a metric that represents a single numerical value that can
// arbitrarily go up and down. Gauges track values that can change over time,
// such as the amount of memory used, the number of requests in progress,
// or the temperature of a device.
func NewGaugeWithLabels(opts MetricsOpts) (*prometheus.GaugeVec, error) {
	if err := Validate(opts); err != nil {
		return nil, err
	}

	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: opts.Namespace,
		Name:      opts.Name,
		Help:      opts.Help,
	}, opts.Labels), nil
}

// NewHistogramWithLabels creates a Prometheus histogram with labels based on the
// provided MetricsOpts.
// A histogram samples observations (usually things like request durations or
// response sizes) and counts them in configurable buckets. It also provides
// a sum of all observed values.
// TODO: make Objectives configurable.
func NewHistogramWithLabels(opts MetricsOpts) (*prometheus.HistogramVec, error) {
	if err := Validate(opts); err != nil {
		return nil, err
	}

	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: opts.Namespace,
		Help:      opts.Help,
		Name:      opts.Name,
	}, opts.Labels), nil
}

// NewSummaryWithLabels creates a Prometheus summary with labels based on the
// provided MetricsOpts.
// A summary samples observations (usually things like request durations and
// response sizes). While it also provides a total count of observations and a
// sum of all observed values, it calculates configurable quantiles over
// a sliding time window.
// TODO: make Buckets and NativeHistogramBucketFactor configurable.
func NewSummaryWithLabels(opts MetricsOpts) (*prometheus.SummaryVec, error) {
	if err := Validate(opts); err != nil {
		return nil, err
	}

	return prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: opts.Namespace,
		Help:      opts.Help,
		Name:      opts.Name,
	}, opts.Labels), nil
}

// RegisterMetrics registers the basic Prometheus metrics and the well known strategies.
func RegisterMetrics(metrics interface{}) error {
	errPromType := errors.New("metric is not a Prometheus Collector")

	switch v := metrics.(type) {
	case *prometheus.CounterVec,
		*prometheus.GaugeVec,
		*prometheus.HistogramVec,
		*prometheus.SummaryVec:
		m, ok := v.(prometheus.Collector)
		if !ok {
			return errPromType
		}

		prometheus.MustRegister(m)

	// registers the complex metrics strategies, metrics built as custom types
	// using structs, like the well known strategies (USE/RED/FourGoldenSignals).
	default:
		if reflect.ValueOf(v).Kind() != reflect.Struct {
			return errors.New("complex metric is not built with a stuct")
		}

		l := reflect.TypeOf(v).NumField()
		if l < 1 {
			return errors.New("metric strategies needs at least one field")
		}

		values := reflect.ValueOf(v)

		for i := 0; i < l; i++ {
			m, ok := values.Field(i).Interface().(prometheus.Collector)
			if !ok {
				return errPromType
			}

			prometheus.MustRegister(m)
		}
	}

	return nil
}

// validate is a helper function to make sure that any required fields exist.
func Validate(opts interface{}) error {
	validate := validator.New()

	return validate.Struct(opts)
}