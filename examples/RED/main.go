package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rabellamy/promstrap"
	"github.com/rabellamy/promstrap/strategy"
)

type workTimeBox struct {
	min float64
	max float64
}

// Allows us to simulate some sort of processing
func doTheWork(wtb workTimeBox) {
	workSim := rand.Float64() * (wtb.max - wtb.min)
	time.Sleep(time.Duration(workSim) * time.Second)
}

func main() {
	var err error

	reporter, _ := strategy.NewRED(strategy.REDOpts{
		RequestType:    "http",
		Namespace:      "bar",
		RequestLabels:  []string{"path", "verb"},
		DurationLabels: []string{"path"},
	})

	// dereference struct to avoid panic
	err = promstrap.RegisterComplexMetrics(*reporter)
	if err != nil {
		fmt.Println(err.Error())
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	r := chi.NewRouter()
	r.Get("/happy", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		// Records that a request took place
		reporter.Requests.WithLabelValues("/happy", "GET").Inc()

		doTheWork(workTimeBox{
			min: .01,
			max: 3.0,
		})

		_, err = w.Write([]byte("You are now happy!!\n"))
		if err != nil {
			// Records the error
			reporter.Errors.WithLabelValues(err.Error()).Inc()
		}

		// Records the duration of the request
		reporter.Duration.WithLabelValues("/happy").Observe(time.Since(t).Seconds())
	})

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err.Error())
	}
}
