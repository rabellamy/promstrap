package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rabellamy/promstrap/strategy"
)

func main() {
	useExample, err := strategy.NewUSE(strategy.USEOpts{
		Namespace: "system",
		UtilizationOpt: strategy.USEUtilizationOpt{
			UtilizationName:   "memory_utilization_ratio",
			UtilizationHelp:   "Memory utilization as a ratio of used to total",
			UtilizationLabels: []string{"type"},
		},
		SaturationOpt: strategy.USESaturationOpt{
			SaturationName:   "memory_saturation_bytes",
			SaturationHelp:   "Amount of memory queued/waiting to be freed",
			SaturationLabels: []string{"type"},
		},
		ErrorsOpt: strategy.USEErrorsOpt{
			ErrorLabels: []string{"type"},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	// Register metrics
	err = useExample.Register()
	if err != nil {
		fmt.Println(err.Error())
	}

	// Expose metrics on /metrics
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
