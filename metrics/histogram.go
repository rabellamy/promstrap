package promstrap

import "github.com/prometheus/client_golang/prometheus"

// HistogramOpts is the options for a Prometheus histogram.
type HistogramOpts struct {
	Namespace string   `validate:"required"`
	Name      string   `validate:"required"`
	Help      string   `validate:"required"`
	Labels    []string `validate:"required"`
}

// NewHistogramWithLabels creates a Prometheus histogram with labels based on the
// provided HistogramOpts.
// A histogram samples observations (usually things like request durations or
// response sizes) and counts them in configurable buckets. It also provides
// a sum of all observed values.
// TODO: make Objectives configurable.
func NewHistogramWithLabels(opts HistogramOpts) (*prometheus.HistogramVec, error) {
	if err := Validate(opts); err != nil {
		return nil, err
	}

	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: opts.Namespace,
		Help:      opts.Help,
		Name:      opts.Name,
	}, opts.Labels), nil
}
