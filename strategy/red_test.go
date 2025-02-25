package strategy

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestNewRED(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		opts    REDOpts
		want    *RED
		wantErr bool
	}{
		"default names": {
			opts: REDOpts{
				Namespace: "bar",
				RequestsOpt: REDRequestsOpt{
					RequestType:   "foo",
					RequestLabels: []string{"jazz", "wiz", "fizz"},
				},
				ErrorsOpt: REDErrorsOpt{
					ErrorLabels: []string{"error"},
				},
				DurationOpt: REDDurationOpt{
					DurationLabels: []string{"cuz", "buzz"},
					Buckets:        []float64{.5, 1.5, 2.0},
					Objectives:     map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
				},
			},
			want: &RED{
				Requests: prometheus.NewCounterVec(prometheus.CounterOpts{
					Namespace: "bar",
					Name:      "foo_requests_total",
					Help:      "Number of requests",
				}, []string{"jazz", "wiz", "fizz"}),
				Errors: prometheus.NewCounterVec(prometheus.CounterOpts{
					Namespace: "bar",
					Name:      "errors_total",
					Help:      "Number of errors",
				}, []string{"error"}),
				Duration: &Distribution{
					Histogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
						Namespace: "bar",
						Name:      "foo_request_duration_seconds_hist",
						Help:      "Duration of request in seconds",
						Buckets:   []float64{.5, 1.5, 2.0},
					}, []string{"cuz", "buzz"}),
					Summary: prometheus.NewSummaryVec(prometheus.SummaryOpts{
						Namespace:  "bar",
						Name:       "foo_request_duration_seconds_sum",
						Help:       "Duration of request in seconds",
						Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
					}, []string{"cuz", "buzz"}),
				},
			},
			wantErr: false,
		},
		"custom names": {
			opts: REDOpts{
				Namespace: "bar",
				RequestsOpt: REDRequestsOpt{
					RequestType:   "foo",
					RequestName:   "custom_requests",
					RequestLabels: []string{"jazz", "wiz", "fizz"},
				},
				ErrorsOpt: REDErrorsOpt{
					ErrorName:   "custom_errors",
					ErrorLabels: []string{"error"},
				},
				DurationOpt: REDDurationOpt{
					DurationName:   "custom_duration",
					DurationLabels: []string{"cuz", "buzz"},
					Buckets:        []float64{.5, 1.5, 2.0},
					Objectives:     map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
				},
			},
			want: &RED{
				Requests: prometheus.NewCounterVec(prometheus.CounterOpts{
					Namespace: "bar",
					Name:      "custom_requests",
					Help:      "Number of requests",
				}, []string{"jazz", "wiz", "fizz"}),
				Errors: prometheus.NewCounterVec(prometheus.CounterOpts{
					Namespace: "bar",
					Name:      "custom_errors",
					Help:      "Number of errors",
				}, []string{"error"}),
				Duration: &Distribution{
					Histogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
						Namespace: "bar",
						Name:      "custom_duration_hist",
						Help:      "Duration of request in seconds",
						Buckets:   []float64{.5, 1.5, 2.0},
					}, []string{"cuz", "buzz"}),
					Summary: prometheus.NewSummaryVec(prometheus.SummaryOpts{
						Namespace:  "bar",
						Name:       "custom_duration_sum",
						Help:       "Duration of request in seconds",
						Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
					}, []string{"cuz", "buzz"}),
				},
			},
			wantErr: false,
		},
		"missing RequestType": {
			opts: REDOpts{
				Namespace: "pineapple",
				RequestsOpt: REDRequestsOpt{
					RequestLabels: []string{"jazz", "wiz", "fizz"},
				},
				ErrorsOpt: REDErrorsOpt{
					ErrorLabels: []string{"error"},
				},
				DurationOpt: REDDurationOpt{
					DurationLabels: []string{"cuz", "buzz"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		"missing Namespace": {
			opts: REDOpts{
				RequestsOpt: REDRequestsOpt{
					RequestType:   "radish",
					RequestLabels: []string{"jazz", "wiz", "fizz"},
				},
				ErrorsOpt: REDErrorsOpt{
					ErrorLabels: []string{"error"},
				},
				DurationOpt: REDDurationOpt{
					DurationLabels: []string{"cuz", "buzz"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		"missing RequestLabels": {
			opts: REDOpts{
				Namespace: "pickle",
				RequestsOpt: REDRequestsOpt{
					RequestType:   "dill",
					RequestLabels: nil,
				},
				ErrorsOpt: REDErrorsOpt{
					ErrorLabels: []string{"error"},
				},
				DurationOpt: REDDurationOpt{
					DurationLabels: []string{"cuz", "buzz"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		"missing DurationLabels": {
			opts: REDOpts{
				Namespace: "house",
				RequestsOpt: REDRequestsOpt{
					RequestType:   "tiny",
					RequestLabels: []string{"jazz", "wiz", "fizz"},
				},
				ErrorsOpt: REDErrorsOpt{
					ErrorLabels: []string{"error"},
				},
				DurationOpt: REDDurationOpt{
					DurationLabels: nil,
				},
			},
			want:    nil,
			wantErr: true,
		},
		"invalid buckets": {
			opts: REDOpts{
				Namespace: "bar",
				RequestsOpt: REDRequestsOpt{
					RequestType:   "foo",
					RequestLabels: []string{"label"},
				},
				DurationOpt: REDDurationOpt{
					DurationLabels: []string{"label"},
					Buckets:        []float64{2.0, 1.0}, // Invalid order
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
			got, err := NewRED(opts)
			if (err != nil) != wantErr {
				t.Errorf("NewRED() error = %v, wantErr %v", err, wantErr)

				return
			}
			if want == nil && got == nil {
				return
			}
			assert.EqualExportedValues(t, *want, *got)
		})
	}
}

func TestREDMetricNames(t *testing.T) {
	t.Parallel()

	// Test for custom names
	red, err := NewRED(REDOpts{
		Namespace: "test",
		RequestsOpt: REDRequestsOpt{
			RequestType:   "http",
			RequestName:   "custom_requests",
			RequestLabels: []string{"method"},
		},
		ErrorsOpt: REDErrorsOpt{
			ErrorName:   "custom_errors",
			ErrorLabels: []string{"type"},
		},
		DurationOpt: REDDurationOpt{
			DurationName:   "custom_duration",
			DurationLabels: []string{"endpoint"},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, "custom_requests", red.RequestMetricName())
	assert.Equal(t, "custom_errors", red.ErrorMetricName())
	assert.Equal(t, "custom_duration", red.DurationMetricName())

	// Test for default names
	redDefault, err := NewRED(REDOpts{
		Namespace: "test",
		RequestsOpt: REDRequestsOpt{
			RequestType:   "http",
			RequestLabels: []string{"method"},
		},
		ErrorsOpt: REDErrorsOpt{
			ErrorLabels: []string{"type"},
		},
		DurationOpt: REDDurationOpt{
			DurationLabels: []string{"endpoint"},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, "http_requests_total", redDefault.RequestMetricName())
	assert.Equal(t, "errors_total", redDefault.ErrorMetricName())
	assert.Equal(t, "http_request_duration_seconds", redDefault.DurationMetricName())
}
