package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestRegisterCollectors(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		collectors []prometheus.Collector
		wantErr    bool
	}{
		"single collector": {
			collectors: []prometheus.Collector{
				prometheus.NewCounter(prometheus.CounterOpts{
					Name: "test_counter",
					Help: "Test counter",
				}),
			},
			wantErr: false,
		},
		"multiple collectors": {
			collectors: []prometheus.Collector{
				prometheus.NewCounter(prometheus.CounterOpts{
					Name: "test_counter_1",
					Help: "Test counter 1",
				}),
				prometheus.NewGauge(prometheus.GaugeOpts{
					Name: "test_gauge_1",
					Help: "Test gauge 1",
				}),
			},
			wantErr: false,
		},
		"duplicate collectors": {
			collectors: []prometheus.Collector{
				prometheus.NewCounter(prometheus.CounterOpts{
					Name: "test_counter_dupe",
					Help: "Test counter duplicate",
				}),
				prometheus.NewCounter(prometheus.CounterOpts{
					Name: "test_counter_dupe",
					Help: "Test counter duplicate",
				}),
			},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		collectors := tt.collectors
		wantErr := tt.wantErr

		t.Run(name, func(t *testing.T) {
			err := RegisterCollectors(collectors...)
			if (err != nil) != wantErr {
				t.Errorf("RegisterCollectors() error = %v, wantErr %v", err, wantErr)
			}
		})
	}

	// Clean up the test metrics from the default registry
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
}
