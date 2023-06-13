# Promstrap

## A simple, lightweight metrics package that teams can use to bootstrap the instrumentaion of Go applications with Prometheus.

## Motivation
Observability is challenging. The instrumentation of applications at scale to enable observability is even more challenging if teams do not understand the fundamentals of what they should be measuring and when to do so. This package was born from having to teach these concepts to various teams and having to  write the same bootstrapping code on many projects. This package aims to provide the guard rails for teams to expedite learning and more easily adopt proven instrumentation strategies while empowering teams to create and being good stewards of **[SL*s](https://sre.google/sre-book/service-level-objectives/)**.

Besides providing a straightforward way of creating Prometheus counters, gauges, histograms and summaries, it also provides away to easily bootstrap well known intrumentation strategies that teams can reason about and leverage in their applications/system. This package was inspired by and uses language in it’s code comments verbatim from these resources:
- [Instrumenting Applications](https://training.promlabs.com/training/instrumenting-applications) - [PromLabs](https://promlabs.com/)
- [Prometheus Metric Types](https://prometheus.io/docs/concepts/metric_types/)
- [The USE Method](https://www.brendangregg.com/usemethod.html)
- [Monitoring Microservices](https://www.slideshare.net/weaveworks/monitoring-microservices)
- [The Four Golden Signals - Monitoring Distributed Systems - The SRE Book](https://sre.google/sre-book/monitoring-distributed-systems/#xref_monitoring_golden-signals)

## Why Prometheus
**[Prometheus](https://prometheus.io)** is an open-source monitoring and alerting toolkit that can be used to collect metrics from a variety of sources, including applications, infrastructure, and cloud services. It is a popular choice for observability because it is highly scalable, flexible, and easy to use.

Prometheus's support for multi-dimensional data collection and querying makes it a particularly good fit for monitoring **[cloud-native](https://www.cncf.io/)** applications, which are often composed of a large number of services. When services are instrumented with Prometheus, they can expose or push metrics that Prometheus can collect and store. By collecting metrics from each service, Prometheus can provide a holistic view of a system's health and performance.

In addition to its monitoring capabilities, Prometheus can also be used to generate alerts. These alerts can be used to notify users when certain metrics exceed predefined thresholds.

Here are some of the benefits of instrumenting applications with Prometheus:
- Increased visibility into the health and performance of applications
- Improved ability to identify and troubleshoot problems
- Enhanced ability to prevent problems from occurring
- Increased agility and responsiveness to changes

## Usage

### Counter
```go
package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rabellamy/promstrap"
)

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	candyCounter, _ := promstrap.NewCounterWithLabels(promstrap.MetricsOpts{
		Namespace: "candy",
		Name:      "mandm_total",
		Help:      "Counts the number of M&Ms",
		Labels:    []string{"color"},
	})

	err := promstrap.RegisterMetrics(candyCounter)
	if err != nil {
		fmt.Println(err.Error())
	}

	r := chi.NewRouter()
	r.Get("/m&ms", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("M&Ms!!\n"))

		candyCounter.WithLabelValues("red").Add(10)
		candyCounter.WithLabelValues("orange").Add(5)
		candyCounter.WithLabelValues("green").Add(5)
		candyCounter.WithLabelValues("blue").Add(5)
		candyCounter.WithLabelValues("yellow").Add(10)
		candyCounter.WithLabelValues("brown").Add(15)
	})

	http.ListenAndServe(":8080", r)
}

```

#### Result
```
# HELP candy_mandm_total Counts the number of M&Ms
# TYPE candy_mandm_total counter
candy_mandm_total{color="blue"} 5
candy_mandm_total{color="brown"} 15
candy_mandm_total{color="green"} 5
candy_mandm_total{color="orange"} 5
candy_mandm_total{color="red"} 10
candy_mandm_total{color="yellow"} 10
```

### Gauge
```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rabellamy/promstrap"
)

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	temp, _ := promstrap.NewGaugeWithLabels(promstrap.MetricsOpts{
		Namespace: "weather",
		Name:      "current_temperature_celsius",
		Help:      "The current temperature Celsius",
		Labels:    []string{"location"},
	})

	err := promstrap.RegisterMetrics(temp)
	if err != nil {
		fmt.Println(err.Error())
	}

	r := chi.NewRouter()
	r.Get("/temp", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("It's hot!!\n"))
		temp.WithLabelValues(time.Now().Format("New York, NY")).Set(42)
	})

	http.ListenAndServe(":8080", r)
}
```

#### Result
```
# HELP weather_current_temperature_celsius The current temperature Celsius
# TYPE weather_current_temperature_celsius gauge
weather_current_temperature_celsius{date="06-30-2023 21:59:38"} 42
```
