package promstrap

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestNewCounterWithLabels(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		opts    MetricsOpts
		want    *prometheus.CounterVec
		wantErr bool
	}{
		"all good": {
			opts: MetricsOpts{
				Namespace: "the_namespace",
				Name:      "the_name_total",
				Help:      "Some help text",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want: prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: "the_namespace",
				Name:      "the_name_total",
				Help:      "Some help text",
			}, []string{"yo", "bro", "flow"}),
			wantErr: false,
		},
		"no Namepace": {
			opts: MetricsOpts{
				Namespace: "",
				Name:      "the_name",
				Help:      "Some help text",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Name": {
			opts: MetricsOpts{
				Namespace: "the_namespace",
				Name:      "",
				Help:      "Some help text",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Help": {
			opts: MetricsOpts{
				Namespace: "the_namespace",
				Name:      "the_name",
				Help:      "",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Labels": {
			opts: MetricsOpts{
				Namespace: "the_namespace",
				Name:      "the_name",
				Help:      "Some help text",
				Labels:    nil,
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
			got, err := NewCounterWithLabels(opts)
			if (err != nil) != wantErr {
				t.Errorf("NewCounterWithLabels() error = %v, wantErr %v", err, wantErr)

				return
			}
			if want == nil && got == nil {
				return
			}
			assert.EqualExportedValues(t, *want, *got)
		})
	}
}

func TestNewGaugeWithLabels(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		opts    MetricsOpts
		want    *prometheus.GaugeVec
		wantErr bool
	}{
		"all good": {
			opts: MetricsOpts{
				Namespace: "the_namespace",
				Name:      "the_name",
				Help:      "Some help text",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: "the_namespace",
				Name:      "the_name",
				Help:      "Some help text",
			}, []string{"yo", "bro", "flow"}),
			wantErr: false,
		},
		"no Namepace": {
			opts: MetricsOpts{
				Namespace: "",
				Name:      "the_name",
				Help:      "Some help text",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Name": {
			opts: MetricsOpts{
				Namespace: "the_namespace",
				Name:      "",
				Help:      "Some help text",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Help": {
			opts: MetricsOpts{
				Namespace: "the_namespace",
				Name:      "the_name",
				Help:      "",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Labels": {
			opts: MetricsOpts{
				Namespace: "the_namespace",
				Name:      "the_name",
				Help:      "Some help text",
				Labels:    nil,
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
			got, err := NewGaugeWithLabels(opts)
			if (err != nil) != wantErr {
				t.Errorf("NewGaugeWithLabels() error = %v, wantErr %v", err, wantErr)

				return
			}
			if want == nil && got == nil {
				return
			}
			assert.EqualExportedValues(t, *want, *got)
		})
	}
}

func TestNewHistogramWithLabels(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		opts    MetricsOpts
		want    *prometheus.HistogramVec
		wantErr bool
	}{
		"all good": {
			opts: MetricsOpts{
				Namespace: "the_namespace",
				Name:      "the_name",
				Help:      "Some help text",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want: prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Namespace: "the_namespace",
				Name:      "the_name",
				Help:      "Some help text",
			}, []string{"yo", "bro", "flow"}),
			wantErr: false,
		},
		"no Namepace": {
			opts: MetricsOpts{
				Namespace: "",
				Name:      "the_name",
				Help:      "Some help text",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Name": {
			opts: MetricsOpts{
				Namespace: "the_namespace",
				Name:      "",
				Help:      "Some help text",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Help": {
			opts: MetricsOpts{
				Namespace: "the_namespace",
				Name:      "the_name",
				Help:      "",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Labels": {
			opts: MetricsOpts{
				Namespace: "the_namespace",
				Name:      "the_name",
				Help:      "Some help text",
				Labels:    nil,
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
			got, err := NewHistogramWithLabels(opts)
			if (err != nil) != wantErr {
				t.Errorf("NewHistogramWithLabels() error = %v, wantErr %v", err, wantErr)

				return
			}
			if want == nil && got == nil {
				return
			}
			assert.EqualExportedValues(t, *want, *got)
		})
	}
}

func TestNewSummaryWithLabels(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		opts    MetricsOpts
		want    *prometheus.SummaryVec
		wantErr bool
	}{
		"all good": {
			opts: MetricsOpts{
				Namespace: "the_namespace",
				Name:      "the_name",
				Help:      "Some help text",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want: prometheus.NewSummaryVec(prometheus.SummaryOpts{
				Namespace: "the_namespace",
				Name:      "the_name",
				Help:      "Some help text",
			}, []string{"yo", "bro", "flow"}),
			wantErr: false,
		},
		"no Namepace": {
			opts: MetricsOpts{
				Namespace: "",
				Name:      "the_name",
				Help:      "Some help text",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Name": {
			opts: MetricsOpts{
				Namespace: "the_namespace",
				Name:      "",
				Help:      "Some help text",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Help": {
			opts: MetricsOpts{
				Namespace: "the_namespace",
				Name:      "the_name",
				Help:      "",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Labels": {
			opts: MetricsOpts{
				Namespace: "the_namespace",
				Name:      "the_name",
				Help:      "Some help text",
				Labels:    nil,
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
			got, err := NewSummaryWithLabels(opts)
			if (err != nil) != wantErr {
				t.Errorf("NewSummaryWithLabels() error = %v, wantErr %v", err, wantErr)

				return
			}
			if want == nil && got == nil {
				return
			}
			assert.EqualExportedValues(t, *want, *got)
		})
	}
}

func TestRegisterMetrics(t *testing.T) {
	t.Parallel()

	testCounter, err := NewCounterWithLabels(MetricsOpts{
		Namespace: "test_counter",
		Name:      "foo",
		Help:      "bar",
		Labels:    []string{"baz"},
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	testComplexMetricCounter, err := NewCounterWithLabels(MetricsOpts{
		Namespace: "test_complex_counter",
		Name:      "foo",
		Help:      "bar",
		Labels:    []string{"baz"},
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	testGauge, err := NewGaugeWithLabels(MetricsOpts{
		Namespace: "test_gauge",
		Name:      "foo",
		Help:      "bar",
		Labels:    []string{"baz"},
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	testHistogram, err := NewHistogramWithLabels(MetricsOpts{
		Namespace: "test_histogram",
		Name:      "foo",
		Help:      "bar",
		Labels:    []string{"baz"},
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	testSummary, err := NewSummaryWithLabels(MetricsOpts{
		Namespace: "test_summary",
		Name:      "foo",
		Help:      "bar",
		Labels:    []string{"baz"},
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	type testValidComplexMetric struct {
		Foo *prometheus.CounterVec
	}

	type testInvalidComplexMetric struct{}

	validComplexMetric := testValidComplexMetric{
		Foo: testComplexMetricCounter,
	}

	invalidComplexMetric := testInvalidComplexMetric{}

	tests := map[string]struct {
		metrics interface{}
		wantErr bool
	}{
		"valid type - Counter": {
			metrics: *testCounter,
			wantErr: false,
		},
		"valid type - Gauge": {
			metrics: *testGauge,
			wantErr: false,
		},
		"valid type - Histogram": {
			metrics: *testHistogram,
			wantErr: false,
		},
		"valid type - Summary": {
			metrics: *testSummary,
			wantErr: false,
		},
		"valid complex": {
			metrics: validComplexMetric,
			wantErr: false,
		},
		"invalid complex": {
			metrics: invalidComplexMetric,
			wantErr: true,
		},
		"string": {
			metrics: "not good",
			wantErr: true,
		},
		"number": {
			metrics: 123,
			wantErr: true,
		},
		"nil": {
			metrics: nil,
			wantErr: true,
		},
	}

	for name, tt := range tests {
		metrics := tt.metrics
		wantErr := tt.wantErr

		t.Run(name, func(t *testing.T) {
			if err := RegisterMetrics(metrics); (err != nil) != wantErr {
				t.Errorf("RegisterMetrics() error = %v, wantErr %v", err, wantErr)
			}
		})
	}
}