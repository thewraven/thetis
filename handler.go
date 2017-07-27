package thetis

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

type Handler struct {
	metrics        []Metric
	wrappedHandler http.Handler
}

type Metric interface {
	GetCollector() prometheus.Collector
	CreateHandler() http.HandlerFunc
}

func NewWithDefaults(router http.Handler) *Handler {
	return &Handler{
		metrics:        []Metric{&httpMetric{}},
		wrappedHandler: router,
	}
}

func NewHandler(router http.Handler) *Handler {
	metrics := make([]Metric, 0)
	return &Handler{
		metrics:        metrics,
		wrappedHandler: router,
	}
}

func (h *Handler) Add(c Metric) {
	h.metrics = append(h.metrics, c)
}

func (h *Handler) RegisterAll() {
	collectors := make([]prometheus.Collector, len(h.metrics))
	for i, m := range h.metrics {
		collectors[i] = m.GetCollector()
	}
	prometheus.MustRegister(collectors...)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sw := &statusWriter{w: w}
	h.wrappedHandler.ServeHTTP(sw, r)
	for _, m := range h.metrics {
		m.CreateHandler().ServeHTTP(sw, r)
	}
}

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
