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
		"all good": {
			opts: REDOpts{
				RequestType:    "foo",
				Namespace:      "bar",
				RequestLabels:  []string{"jazz", "wiz", "fizz"},
				DurationLabels: []string{"cuz", "buzz"},
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
				Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
					Namespace: "bar",
					Name:      "foo_request_duration_seconds",
					Help:      "Duration of request in seconds",
				}, []string{"cuz", "buzz"}),
			},
			wantErr: false,
		},
		"missing RequestType": {
			opts: REDOpts{
				RequestType:    "",
				Namespace:      "pineapple",
				RequestLabels:  []string{"jazz", "wiz", "fizz"},
				DurationLabels: []string{"cuz", "buzz"},
			},
			want:    nil,
			wantErr: true,
		},
		"missing Namespace": {
			opts: REDOpts{
				RequestType:    "radish",
				Namespace:      "",
				RequestLabels:  []string{"jazz", "wiz", "fizz"},
				DurationLabels: []string{"cuz", "buzz"},
			},
			want:    nil,
			wantErr: true,
		},
		"missing RequestLabels": {
			opts: REDOpts{
				RequestType:    "dill",
				Namespace:      "pickle",
				RequestLabels:  nil,
				DurationLabels: []string{"cuz", "buzz"},
			},
			want:    nil,
			wantErr: true,
		},
		"missing DurationLabels": {
			opts: REDOpts{
				RequestType:    "tiny",
				Namespace:      "house",
				RequestLabels:  []string{"jazz", "wiz", "fizz"},
				DurationLabels: nil,
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
