package thetis

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//Handler is a custom implementation of an http.Handler
//that calls the custom metrics on every route.
type Handler struct {
	metrics        []Metric
	wrappedHandler http.Handler
}

//Metric is a definition of a prometheus metric.
//It will return the prometheus.Collector to be used
//And Call the returned http.HandlerFunc
//for every request handled
type Metric interface {
	GetCollector() prometheus.Collector
	CreateHandler() http.HandlerFunc
}

//NewWithDefaults returns an http.Handler that process
//all the requests. Receives a router, which contains
//your logic.
func NewWithDefaults(router http.Handler) *Handler {
	return &Handler{
		metrics:        []Metric{&httpMetric{}},
		wrappedHandler: router,
	}
}

//NewHandler returns an http.Handler without metrics.
//You need to add your custom metrics using Add()
func NewHandler(router http.Handler) *Handler {
	metrics := make([]Metric, 0)
	return &Handler{
		metrics:        metrics,
		wrappedHandler: router,
	}
}

//Add appends your custom metrics to the chain of
//callable handlers
func (h *Handler) Add(c ...Metric) {
	h.metrics = append(h.metrics, c...)
}

//RegisterAll registers the given Collectors.
//calls to prometheus.MustRegister could panic, please check the documentation
func (h *Handler) RegisterAll() http.Handler {
	collectors := make([]prometheus.Collector, len(h.metrics))
	for i, m := range h.metrics {
		collectors[i] = m.GetCollector()
	}
	prometheus.MustRegister(collectors...)
	return promhttp.Handler()
}

//ServeHTTP implements the http.Handler interface
//it must be used after a single call of RegisterAll
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sw := &statusWriter{w: w}
	h.wrappedHandler.ServeHTTP(sw, r)
	for _, m := range h.metrics {
		m.CreateHandler().ServeHTTP(sw, r)
	}
}

//statusWriter is a custom implementation of
//an http.ResponseWriter that can
//preserve the http status written in the response
type statusWriter struct {
	status int
	w      http.ResponseWriter
}

func (s *statusWriter) WriteHeader(status int) {
	s.status = status
	s.w.WriteHeader(status)
}

func (s *statusWriter) Write(data []byte) (int, error) {
	return s.w.Write(data)
}

func (s *statusWriter) Header() http.Header {
	return s.w.Header()
}
