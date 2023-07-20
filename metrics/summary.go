package promstrap

import "github.com/prometheus/client_golang/prometheus"

// SummaryOpts is the options for a Prometheus summary.
type SummaryOpts struct {
	Namespace string   `validate:"required"`
	Name      string   `validate:"required"`
	Help      string   `validate:"required"`
	Labels    []string `validate:"required"`
	// Objectives defines the quantile rank estimates with their respective
	// absolute error.
	Objectives map[float64]float64
}

// NewSummaryWithLabels creates a Prometheus summary with labels based on the
// provided SummaryOpts.
// A summary samples observations (usually things like request durations and
// response sizes). While it also provides a total count of observations and a
// sum of all observed values, it calculates configurable quantiles over
// a sliding time window.
func NewSummaryWithLabels(opts SummaryOpts) (*prometheus.SummaryVec, error) {
	if err := Validate(opts); err != nil {
		return nil, err
	}

	pOpts := prometheus.SummaryOpts{
		Namespace: opts.Namespace,
		Help:      opts.Help,
		Name:      opts.Name,
	}

	if opts.Objectives != nil {
		pOpts.Objectives = opts.Objectives
	}

	return prometheus.NewSummaryVec(pOpts, opts.Labels), nil
}
