package promstrap

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestNewCounterWithLabels(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		opts    CounterOpts
		want    *prometheus.CounterVec
		wantErr bool
	}{
		"all good": {
			opts: CounterOpts{
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
			opts: CounterOpts{
				Namespace: "",
				Name:      "the_name",
				Help:      "Some help text",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Name": {
			opts: CounterOpts{
				Namespace: "the_namespace",
				Name:      "",
				Help:      "Some help text",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Help": {
			opts: CounterOpts{
				Namespace: "the_namespace",
				Name:      "the_name",
				Help:      "",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Labels": {
			opts: CounterOpts{
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
