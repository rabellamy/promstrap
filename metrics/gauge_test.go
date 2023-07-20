package promstrap

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestNewGaugeWithLabels(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		opts    GaugeOpts
		want    *prometheus.GaugeVec
		wantErr bool
	}{
		"all good": {
			opts: GaugeOpts{
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
			opts: GaugeOpts{
				Namespace: "",
				Name:      "the_name",
				Help:      "Some help text",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Name": {
			opts: GaugeOpts{
				Namespace: "the_namespace",
				Name:      "",
				Help:      "Some help text",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Help": {
			opts: GaugeOpts{
				Namespace: "the_namespace",
				Name:      "the_name",
				Help:      "",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Labels": {
			opts: GaugeOpts{
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
