package prometheushelper

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/*
Rate:
sum(rate(request_duration_seconds_count{job=“...”}[1m]))

Errors:
sum(rate(request_duration_seconds_count{job=“...”,status_code!~”2..”}[1m]))

Duration:
histogram_quantile(0.99, sum(rate(request_duration_seconds_bucket{job=“...}[1m])) by (le))
*/
type HttpMetrics struct {
	RequestTotal             *prometheus.CounterVec
	RequestDurationHistogram *prometheus.HistogramVec
}

func NewHttpMetrics() HttpMetrics {
	return HttpMetrics{
		RequestTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "sample",
			Subsystem: "http",
			Name:      "request_total",
			Help:      "total HTTP requests processed",
		}, []string{"code", "method"}),
		RequestDurationHistogram: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "sample",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "Seconds spent serving HTTP requests",
			Buckets:   prometheus.DefBuckets,
		}, []string{"code", "method"}),
	}
}

func InstrumentHandler(next http.HandlerFunc, _http HttpMetrics) http.HandlerFunc {
	return promhttp.InstrumentHandlerCounter(_http.RequestTotal, promhttp.InstrumentHandlerDuration(_http.RequestDurationHistogram, next))
}
