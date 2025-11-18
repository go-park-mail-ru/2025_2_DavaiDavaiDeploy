package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	HitsTotal *prometheus.CounterVec
	name      string
	Times     *prometheus.HistogramVec
	Errors    *prometheus.CounterVec
}
