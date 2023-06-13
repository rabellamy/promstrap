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
				Namespace:        "foobar",
				LatencyName:      "foo",
				LatencyHelp:      "bar",
				LatencyLabels:    []string{"baz", "qux", "quux"},
				TrafficName:      "corge_total",
				TrafficHelp:      "grault",
				TrafficLabels:    []string{"garply", "waldo", "fred"},
				SaturationName:   "plugh",
				SaturationHelp:   "xyzzy,",
				SaturationLabels: []string{"thud"},
			},
			want: &FourGoldenSignals{
				Latency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
					Namespace: "foobar",
					Name:      "foo",
					Help:      "bar",
				}, []string{"baz", "qux", "quux"}),
				Traffic: prometheus.NewCounterVec(prometheus.CounterOpts{
					Namespace: "foobar",
					Name:      "corge_total",
					Help:      "grault",
				}, []string{"garply", "waldo", "fred"}),
				Errors: prometheus.NewCounterVec(prometheus.CounterOpts{
					Namespace: "foobar",
					Name:      "errors_total",
					Help:      "Number of errors",
				}, []string{"error"}),
				Saturation: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: "foobar",
					Name:      "plugh",
					Help:      "xyzzy",
				}, []string{"thud"}),
			},
			wantErr: false,
		},
		"missing Namespace": {
			opts: FourGoldenSignalsOpts{
				Namespace:        "",
				LatencyName:      "foo",
				LatencyHelp:      "bar",
				LatencyLabels:    []string{"baz", "qux", "quux"},
				TrafficName:      "corge",
				TrafficHelp:      "grault",
				TrafficLabels:    []string{"garply", "waldo", "fred"},
				SaturationName:   "plugh",
				SaturationHelp:   "xyzzy",
				SaturationLabels: []string{"thud"},
			},
			want:    nil,
			wantErr: true,
		},
		"missing LatencyName": {
			opts: FourGoldenSignalsOpts{
				Namespace:        "foobar",
				LatencyName:      "",
				LatencyHelp:      "bar",
				LatencyLabels:    []string{"baz", "qux", "quux"},
				TrafficName:      "corge",
				TrafficHelp:      "grault",
				TrafficLabels:    []string{"garply", "waldo", "fred"},
				SaturationName:   "plugh",
				SaturationHelp:   "xyzzy",
				SaturationLabels: []string{"thud"},
			},
			want:    nil,
			wantErr: true,
		},
		"missing LatencyHelp": {
			opts: FourGoldenSignalsOpts{
				Namespace:        "foobar",
				LatencyName:      "foo",
				LatencyHelp:      "",
				LatencyLabels:    []string{"baz", "qux", "quux"},
				TrafficName:      "corge",
				TrafficHelp:      "grault",
				TrafficLabels:    []string{"garply", "waldo", "fred"},
				SaturationName:   "plugh",
				SaturationHelp:   "xyzzy",
				SaturationLabels: []string{"thud"},
			},
			want:    nil,
			wantErr: true,
		},
		"missing LatencyLabels": {
			opts: FourGoldenSignalsOpts{
				Namespace:        "foobar",
				LatencyName:      "foo",
				LatencyHelp:      "bar",
				LatencyLabels:    nil,
				TrafficName:      "corge",
				TrafficHelp:      "grault",
				TrafficLabels:    []string{"garply", "waldo", "fred"},
				SaturationName:   "plugh",
				SaturationHelp:   "xyzzy",
				SaturationLabels: []string{"thud"},
			},
			want:    nil,
			wantErr: true,
		},
		"missing TrafficName": {
			opts: FourGoldenSignalsOpts{
				Namespace:        "foobar",
				LatencyName:      "foo",
				LatencyHelp:      "bar",
				LatencyLabels:    []string{"baz", "qux", "quux"},
				TrafficName:      "",
				TrafficHelp:      "grault",
				TrafficLabels:    []string{"garply", "waldo", "fred"},
				SaturationName:   "plugh",
				SaturationHelp:   "xyzzy",
				SaturationLabels: []string{"thud"},
			},
			want:    nil,
			wantErr: true,
		},
		"missing TrafficHelp": {
			opts: FourGoldenSignalsOpts{
				Namespace:        "foobar",
				LatencyName:      "foo",
				LatencyHelp:      "bar",
				LatencyLabels:    []string{"baz", "qux", "quux"},
				TrafficName:      "corge",
				TrafficHelp:      "",
				TrafficLabels:    []string{"garply", "waldo", "fred"},
				SaturationName:   "plugh",
				SaturationHelp:   "xyzzy",
				SaturationLabels: []string{"thud"},
			},
			want:    nil,
			wantErr: true,
		},
		"missing TrafficLabels": {
			opts: FourGoldenSignalsOpts{
				Namespace:        "foobar",
				LatencyName:      "foo",
				LatencyHelp:      "bar",
				LatencyLabels:    []string{"baz", "qux", "quux"},
				TrafficName:      "corge",
				TrafficHelp:      "grault",
				TrafficLabels:    nil,
				SaturationName:   "plugh",
				SaturationHelp:   "xyzzy",
				SaturationLabels: []string{"thud"},
			},
			want:    nil,
			wantErr: true,
		},
		"missing SaturationName": {
			opts: FourGoldenSignalsOpts{
				Namespace:        "foobar",
				LatencyName:      "foo",
				LatencyHelp:      "bar",
				LatencyLabels:    []string{"baz", "qux", "quux"},
				TrafficName:      "corge",
				TrafficHelp:      "grault",
				TrafficLabels:    []string{"garply", "waldo", "fred"},
				SaturationName:   "",
				SaturationHelp:   "xyzzy",
				SaturationLabels: []string{"thud"},
			},
			want:    nil,
			wantErr: true,
		},
		"missing SaturationHelp": {
			opts: FourGoldenSignalsOpts{
				Namespace:        "foobar",
				LatencyName:      "foo",
				LatencyHelp:      "bar",
				LatencyLabels:    []string{"baz", "qux", "quux"},
				TrafficName:      "corge",
				TrafficHelp:      "grault",
				TrafficLabels:    []string{"garply", "waldo", "fred"},
				SaturationName:   "plugh",
				SaturationHelp:   "",
				SaturationLabels: []string{"thud"},
			},
			want:    nil,
			wantErr: true,
		},
		"missing SaturationLabels": {
			opts: FourGoldenSignalsOpts{
				Namespace:        "foobar",
				LatencyName:      "foo",
				LatencyHelp:      "bar",
				LatencyLabels:    []string{"baz", "qux", "quux"},
				TrafficName:      "corge",
				TrafficHelp:      "grault",
				TrafficLabels:    []string{"garply", "waldo", "fred"},
				SaturationName:   "plugh",
				SaturationHelp:   "xyzzy",
				SaturationLabels: nil,
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
