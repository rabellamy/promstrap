package strategy

import (
	"fmt"

	"github.com/go-playground/validator"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rabellamy/promstrap/metrics"
)

// RED describes a set of metrics that work well for monitoring request-handling
// services (like an HTTP API server or a database server).
// https://www.slideshare.net/weaveworks/monitoring-microservices
type RED struct {
	// The number of requests per second.
	Requests *prometheus.CounterVec
	// Errors is the rate of requests that fail, either explicitly (e.g., HTTP 500s),
	// implicitly (for example, an HTTP 200 success response, but coupled with the wrong content).
	Errors *prometheus.CounterVec
	// Distributions of the amount of time each request takes
	Duration *Distribution

	redOpts REDOpts
}

type REDRequestsOpt struct {
	// RequestName is the name of the requests metric. If not specified, defaults to "{RequestType}_requests_total".
	RequestName string
	// RequestType is the type of the request (e.g., "http", "grpc").
	RequestType string `validate:"required"`
	// RequestLabels are the labels to attach to the requests metric.
	RequestLabels []string `validate:"required"`
}

type REDErrorsOpt struct {
	// ErrorName is the name of the errors metric. If not specified, defaults to "errors_total".
	ErrorName string
	// ErrorLabels are the labels to attach to the errors metric.
	ErrorLabels []string `validate:"required"`
}

type REDDurationOpt struct {
	// DurationName is the name of the duration metric. If not specified, defaults to "{RequestType}_request_duration_seconds".
	DurationName string
	// DurationLabels are the labels to attach to the duration metric.
	DurationLabels []string `validate:"required"`
	// Buckets defines the histogram buckets into which observations are counted. Each
	// element in the slice is the upper inclusive bound of a bucket.
	Buckets []float64
	// Objectives defines the summary quantile rank estimates with their respective
	// absolute error.
	Objectives map[float64]float64
}

// REDOpts is the options to create a RED strategy.
type REDOpts struct {
	Namespace   string         `validate:"required"`
	RequestsOpt REDRequestsOpt `validate:"required"`
	ErrorsOpt   REDErrorsOpt   `validate:"required"`
	DurationOpt REDDurationOpt `validate:"required"`
}

// NewRED creates a RED strategy.
func NewRED(opts REDOpts) (*RED, error) {
	validate := validator.New()
	if err := validate.Struct(opts); err != nil {
		return nil, err
	}

	requestsName := getREDRequestsMetricName(opts)
	requests, err := metrics.NewCounterWithLabels(metrics.CounterOpts{
		Namespace: opts.Namespace,
		Name:      requestsName,
		Help:      "Number of requests",
		Labels:    opts.RequestsOpt.RequestLabels,
	})
	if err != nil {
		return nil, err
	}

	errorsName := getREDErrorsMetricName(opts)
	errors, err := metrics.NewCounterWithLabels(metrics.CounterOpts{
		Namespace: opts.Namespace,
		Name:      errorsName,
		Help:      "Number of errors, RED",
		Labels:    opts.ErrorsOpt.ErrorLabels,
	})
	if err != nil {
		return nil, err
	}

	durationName := getREDDurationMetricName(opts)
	duration, err := NewDistribution(DistributionOpts{
		Namespace:  opts.Namespace,
		Name:       durationName,
		Help:       "Duration of request in seconds",
		Labels:     opts.DurationOpt.DurationLabels,
		Buckets:    opts.DurationOpt.Buckets,
		Objectives: opts.DurationOpt.Objectives,
	})
	if err != nil {
		return nil, err
	}

	return &RED{
		Requests: requests,
		Errors:   errors,
		Duration: duration,
		redOpts:  opts,
	}, nil
}

// Register registers the RED strategy with the Prometheus DefaultRegisterer.
func (r RED) Register() error {
	err := RegisterStrategyFields(r)
	if err != nil {
		return err
	}

	return nil
}

func (r RED) RequestCounterName() string {
	return getREDRequestsMetricName(r.redOpts)
}

func (r RED) ErrorCounterName() string {
	return getREDErrorsMetricName(r.redOpts)
}

func (r RED) DurationHistogramName() string {
	return getREDDurationMetricName(r.redOpts)
}

func getREDRequestsMetricName(opts REDOpts) string {
	if opts.RequestsOpt.RequestName != "" {
		return opts.RequestsOpt.RequestName
	}
	return fmt.Sprintf("%s_requests_total", opts.RequestsOpt.RequestType)
}

func getREDErrorsMetricName(opts REDOpts) string {
	if opts.ErrorsOpt.ErrorName != "" {
		return opts.ErrorsOpt.ErrorName
	}
	return "errors_total"
}

func getREDDurationMetricName(opts REDOpts) string {
	if opts.DurationOpt.DurationName != "" {
		return opts.DurationOpt.DurationName
	}
	return fmt.Sprintf("%s_request_duration_seconds", opts.RequestsOpt.RequestType)
}
