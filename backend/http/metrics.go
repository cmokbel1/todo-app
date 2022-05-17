package http

import (
	"net/http"
	_ "net/http/pprof"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ListenAndServeDebug() error {
	h := http.NewServeMux()
	h.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":6060", h)
}

var metrics = struct {
	routeRequestCount *prometheus.CounterVec
	routeRequestTime  *prometheus.CounterVec
}{
	routeRequestCount: promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "todo_http_route_request_count",
		Help: "Total number of requests per route",
	}, []string{"method", "path"}),
	routeRequestTime: promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "todo_http_route_request_time",
		Help: "Total number of request time (seconds) per route",
	}, []string{"method", "path"}),
}

func monitorMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		if strings.HasPrefix(r.URL.Path, "/api") {
			// route pattern is built as request passes through the chain of handlers
			// so we have to wait until here to determine the full route.
			// see: https://github.com/go-chi/chi/issues/150#issuecomment-278850733
			route := chi.RouteContext(r.Context()).RoutePattern()
			metrics.routeRequestCount.WithLabelValues(r.Method, route).Inc()
			metrics.routeRequestTime.WithLabelValues(r.Method, route).Add(float64(time.Since(start).Seconds()))
		}
	})
}
