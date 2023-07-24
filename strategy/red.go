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
}

// REDOpts is the options to create a RED strategy.
type REDOpts struct {
	RequestType    string   `validate:"required"`
	Namespace      string   `validate:"required"`
	RequestLabels  []string `validate:"required"`
	DurationLabels []string `validate:"required"`
	// Buckets defines the histogram buckets into which observations are counted. Each
	// element in the slice is the upper inclusive bound of a bucket.
	Buckets []float64
	// Objectives defines the summary quantile rank estimates with their respective
	// absolute error.
	Objectives map[float64]float64
}

// NewRED creates a RED strategy.
func NewRED(opts REDOpts) (*RED, error) {
	validate := validator.New()
	if err := validate.Struct(opts); err != nil {
		return nil, err
	}

	requests, err := metrics.NewCounterWithLabels(metrics.CounterOpts{
		Namespace: opts.Namespace,
		Name:      fmt.Sprintf("%s_requests_total", opts.RequestType),
		Help:      "Number of requests",
		Labels:    opts.RequestLabels,
	})
	if err != nil {
		return nil, err
	}

	errors, err := metrics.NewCounterWithLabels(metrics.CounterOpts{
		Namespace: opts.Namespace,
		Name:      "errors_total",
		Help:      "Number of errors",
		Labels:    []string{"error"},
	})
	if err != nil {
		return nil, err
	}

	duration, err := NewDistribution(DistributionOpts{
		Namespace:  opts.Namespace,
		Name:       fmt.Sprintf("%s_request_duration_seconds_total", opts.RequestType),
		Help:       "Duration of request in seconds",
		Labels:     opts.DurationLabels,
		Buckets:    opts.Buckets,
		Objectives: opts.Objectives,
	})
	if err != nil {
		return nil, err
	}

	return &RED{
		Requests: requests,
		Errors:   errors,
		Duration: duration,
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
