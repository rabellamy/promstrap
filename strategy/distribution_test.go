package strategy

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestNewDistribution(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		opts    DistributionOpts
		want    *Distribution
		wantErr bool
	}{
		"all good": {
			opts: DistributionOpts{
				Namespace:  "foobar",
				Name:       "foo",
				Help:       "bar",
				Labels:     []string{"baz", "qux", "quux"},
				Buckets:    []float64{.5, 1.5, 2.0},
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			want: &Distribution{
				Histogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
					Namespace: "foobar",
					Name:      "foo_hist",
					Help:      "bar",
					Buckets:   []float64{.5, 1.5, 2.0},
				}, []string{"baz", "qux", "quux"}),
				Summary: prometheus.NewSummaryVec(prometheus.SummaryOpts{
					Namespace:  "foobar",
					Name:       "foo_sum",
					Help:       "bar",
					Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
				}, []string{"baz", "qux", "quux"}),
			},
			wantErr: false,
		},
	}
	for name, tt := range tests {
		opts := tt.opts
		wantErr := tt.wantErr
		want := tt.want

		t.Run(name, func(t *testing.T) {
			got, err := NewDistribution(opts)
			if (err != nil) != wantErr {
				t.Errorf("NewDistribution error = %v, wantErr %v", err, wantErr)

				return
			}
			if want == nil && got == nil {
				return
			}
			assert.EqualExportedValues(t, *want, *got)
		})
	}
}

func TestDistributionCollectors(t *testing.T) {
	t.Parallel()

	dist, err := NewDistribution(DistributionOpts{
		Namespace:  "test",
		Name:       "request_duration_seconds",
		Help:       "Request duration in seconds",
		Labels:     []string{"method", "path"},
		Buckets:    []float64{.005, .01, .025, .05, .1, .25, .5, 1},
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
	assert.NoError(t, err)

	collectors := dist.Collectors()

	assert.Equal(t, 2, len(collectors))

	for _, collector := range collectors {
		switch collector.(type) {
		case *prometheus.HistogramVec, *prometheus.SummaryVec:
		default:
			t.Errorf("Unexpected collector type: %T", collector)
		}
	}
}
