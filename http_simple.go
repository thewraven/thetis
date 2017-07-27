package thetis

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

type httpMetric struct {
	collector *prometheus.CounterVec
}

//GetCollector returns the default collector to be registered
//for the Handler
func (h *httpMetric) GetCollector() prometheus.Collector {
	responses := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_status",
		Help: "HTTP Status per request",
	}, []string{"code", "route", "method"},
	)
	h.collector = responses
	return responses
}

func (h *httpMetric) CreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		//ignoring the metrics route by default
		if uri == "/metrics" {
			return
		}
		method := r.Method
		status := "0"
		if sw, ok := w.(*statusWriter); ok {
			status = fmt.Sprint(sw.status)
		}
		h.collector.WithLabelValues(status, uri, method).Inc()
	}
}
