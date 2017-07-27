package thetis

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

type httpMetric struct {
	collector *prometheus.CounterVec
}

func (h *httpMetric) GetCollector() prometheus.Collector {
	responses := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_status",
		Help: "HTTP Status per request",
	}, []string{"code", "route", "method"},
	)
	/*histogram: = prometheus.NewHistogram(prometheus.HistogramOpts {
	    Name: "http_requests",
	    Help: "HTTP Requests status and size per route.",
	})*/
	h.collector = responses
	return responses
}

func (h *httpMetric) CreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		status := "0"
		if sw, ok := w.(*statusWriter); ok {
			status = fmt.Sprint(sw.status)
		}
		uri := r.RequestURI
		h.collector.WithLabelValues(status, uri, method).Inc()
	}
}
