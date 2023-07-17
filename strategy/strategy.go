package strategy

import (
	"errors"
	"reflect"

	"github.com/prometheus/client_golang/prometheus"
)

// Strategy describes a collection of metrics built with a struct, allowing any multi-level
// combination of the 4 Prometheus metric types and other Strategies.
// e.g. RED - describes a set of metrics that work well for monitoring
// request-handling services.
//
//	type RED struct {
//		Requests *prometheus.CounterVec
//		Errors *prometheus.CounterVec
//		Duration *prometheus.HistogramVec
//	}
type Strategy interface {
	Register() error
}

// RegisterStrategyFields registers the fields of metrics and Strategies that
// comprise a Strategy with the Prometheus DefaultRegisterer.
func RegisterStrategyFields(s Strategy) error {
	l := reflect.TypeOf(s).NumField()
	if l < 1 {
		return errors.New("strategies needs at least one field")
	}

	values := reflect.ValueOf(s)

	for i := 0; i < l; i++ {
		switch v := values.Field(i).Interface().(type) {
		case prometheus.Collector:
			prometheus.MustRegister(v)

		// Allows for the composability of strategies
		case Strategy:
			err := v.Register()
			if err != nil {
				return err
			}

		default:
			return errors.New("field is not a Prometheus Collector nor Strategy")
		}
	}

	return nil
}
