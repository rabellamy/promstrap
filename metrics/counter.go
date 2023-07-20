package promstrap

import "github.com/prometheus/client_golang/prometheus"

// CounterOpts is the options for a Prometheus counter.
type CounterOpts struct {
	Namespace string   `validate:"required"`
	Name      string   `validate:"required"`
	Help      string   `validate:"required"`
	Labels    []string `validate:"required"`
}

// NewCounterWithLabels creates a Prometheus counter with labels based on the
// provided CounterOpts.
// A counter is a cumulative metric that represents a single monotonically
// increasing counter whose value can only increase or be reset to zero on restart.
// Counters are for tracking cumulative totals over time, like the total number
// of HTTP requests or the number of errors.
func NewCounterWithLabels(opts CounterOpts) (*prometheus.CounterVec, error) {
	if err := Validate(opts); err != nil {
		return nil, err
	}

	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: opts.Namespace,
		Name:      opts.Name,
		Help:      opts.Help,
	}, opts.Labels), nil
}
