package strategy

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestNewFourGoldenSignals(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		opts    FourGoldenSignalsOpts
		want    *FourGoldenSignals
		wantErr bool
	}{
		"all good": {
			opts: FourGoldenSignalsOpts{
				Namespace: "test",
				LatencyOpt: FGSLatencyOpt{
					LatencyName:   "http_request_latency_seconds",
					LatencyType:   "http",
					LatencyHelp:   "Request latency in seconds",
					LatencyLabels: []string{"method", "path"},
					Buckets:       []float64{.005, .01, .025, .05, .1, .25, .5, 1},
					Objectives:    map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
				},
				TrafficOpt: FGSTrafficOpt{
					TrafficName:   "http_requests_total",
					TrafficType:   "http",
					TrafficHelp:   "Total number of HTTP requests",
					TrafficLabels: []string{"method", "path", "status"},
				},
				ErrorsOpt: FGSErrorsOpt{
					ErrorHelp:   "Number of errors",
					ErrorLabels: []string{"type"},
				},
				SaturationOpt: FGSSaturationOpt{
					SaturationName:   "memory_heap_saturation_bytes",
					SaturationHelp:   "Memory heap usage in bytes",
					SaturationLabels: []string{"gc_type"},
				},
			},
			want: &FourGoldenSignals{
				Latency: &Distribution{
					Histogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
						Namespace: "test",
						Name:      "http_request_latency_seconds_hist",
						Help:      "Request latency in seconds",
						Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1},
					}, []string{"method", "path"}),
					Summary: prometheus.NewSummaryVec(prometheus.SummaryOpts{
						Namespace:  "test",
						Name:       "http_request_latency_seconds_sum",
						Help:      "Request latency in seconds",
						Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
					}, []string{"method", "path"}),
				},
				Traffic: prometheus.NewCounterVec(prometheus.CounterOpts{
					Namespace: "test",
					Name:      "http_requests_total",
					Help:      "Total number of HTTP requests",
				}, []string{"method", "path", "status"}),
				Errors: prometheus.NewCounterVec(prometheus.CounterOpts{
					Namespace: "test",
					Name:      "errors_total",
					Help:      "Number of errors",
				}, []string{"type"}),
				Saturation: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: "test",
					Name:      "memory_heap_saturation_bytes",
					Help:      "Memory heap usage in bytes",
				}, []string{"gc_type"}),
			},
			wantErr: false,
		},
		"missing namespace": {
			opts: FourGoldenSignalsOpts{
				LatencyOpt: FGSLatencyOpt{
					LatencyName:   "http_request_latency_seconds",
					LatencyType:   "http",
					LatencyHelp:   "Request latency in seconds",
					LatencyLabels: []string{"method", "path"},
				},
				TrafficOpt: FGSTrafficOpt{
					TrafficName:   "http_requests_total",
					TrafficType:   "http",
					TrafficHelp:   "Total number of HTTP requests",
					TrafficLabels: []string{"method", "path", "status"},
				},
				ErrorsOpt: FGSErrorsOpt{
					ErrorHelp:   "Number of errors",
					ErrorLabels: []string{"type"},
				},
				SaturationOpt: FGSSaturationOpt{
					SaturationName:   "memory_heap_saturation_bytes",
					SaturationHelp:   "Memory heap usage in bytes",
					SaturationLabels: []string{"gc_type"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		"missing latency name": {
			opts: FourGoldenSignalsOpts{
				Namespace: "test",
				LatencyOpt: FGSLatencyOpt{
					LatencyType:   "http",
					LatencyHelp:   "Request latency in seconds",
					LatencyLabels: []string{"method", "path"},
				},
				TrafficOpt: FGSTrafficOpt{
					TrafficName:   "http_requests_total",
					TrafficType:   "http",
					TrafficHelp:   "Total number of HTTP requests",
					TrafficLabels: []string{"method", "path", "status"},
				},
				ErrorsOpt: FGSErrorsOpt{
					ErrorHelp:   "Number of errors",
					ErrorLabels: []string{"type"},
				},
				SaturationOpt: FGSSaturationOpt{
					SaturationName:   "memory_heap_saturation_bytes",
					SaturationHelp:   "Memory heap usage in bytes",
					SaturationLabels: []string{"gc_type"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		"missing traffic name": {
			opts: FourGoldenSignalsOpts{
				Namespace: "test",
				LatencyOpt: FGSLatencyOpt{
					LatencyName:   "http_request_latency_seconds",
					LatencyType:   "http",
					LatencyHelp:   "Request latency in seconds",
					LatencyLabels: []string{"method", "path"},
				},
				TrafficOpt: FGSTrafficOpt{
					TrafficType:   "http",
					TrafficHelp:   "Total number of HTTP requests",
					TrafficLabels: []string{"method", "path", "status"},
				},
				ErrorsOpt: FGSErrorsOpt{
					ErrorHelp:   "Number of errors",
					ErrorLabels: []string{"type"},
				},
				SaturationOpt: FGSSaturationOpt{
					SaturationName:   "memory_heap_saturation_bytes",
					SaturationHelp:   "Memory heap usage in bytes",
					SaturationLabels: []string{"gc_type"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		"missing saturation name": {
			opts: FourGoldenSignalsOpts{
				Namespace: "test",
				LatencyOpt: FGSLatencyOpt{
					LatencyName:   "http_request_latency_seconds",
					LatencyType:   "http",
					LatencyHelp:   "Request latency in seconds",
					LatencyLabels: []string{"method", "path"},
				},
				TrafficOpt: FGSTrafficOpt{
					TrafficName:   "http_requests_total",
					TrafficType:   "http",
					TrafficHelp:   "Total number of HTTP requests",
					TrafficLabels: []string{"method", "path", "status"},
				},
				ErrorsOpt: FGSErrorsOpt{
					ErrorHelp:   "Number of errors",
					ErrorLabels: []string{"type"},
				},
				SaturationOpt: FGSSaturationOpt{
					SaturationHelp:   "Memory heap usage in bytes",
					SaturationLabels: []string{"gc_type"},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for name, tt := range tests {
		opts := tt.opts
		wantErr := tt.wantErr
		want := tt.want

		t.Run(name, func(t *testing.T) {
			got, err := NewFourGoldenSignals(opts)
			if (err != nil) != wantErr {
				t.Errorf("NewFourGoldenSignals() error = %v, wantErr %v", err, wantErr)
				return
			}
			if want == nil && got == nil {
				return
			}
			assert.EqualExportedValues(t, *want, *got)
		})
	}
}

func TestFGSMetricNames(t *testing.T) {
	t.Parallel()

	fgs, err := NewFourGoldenSignals(FourGoldenSignalsOpts{
		Namespace: "test",
		LatencyOpt: FGSLatencyOpt{
			LatencyName:   "http_request_latency_seconds",
			LatencyType:   "http",
			LatencyHelp:   "Request latency in seconds",
			LatencyLabels: []string{"method"},
		},
		TrafficOpt: FGSTrafficOpt{
			TrafficName:   "http_requests_total",
			TrafficType:   "http",
			TrafficHelp:   "Total number of requests",
			TrafficLabels: []string{"method"},
		},
		ErrorsOpt: FGSErrorsOpt{
			ErrorName:   "custom_errors",
			ErrorHelp:   "Number of errors",
			ErrorLabels: []string{"type"},
		},
		SaturationOpt: FGSSaturationOpt{
			SaturationName:   "memory_heap_saturation_bytes",
			SaturationHelp:   "Memory heap usage",
			SaturationLabels: []string{"type"},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, "http_request_latency_seconds", fgs.LatencyMetricName())
	assert.Equal(t, "http_requests_total", fgs.TrafficMetricName())
	assert.Equal(t, "custom_errors", fgs.ErrorMetricName())
	assert.Equal(t, "memory_heap_saturation_bytes", fgs.SaturationMetricName())

	// Test with default error name
	fgsDefaultError, err := NewFourGoldenSignals(FourGoldenSignalsOpts{
		Namespace: "test",
		LatencyOpt: FGSLatencyOpt{
			LatencyName:   "http_request_latency_seconds",
			LatencyType:   "http",
			LatencyHelp:   "Request latency in seconds",
			LatencyLabels: []string{"method"},
		},
		TrafficOpt: FGSTrafficOpt{
			TrafficName:   "http_requests_total",
			TrafficType:   "http",
			TrafficHelp:   "Total number of requests",
			TrafficLabels: []string{"method"},
		},
		ErrorsOpt: FGSErrorsOpt{
			ErrorHelp:   "Number of errors",
			ErrorLabels: []string{"type"},
		},
		SaturationOpt: FGSSaturationOpt{
			SaturationName:   "memory_heap_saturation_bytes",
			SaturationHelp:   "Memory heap usage",
			SaturationLabels: []string{"type"},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, "errors_total", fgsDefaultError.ErrorMetricName())
}
