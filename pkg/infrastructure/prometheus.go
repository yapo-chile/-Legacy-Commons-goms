package infrastructure

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.schibsted.io/Yapo/goms/pkg/interfaces/loggers"
)

// Prometheus provides both, a way to instrument http.HandlerFunc with
// Prometheus, and a Prometheus http.server that can exposes metrics in a given port
type Prometheus struct {
	// common  metrics for handlers
	// counter metric of HTTP request qty
	counter *prometheus.CounterVec
	// duration metric of request duration using http
	duration prometheus.ObserverVec
	// inFlight metric of requests currently being served
	inFlight prometheus.Gauge
	// requestSize   metric of HTTP request size
	requestSize prometheus.ObserverVec
	// responseSize  metric of HTTP response size
	responseSize prometheus.ObserverVec
	// Exporter params
	// server exposes all metrics on /metrics path using a given port
	server *http.Server
	// logger logs runtime messages
	logger loggers.Logger
	// enabled enables prometheus exporter
	enabled bool
}

// MakePrometheusExporter Builds a fresh Prometheus, initializing its
// metrics
func MakePrometheusExporter(port string, enabled bool, logger loggers.Logger) *Prometheus {
	p := Prometheus{
		counter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_requests_total",
				Help: "A counter for requests to the wrapped handler.",
			},
			[]string{"handler", "method"},
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
			[]string{"handler", "method"},
		),
		responseSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "response_size_bytes",
				Help:    "A histogram of response sizes for requests.",
				Buckets: []float64{200, 500, 900, 1500},
			},
			[]string{"handler", "method"},
		),
		enabled: enabled,
		logger:  logger,
	}

	// Register all of the common metrics in the standard registry
	prometheus.MustRegister(p.counter, p.duration, p.inFlight, p.requestSize, p.responseSize)

	// start prometheus exposer server in /metrics endpoint
	p.expose(port)
	return &p
}

// TrackHandlerFunc instruments handler with Prometheus, adding every
// configured metric
func (p *Prometheus) TrackHandlerFunc(handlerName string, handler http.HandlerFunc) http.HandlerFunc {
	if !p.enabled {
		return handler
	}
	// In Flight requests
	handler = promhttp.InstrumentHandlerInFlight(p.inFlight, handler).(http.HandlerFunc)

	// Request Counter
	handler = promhttp.InstrumentHandlerCounter(
		p.counter.MustCurryWith(prometheus.Labels{"handler": handlerName}), handler)

	// Duration
	handler = promhttp.InstrumentHandlerDuration(
		p.duration.MustCurryWith(prometheus.Labels{"handler": handlerName}), handler)

	// Request Size
	handler = promhttp.InstrumentHandlerRequestSize(
		p.requestSize.MustCurryWith(prometheus.Labels{"handler": handlerName}), handler)

	// Response Size
	handler = promhttp.InstrumentHandlerResponseSize(
		p.responseSize.MustCurryWith(prometheus.Labels{"handler": handlerName}),
		handler,
	).(http.HandlerFunc)

	// return tracked handler
	return handler.ServeHTTP
}

// NewCounterVector creates a new instance of CounterVector
func (*Prometheus) NewCounterVector(name, help string, labels []string) loggers.CounterVector {
	vector := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		},
		labels,
	)
	prometheus.MustRegister(vector)
	return &CounterVector{
		vector: vector,
	}
}

// CounterVector is a Collector that bundles a set of Counters that all share the
// same Desc, but have different values for their variable labels. This is used
// if you want to count the same thing partitioned by various dimensions
// (e.g. number of HTTP requests, partitioned by response code and
// method). Create instances with NewCounterVec.
type CounterVector struct {
	vector *prometheus.CounterVec
}

// WithLabelValues selects counter with label
// ex: myVec.WithLabelValues("404", "GET").Add(42)
func (v *CounterVector) WithLabelValues(labels ...string) loggers.Counter {
	return v.vector.WithLabelValues(labels...)
}

// expose starts prometheus exporter metrics server exposing metrics in "/metrics" path
func (p *Prometheus) expose(port string) {
	if !p.enabled {
		return
	}
	p.server = &http.Server{Addr: ":" + port}
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := p.server.ListenAndServe(); err != http.ErrServerClosed {
			p.logger.Error("Prometheus error: %s", err)
		}
	}()
}

// Close closes prometheus server
func (p *Prometheus) Close() error {
	if !p.enabled {
		return nil
	}
	return p.server.Close()
}
