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

## Basic Usage

### Counter
```go
// Creates a Prometheus counter with labels to count the number of errors.
// A counter is a cumulative metric that represents a single monotonically
// increasing counter whose value can only increase or be reset to zero on restart.
errors, err := metrics.NewCounterWithLabels(metrics.CounterOpts{
	Namespace: "service_name",
	Name:      "errors_total",
	Help:      "Number of errors",
	Labels:    []string{"error"},
})
if err != nil {
	return nil, err
}


// Registers errors with the Prometheus DefaultRegisterer
prometheus.MustRegister(errors)

// Increment by 1
errors.WithLabelValues("my error").Inc()
```

### Gauge
```go
// Creates a Prometheus Gauge with labels to measure the length of a queue.
// A gauge is a metric that represents a single numerical value that can
// arbitrarily go up and down.
length, err := metrics.NewGaugeWithLabels(metrics.GaugeOpts{
	Namespace: "service_name",
	Name:      "queue_length",
	Help:      "The number of items in the queue.",
	Labels:    []string{"queue"},
})
if err != nil {
	return nil, err
}

// Registers length with the Prometheus DefaultRegisterer
prometheus.MustRegister(length)

// Use Set() when you know the absolute value
length.WithLabelValues("foo").Set(0)

// Increment by 1
length.WithLabelValues("foo").Inc()
// Decrement by 1
length.WithLabelValues("foo").Dec()
// Increment by 14
length.WithLabelValues("foo").Add(14)
// Decrement by 7
length.WithLabelValues("foo").Sub(7)
```

### Histogram
```go
// Creates a Prometheus Histogram with labels to measure the duration of request
// in seconds.
// A histogram samples observations (usually things like request durations or
// response sizes) and counts them in configurable buckets.
duration, err := metrics.NewHistogramWithLabels(metrics.HistogramOpts{
	Namespace: "service_name",
	Name:      "http_request_duration_seconds",
	Help:      "Duration of request in seconds",
	Labels:    []string{"path"},
	Buckets:   []float64{0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
})
if err != nil {
	return nil, err
}

// Registers duration with the Prometheus DefaultRegisterer
prometheus.MustRegister(duration)

// Record that the HTTP request to the /happy path took 0.5 seconds to serve
duration.WithLabelValues("/happy").Observe(0.5)
```

### Summary
```go
// Creates a Prometheus Summary with labels to measure the duration of request
// in seconds.
// A summary samples observations (usually things like request durations and
// response sizes).
duration, err := metrics.NewSummaryWithLabels(metrics.SummaryOpts{
	Namespace: "service_name",
	Name:      "http_request_duration_seconds",
	Help:      "Duration of request in seconds",
	Labels:    []string{"path"},
	// Objectives defines the quantile rank estimates with their respective
	// absolute error.
	Objectives: map[float64]float64{
		0.5:  0.05,  // 50th percentile with a max. absolute error of 0.05.
		0.9:  0.01,  // 90th percentile with a max. absolute error of 0.01.
		0.99: 0.001, // 99th percentile with a max. absolute error of 0.001.
	},
})
if err != nil {
	return nil, err
}

// Registers duration with the Prometheus DefaultRegisterer
prometheus.MustRegister(duration)

// Record that the HTTP request to the /happy path took 0.5 seconds to serve
duration.WithLabelValues("/happy").Observe(0.5)
```
