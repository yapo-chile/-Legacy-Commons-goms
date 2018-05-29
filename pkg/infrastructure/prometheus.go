package infrastructure

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusHandler struct {
	counter      *prometheus.CounterVec
	duration     prometheus.ObserverVec
	inFlight     prometheus.Gauge
	requestSize  prometheus.ObserverVec
	responseSize prometheus.ObserverVec
}

func MakePrometheusHandler() PrometheusHandler {
	h := PrometheusHandler{
		counter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_requests_total",
				Help: "A counter for requests to the wrapped handler.",
			},
			[]string{"code", "method"},
		),
		duration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "request_duration_seconds",
				Help:    "A histogram of latencies for requests.",
				Buckets: []float64{.25, .5, 1, 2.5, 5, 10},
			},
			[]string{"handler", "method"},
		),
		inFlight: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "in_flight_requests",
			Help: "A gauge of requests currently being served by the wrapped handler.",
		}),
		requestSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "request_size_bytes",
				Help:    "A histogram of request sizes for requests.",
				Buckets: []float64{50, 100, 200, 500, 1000, 1500},
			},
			[]string{},
		),
		responseSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "response_size_bytes",
				Help:    "A histogram of response sizes for requests.",
				Buckets: []float64{200, 500, 900, 1500},
			},
			[]string{},
		),
	}

	// Register all of the metrics in the standard registry.
	prometheus.MustRegister(h.counter, h.duration, h.inFlight, h.requestSize, h.responseSize)
	return h
}

func (h *PrometheusHandler) TrackHandlerFunc(pattern string, handler http.HandlerFunc) http.HandlerFunc {
	return promhttp.InstrumentHandlerInFlight(h.inFlight,
		promhttp.InstrumentHandlerCounter(h.counter,
			promhttp.InstrumentHandlerDuration(h.duration.MustCurryWith(prometheus.Labels{"handler": pattern}),
				promhttp.InstrumentHandlerRequestSize(h.requestSize,
					promhttp.InstrumentHandlerResponseSize(h.responseSize, handler),
				),
			),
		),
	).ServeHTTP
}

func (h *PrometheusHandler) Handler() http.Handler {
	return promhttp.Handler()
}
