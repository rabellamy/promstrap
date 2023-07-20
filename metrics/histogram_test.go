package promstrap

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestNewHistogramWithLabels(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		opts    HistogramOpts
		want    *prometheus.HistogramVec
		wantErr bool
	}{
		"all good": {
			opts: HistogramOpts{
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
			opts: HistogramOpts{
				Namespace: "",
				Name:      "the_name",
				Help:      "Some help text",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Name": {
			opts: HistogramOpts{
				Namespace: "the_namespace",
				Name:      "",
				Help:      "Some help text",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Help": {
			opts: HistogramOpts{
				Namespace: "the_namespace",
				Name:      "the_name",
				Help:      "",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Labels": {
			opts: HistogramOpts{
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
