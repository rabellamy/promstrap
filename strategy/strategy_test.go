package strategy

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rabellamy/promstrap/metrics"
)

type testValidStrategy struct {
	Foo *prometheus.CounterVec
}

func (m testValidStrategy) Register() error {
	err := RegisterStrategyFields(m)
	if err != nil {
		return err
	}

	return nil
}

func TestRegisterStrategyFields(t *testing.T) {
	t.Parallel()

	var err error

	testCounter, err := metrics.NewCounterWithLabels(metrics.CounterOpts{
		Namespace: "test_counter",
		Name:      "foo",
		Help:      "bar",
		Labels:    []string{"baz"},
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	validComplexMetric := testValidStrategy{
		Foo: testCounter,
	}

	tests := map[string]struct {
		strategy Strategy
		wantErr  bool
	}{
		"valid strategy": {
			strategy: validComplexMetric,
			wantErr:  false,
		},
		// TODO: more test cases
	}

	for name, tt := range tests {
		strategy := tt.strategy
		wantErr := tt.wantErr

		t.Run(name, func(t *testing.T) {
			if err := RegisterStrategyFields(strategy); (err != nil) != wantErr {
				t.Errorf("RegisterStrategyFields() error = %v, wantErr %v", err, wantErr)
			}
		})
	}
}
