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
				Namespace:         "foo",
				UtilizationName:   "qux",
				UtilizationHelp:   "quux",
				UtilizationLabels: []string{"fred, plugh, xyzzy"},
				SaturationName:    "bar",
				SaturationHelp:    "baz",
				SaturationLabels:  []string{"grault, garply, waldo"},
			},
			want: &USE{
				Utilization: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: "foo",
					Name:      "qux",
					Help:      "quux",
				}, []string{"fred, plugh, xyzzy"}),
				Saturation: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: "foo",
					Name:      "bar",
					Help:      "baz",
				}, []string{"grault, garply, waldo"}),
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
				Namespace:         "",
				UtilizationName:   "qux",
				UtilizationHelp:   "quux",
				UtilizationLabels: []string{"fred, plugh, xyzzy"},
				SaturationName:    "bar",
				SaturationHelp:    "baz",
				SaturationLabels:  []string{"grault, garply, waldo"},
			},
			want:    nil,
			wantErr: true,
		},
		"missing UtilizationName": {
			opts: USEOpts{
				Namespace:         "foo",
				UtilizationName:   "",
				UtilizationHelp:   "quux",
				UtilizationLabels: []string{"fred, plugh, xyzzy"},
				SaturationName:    "bar",
				SaturationHelp:    "baz",
				SaturationLabels:  []string{"grault, garply, waldo"},
			},
			want:    nil,
			wantErr: true,
		},
		"missing UtilizationHelp": {
			opts: USEOpts{
				Namespace:         "foo",
				UtilizationName:   "qux",
				UtilizationHelp:   "",
				UtilizationLabels: []string{"fred, plugh, xyzzy"},
				SaturationName:    "bar",
				SaturationHelp:    "baz",
				SaturationLabels:  []string{"grault, garply, waldo"},
			},
			want:    nil,
			wantErr: true,
		},
		"missing UtilizationLabels": {
			opts: USEOpts{
				Namespace:         "foo",
				UtilizationName:   "qux",
				UtilizationHelp:   "quux",
				UtilizationLabels: nil,
				SaturationName:    "bar",
				SaturationHelp:    "baz",
				SaturationLabels:  []string{"grault, garply, waldo"},
			},
			want:    nil,
			wantErr: true,
		},
		"missing SaturationName": {
			opts: USEOpts{
				Namespace:         "foo",
				UtilizationName:   "qux",
				UtilizationHelp:   "quux",
				UtilizationLabels: []string{"fred, plugh, xyzzy"},
				SaturationName:    "",
				SaturationHelp:    "baz",
				SaturationLabels:  []string{"grault, garply, waldo"},
			},
			want:    nil,
			wantErr: true,
		},
		"missing SaturationHelp": {
			opts: USEOpts{
				Namespace:         "foo",
				UtilizationName:   "qux",
				UtilizationHelp:   "quux",
				UtilizationLabels: []string{"fred, plugh, xyzzy"},
				SaturationName:    "bar",
				SaturationHelp:    "",
				SaturationLabels:  []string{"grault, garply, waldo"},
			},
			want:    nil,
			wantErr: true,
		},
		"missing SaturationLabels": {
			opts: USEOpts{
				Namespace:         "foo",
				UtilizationName:   "qux",
				UtilizationHelp:   "quux",
				UtilizationLabels: []string{"fred, plugh, xyzzy"},
				SaturationName:    "bar",
				SaturationHelp:    "baz",
				SaturationLabels:  nil,
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

	// Test for default names
	useDefault, err := NewUSE(USEOpts{
		Namespace: "test",
		UtilizationOpt: USEUtilizationOpt{
			UtilizationName:   "utilization",
			UtilizationHelp:   "Test utilization help",
			UtilizationLabels: []string{"label"},
		},
		SaturationOpt: USESaturationOpt{
			SaturationName:   "saturation",
			SaturationHelp:   "Test saturation help",
			SaturationLabels: []string{"label"},
		},
		ErrorsOpt: USEErrorsOpt{
			ErrorLabels: []string{"error"},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, "utilization", useDefault.UtilizationMetricName())
	assert.Equal(t, "saturation", useDefault.SaturationMetricName())
	assert.Equal(t, "errors_total", useDefault.ErrorMetricName())
}
