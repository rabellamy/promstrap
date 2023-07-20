package promstrap

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestNewSummaryWithLabels(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		opts    SummaryOpts
		want    *prometheus.SummaryVec
		wantErr bool
	}{
		"all good": {
			opts: SummaryOpts{
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
			opts: SummaryOpts{
				Namespace: "",
				Name:      "the_name",
				Help:      "Some help text",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Name": {
			opts: SummaryOpts{
				Namespace: "the_namespace",
				Name:      "",
				Help:      "Some help text",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Help": {
			opts: SummaryOpts{
				Namespace: "the_namespace",
				Name:      "the_name",
				Help:      "",
				Labels:    []string{"yo", "bro", "flow"},
			},
			want:    nil,
			wantErr: true,
		},
		"no Labels": {
			opts: SummaryOpts{
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
