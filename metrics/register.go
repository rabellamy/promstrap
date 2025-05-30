package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// RegisterCollectors registers the provided Prometheus collectors with the default registry.
// It returns the first error encountered, if any.
func RegisterCollectors(collectors ...prometheus.Collector) error {
	for _, collector := range collectors {
		if err := prometheus.DefaultRegisterer.Register(collector); err != nil {
			return err
		}
	}
	return nil
}
