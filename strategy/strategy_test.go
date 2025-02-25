package strategy

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"gitlab.alticeustech.com/platform-engineering/observability-infrastructure/promstrap/metrics"
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

func (m testValidStrategy) Collectors() map[string]prometheus.Collector {
	return map[string]prometheus.Collector{
		"foo": m.Foo,
	}
}

type testInvalidStrategy struct {
	CounterOne *prometheus.CounterVec
	CounterTwo *prometheus.CounterVec
}

func (m testInvalidStrategy) Register() error {
	err := RegisterStrategyFields(m)
	if err != nil {
		return err
	}

	return nil
}

func (m testInvalidStrategy) Collectors() map[string]prometheus.Collector {
	return map[string]prometheus.Collector{}
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
		t.Error(err)
	}

	duplicateCounter, err := metrics.NewCounterWithLabels(metrics.CounterOpts{
		Namespace: "duplicate_counter",
		Name:      "duplicate",
		Help:      "duplicate test",
		Labels:    []string{"baz"},
	})

	if err != nil {
		t.Error(err)
	}

	validComplexMetric := testValidStrategy{
		Foo: testCounter,
	}

	invalidStrategy := testInvalidStrategy{
		CounterOne: duplicateCounter,
		CounterTwo: duplicateCounter,
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
		"invalid strategy": {
			strategy: invalidStrategy,
			wantErr:  true,
		},
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
