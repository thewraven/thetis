package thetis

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

type httpMetric struct {
	collector *prometheus.CounterVec
}

//GetCollector reports the status,route and method
//of every request
func (h *httpMetric) GetCollector() prometheus.Collector {
	responses := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_status",
		Help: "HTTP Status per request",
	}, []string{"code", "route", "method"},
	)
	h.collector = responses
	return responses
}

//CreateHandler returns a handler that will be called
//every request, and it is responsible of writing
//the desired metrics to the handler
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
