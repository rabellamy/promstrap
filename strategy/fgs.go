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
	// Latency is the time it takes to service a request. It’s important to
	// distinguish between the latency of successful requests and the latency of failed
	// requests.
	Latency *Distribution
	// Traffic is a measure of how much demand is being placed on your system,
	// measured in a high-level system-specific metric. For a web service, this
	// measurement is usually HTTP requests per second, perhaps broken out by the
	// nature of the requests (e.g., static versus dynamic content).
	Traffic *prometheus.CounterVec
	// Errors is the rate of requests that fail, either explicitly (e.g., HTTP 500s),
	// implicitly (for example, an HTTP 200 success response, but coupled with the wrong content)
	Errors *prometheus.CounterVec
	// Saturation is How "full" your service is. A measure of your system fraction,
	// emphasizing the resources that are most constrained. The degree to which extra
	// work is queued (or denied) that can't be serviced (e.g., in a memory-constrained system,
	// show memory; in an I/O-constrained system, show I/O or another example could be
	// scheduler run queue length).
	Saturation *prometheus.GaugeVec
}

// FourGoldenSignalsOpts is the options for a FourGoldenSignals strategy.
type FourGoldenSignalsOpts struct {
	Namespace     string   `validate:"required"`
	LatencyName   string   `validate:"required"`
	LatencyHelp   string   `validate:"required"`
	LatencyLabels []string `validate:"required"`
	// Buckets defines the histogram buckets into which observations are counted. Each
	// element in the slice is the upper inclusive bound of a bucket.
	LatencyBuckets []float64
	// Objectives defines the summary quantile rank estimates with their respective
	// absolute error.
	LatencyObjectives map[float64]float64
	TrafficName       string   `validate:"required"`
	TrafficHelp       string   `validate:"required"`
	TrafficLabels     []string `validate:"required"`
	SaturationName    string   `validate:"required"`
	SaturationHelp    string   `validate:"required"`
	SaturationLabels  []string `validate:"required"`
}

// NewFourGoldenSignals creates a new FourGoldenSignals strategy.
func NewFourGoldenSignals(opts FourGoldenSignalsOpts) (*FourGoldenSignals, error) {
	validate := validator.New()
	if err := validate.Struct(opts); err != nil {
		return nil, err
	}

	latency, err := NewDistribution(DistributionOpts{
		Namespace:  opts.Namespace,
		Name:       opts.LatencyName,
		Help:       opts.LatencyHelp,
		Labels:     opts.LatencyLabels,
		Buckets:    opts.LatencyBuckets,
		Objectives: opts.LatencyObjectives,
	})
	if err != nil {
		return nil, err
	}

	traffic, err := metrics.NewCounterWithLabels(metrics.CounterOpts{
		Namespace: opts.Namespace,
		Name:      opts.TrafficName,
		Help:      opts.TrafficHelp,
		Labels:    opts.TrafficLabels,
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

	saturation, err := metrics.NewGaugeWithLabels(metrics.GaugeOpts{
		Namespace: opts.Namespace,
		Name:      opts.SaturationName,
		Help:      opts.SaturationHelp,
		Labels:    opts.SaturationLabels,
	})
	if err != nil {
		return nil, err
	}

	return &FourGoldenSignals{
		Latency:    latency,
		Traffic:    traffic,
		Errors:     errors,
		Saturation: saturation,
	}, nil
}

// Register registers the FourGoldenSignals strategy with the Prometheus
// DefaultRegisterer.
func (f FourGoldenSignals) Register() error {
	err := RegisterStrategyFields(f)
	if err != nil {
		return err
	}

	return nil
}
