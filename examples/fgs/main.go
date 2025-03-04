package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.alticeustech.com/platform-engineering/observability-infrastructure/promstrap/strategy"
)

func main() {
	fgsExample, err := strategy.NewFourGoldenSignals(strategy.FourGoldenSignalsOpts{
		Namespace: "service",
		LatencyOpt: strategy.FGSLatencyOpt{
			LatencyName:   "http_request_latency_seconds",
			LatencyType:   "http",
			LatencyHelp:   "HTTP request latency in seconds",
			LatencyLabels: []string{"method", "path"},
			Buckets:       []float64{.005, .01, .025, .05, .1, .25, .5, 1},
		},
		TrafficOpt: strategy.FGSTrafficOpt{
			TrafficName:   "http_server_requests_total",
			TrafficType:   "http",
			TrafficHelp:   "Total number of HTTP requests",
			TrafficLabels: []string{"method", "path", "status"},
		},
		ErrorsOpt: strategy.FGSErrorsOpt{
			ErrorHelp:   "Number of errors",
			ErrorLabels: []string{"type"},
		},
		SaturationOpt: strategy.FGSSaturationOpt{
			SaturationName:   "memory_heap_saturation_bytes",
			SaturationHelp:   "Memory heap usage in bytes",
			SaturationLabels: []string{"gc_type"},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	// Register metrics
	err = fgsExample.Register()
	if err != nil {
		fmt.Println(err.Error())
	}

	// Expose metrics on /metrics
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
