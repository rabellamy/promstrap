package strategy

import (
	"fmt"

	"github.com/go-playground/validator"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rabellamy/promstrap/metrics"
)

// Distribution encapsulates the two Prometheus metric types used
// to track the distribution of a set of observed values.
type Distribution struct {
	Histogram *prometheus.HistogramVec
	Summary   *prometheus.SummaryVec
	opts      DistributionOpts
}

// DistributionOpts is the options to create a Distribution
// strategy.
type DistributionOpts struct {
	Namespace string   `validate:"required"`
	Name      string   `validate:"required"`
	Help      string   `validate:"required"`
	Labels    []string `validate:"required"`
	// Buckets defines the histogram buckets into which observations are counted.
	// Each element in the slice is the upper inclusive bound of a bucket.
	Buckets []float64
	// Objectives defines the summary quantile rank estimates with their respective
	// absolute error.
	Objectives map[float64]float64
}

// NewDistribution creates a Distribution.
func NewDistribution(opts DistributionOpts) (*Distribution, error) {
	validate := validator.New()
	if err := validate.Struct(opts); err != nil {
		return nil, err
	}

	histogramName := getDistributionHistogramName(opts)
	histogram, err := metrics.NewHistogramWithLabels(metrics.HistogramOpts{
		Namespace: opts.Namespace,
		Name:      histogramName,
		Help:      opts.Help,
		Labels:    opts.Labels,
		Buckets:   opts.Buckets,
	})
	if err != nil {
		return nil, err
	}

	summaryName := getDistributionSummaryName(opts)
	summary, err := metrics.NewSummaryWithLabels(metrics.SummaryOpts{
		Namespace:  opts.Namespace,
		Name:       summaryName,
		Help:       opts.Help,
		Labels:     opts.Labels,
		Objectives: opts.Objectives,
	})
	if err != nil {
		return nil, err
	}

	return &Distribution{
		Histogram: histogram,
		Summary:   summary,
	}, nil
}

// Register registers the Distribution strategy with the Prometheus
// DefaultRegisterer.
func (r Distribution) Register() error {
	err := RegisterStrategyFields(r)
	if err != nil {
		return err
	}

	return nil
}

func (r Distribution) HistogramName() string {
	return getDistributionHistogramName(r.opts)
}

func (r Distribution) SummaryName() string {
	return getDistributionSummaryName(r.opts)
}

func getDistributionHistogramName(opts DistributionOpts) string {
	return fmt.Sprintf("%s_hist", opts.Name)
}

func getDistributionSummaryName(opts DistributionOpts) string {
	return fmt.Sprintf("%s_sum", opts.Name)
}
