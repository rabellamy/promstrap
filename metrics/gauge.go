package promstrap

import "github.com/prometheus/client_golang/prometheus"

// GaugeOpts is the options for Prometheus a Prometheus Gauge.
type GaugeOpts struct {
	Namespace string   `validate:"required"`
	Name      string   `validate:"required"`
	Help      string   `validate:"required"`
	Labels    []string `validate:"required"`
}

// NewGaugeWithLabels creates a Prometheus Gauge with labels based on the
// provided GaugeOpts.
// A gauge is a metric that represents a single numerical value that can
// arbitrarily go up and down. Gauges track values that can change over time,
// such as the amount of memory used, the number of requests in progress,
// or the temperature of a device.
func NewGaugeWithLabels(opts GaugeOpts) (*prometheus.GaugeVec, error) {
	if err := Validate(opts); err != nil {
		return nil, err
	}

	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: opts.Namespace,
		Name:      opts.Name,
		Help:      opts.Help,
	}, opts.Labels), nil
}
