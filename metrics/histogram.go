package metrics

import (
	"github.com/go-playground/validator"
	"github.com/prometheus/client_golang/prometheus"
)

// HistogramOpts is the options for a Prometheus histogram.
type HistogramOpts struct {
	Namespace string   `validate:"required"`
	Name      string   `validate:"required"`
	Help      string   `validate:"required"`
	Labels    []string `validate:"required"`
	// Buckets defines the buckets into which observations are counted. Each
	// element in the slice is the upper inclusive bound of a bucket.
	Buckets []float64
}

// NewHistogramWithLabels creates a Prometheus histogram with labels based on the
// provided HistogramOpts.
// A histogram samples observations (usually things like request durations or
// response sizes) and counts them in configurable buckets. It also provides
// a sum of all observed values.
func NewHistogramWithLabels(opts HistogramOpts) (*prometheus.HistogramVec, error) {
	validate := validator.New()
	if err := validate.Struct(opts); err != nil {
		return nil, err
	}

	pOpts := prometheus.HistogramOpts{
		Namespace: opts.Namespace,
		Help:      opts.Help,
		Name:      opts.Name,
	}

	if opts.Buckets != nil {
		pOpts.Buckets = opts.Buckets
	}

	return prometheus.NewHistogramVec(pOpts, opts.Labels), nil
}
