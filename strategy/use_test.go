package strategy

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestNewUSE(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		opts    USEOpts
		want    *USE
		wantErr bool
	}{
		"all good": {
			opts: USEOpts{
				Namespace: "foo",
				UtilizationOpt: USEUtilizationOpt{
					UtilizationName:   "qux",
					UtilizationHelp:   "quux",
					UtilizationLabels: []string{"fred", "plugh", "xyzzy"},
				},
				SaturationOpt: USESaturationOpt{
					SaturationName:   "bar",
					SaturationHelp:   "baz",
					SaturationLabels: []string{"grault", "garply", "waldo"},
				},
				ErrorsOpt: USEErrorsOpt{
					ErrorLabels: []string{"error"},
				},
			},
			want: &USE{
				Utilization: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: "foo",
					Name:      "qux",
					Help:      "quux",
				}, []string{"fred", "plugh", "xyzzy"}),
				Saturation: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: "foo",
					Name:      "bar",
					Help:      "baz",
				}, []string{"grault", "garply", "waldo"}),
				Errors: prometheus.NewCounterVec(prometheus.CounterOpts{
					Namespace: "foo",
					Name:      "errors_total",
					Help:      "Number of errors",
				}, []string{"error"}),
			},
			wantErr: false,
		},
		"missing Namespace": {
			opts: USEOpts{
				UtilizationOpt: USEUtilizationOpt{
					UtilizationName:   "qux",
					UtilizationHelp:   "quux",
					UtilizationLabels: []string{"fred", "plugh", "xyzzy"},
				},
				SaturationOpt: USESaturationOpt{
					SaturationName:   "bar",
					SaturationHelp:   "baz",
					SaturationLabels: []string{"grault", "garply", "waldo"},
				},
				ErrorsOpt: USEErrorsOpt{
					ErrorLabels: []string{"error"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		"missing UtilizationName": {
			opts: USEOpts{
				Namespace: "foo",
				UtilizationOpt: USEUtilizationOpt{
					UtilizationHelp:   "quux",
					UtilizationLabels: []string{"fred", "plugh", "xyzzy"},
				},
				SaturationOpt: USESaturationOpt{
					SaturationName:   "bar",
					SaturationHelp:   "baz",
					SaturationLabels: []string{"grault", "garply", "waldo"},
				},
				ErrorsOpt: USEErrorsOpt{
					ErrorLabels: []string{"error"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		"missing UtilizationHelp": {
			opts: USEOpts{
				Namespace: "foo",
				UtilizationOpt: USEUtilizationOpt{
					UtilizationName:   "qux",
					UtilizationLabels: []string{"fred", "plugh", "xyzzy"},
				},
				SaturationOpt: USESaturationOpt{
					SaturationName:   "bar",
					SaturationHelp:   "baz",
					SaturationLabels: []string{"grault", "garply", "waldo"},
				},
				ErrorsOpt: USEErrorsOpt{
					ErrorLabels: []string{"error"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		"missing UtilizationLabels": {
			opts: USEOpts{
				Namespace: "foo",
				UtilizationOpt: USEUtilizationOpt{
					UtilizationName: "qux",
					UtilizationHelp: "quux",
				},
				SaturationOpt: USESaturationOpt{
					SaturationName:   "bar",
					SaturationHelp:   "baz",
					SaturationLabels: []string{"grault", "garply", "waldo"},
				},
				ErrorsOpt: USEErrorsOpt{
					ErrorLabels: []string{"error"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		"missing SaturationName": {
			opts: USEOpts{
				Namespace: "foo",
				UtilizationOpt: USEUtilizationOpt{
					UtilizationName:   "qux",
					UtilizationHelp:   "quux",
					UtilizationLabels: []string{"fred", "plugh", "xyzzy"},
				},
				SaturationOpt: USESaturationOpt{
					SaturationHelp:   "baz",
					SaturationLabels: []string{"grault", "garply", "waldo"},
				},
				ErrorsOpt: USEErrorsOpt{
					ErrorLabels: []string{"error"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		"missing SaturationHelp": {
			opts: USEOpts{
				Namespace: "foo",
				UtilizationOpt: USEUtilizationOpt{
					UtilizationName:   "qux",
					UtilizationHelp:   "quux",
					UtilizationLabels: []string{"fred", "plugh", "xyzzy"},
				},
				SaturationOpt: USESaturationOpt{
					SaturationName:   "bar",
					SaturationLabels: []string{"grault", "garply", "waldo"},
				},
				ErrorsOpt: USEErrorsOpt{
					ErrorLabels: []string{"error"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		"missing SaturationLabels": {
			opts: USEOpts{
				Namespace: "foo",
				UtilizationOpt: USEUtilizationOpt{
					UtilizationName:   "qux",
					UtilizationHelp:   "quux",
					UtilizationLabels: []string{"fred", "plugh", "xyzzy"},
				},
				SaturationOpt: USESaturationOpt{
					SaturationName: "bar",
					SaturationHelp: "baz",
				},
				ErrorsOpt: USEErrorsOpt{
					ErrorLabels: []string{"error"},
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
			got, err := NewUSE(opts)
			if (err != nil) != wantErr {
				t.Errorf("NewUSE() error = %v, wantErr %v", err, wantErr)

				return
			}
			if want == nil && got == nil {
				return
			}
			assert.EqualExportedValues(t, *want, *got)
		})
	}
}

func TestUSEMetricNames(t *testing.T) {
	t.Parallel()

	// Test for custom names
	use, err := NewUSE(USEOpts{
		Namespace: "test",
		UtilizationOpt: USEUtilizationOpt{
			UtilizationName:   "custom_utilization",
			UtilizationHelp:   "Test utilization help",
			UtilizationLabels: []string{"label"},
		},
		SaturationOpt: USESaturationOpt{
			SaturationName:   "custom_saturation",
			SaturationHelp:   "Test saturation help",
			SaturationLabels: []string{"label"},
		},
		ErrorsOpt: USEErrorsOpt{
			ErrorName:   "custom_errors",
			ErrorLabels: []string{"error"},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, "custom_utilization", use.UtilizationMetricName())
	assert.Equal(t, "custom_saturation", use.SaturationMetricName())
	assert.Equal(t, "custom_errors", use.ErrorMetricName())

}

func TestUSECollectors(t *testing.T) {
	t.Parallel()

	use, err := NewUSE(USEOpts{
		Namespace: "test",
		UtilizationOpt: USEUtilizationOpt{
			UtilizationName:   "cpu_utilization_ratio",
			UtilizationHelp:   "CPU utilization ratio",
			UtilizationLabels: []string{"cpu"},
		},
		SaturationOpt: USESaturationOpt{
			SaturationName:   "memory_saturation_bytes",
			SaturationHelp:   "Memory saturation in bytes",
			SaturationLabels: []string{"type"},
		},
		ErrorsOpt: USEErrorsOpt{
			ErrorLabels: []string{"type"},
		},
	})
	assert.NoError(t, err)

	collectors := use.Collectors()

	assert.Equal(t, 3, len(collectors))

	gaugeCount := 0
	counterCount := 0

	for _, collector := range collectors {
		switch collector.(type) {
		case *prometheus.GaugeVec:
			gaugeCount++
		case *prometheus.CounterVec:
			counterCount++
		default:
			t.Errorf("Unexpected collector type: %T", collector)
		}
	}

	assert.Equal(t, 2, gaugeCount, "Should have 2 GaugeVec")
	assert.Equal(t, 1, counterCount, "Should have 1 CounterVec")
}
