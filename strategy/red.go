package strategy

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rabellamy/promstrap"
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
	Duration *prometheus.HistogramVec
}

// REDOpts is the options to create a RED strategy.
type REDOpts struct {
	RequestType    string   `validate:"required"`
	Namespace      string   `validate:"required"`
	RequestLabels  []string `validate:"required"`
	DurationLabels []string `validate:"required"`
}

// NewRED creates a RED strategy.
func NewRED(opts REDOpts) (*RED, error) {
	if err := promstrap.Validate(opts); err != nil {
		return nil, err
	}

	requestsCounter, err := promstrap.NewCounterWithLabels(promstrap.MetricsOpts{
		Namespace: opts.Namespace,
		Name:      fmt.Sprintf("%s_requests_total", opts.RequestType),
		Help:      "Number of requests",
		Labels:    opts.RequestLabels,
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

	durationHistogram, err := promstrap.NewHistogramWithLabels(promstrap.MetricsOpts{
		Namespace: opts.Namespace,
		Name:      fmt.Sprintf("%s_request_duration_seconds_total", opts.RequestType),
		Help:      "Duration of request in seconds",
		Labels:    opts.DurationLabels,
	})
	if err != nil {
		return nil, err
	}

	return &RED{
		Requests: requestsCounter,
		Errors:   errorsCounter,
		Duration: durationHistogram,
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
