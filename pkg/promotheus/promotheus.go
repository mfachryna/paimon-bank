package promotheus

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	histogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "paimon_bank_requests",
		Help:    "Duration of HTTP requests in seconds.",
		Buckets: prometheus.LinearBuckets(1, 1, 10),
	}, []string{"path", "method", "status"})
	reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "chi_pattern_requests_total",
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path (with patterns).",
			ConstLabels: prometheus.Labels{"service": "paimon_bank_requests"},
		},
		[]string{"code", "method", "path"},
	)
	latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "chi_pattern_request_duration_milliseconds",
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path (with patterns).",
		ConstLabels: prometheus.Labels{"service": "paimon_bank_requests"},
		Buckets:     []float64{300, 1200, 5000},
	},
		[]string{"code", "method", "path"},
	)
)

func PrometheusMiddleware(next http.Handler) http.Handler {

	prometheus.MustRegister(reqs)
	prometheus.MustRegister(latency)
	registry := prometheus.NewRegistry()

	registry.MustRegister(histogram)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		rw := &responseWriter{w, http.StatusOK}

		defer func(startTime time.Time, rw *responseWriter) {
			duration := time.Since(startTime).Seconds()
			histogram.WithLabelValues(r.URL.Path, r.Method, strconv.Itoa(rw.status)).Observe(duration)
			reqs.WithLabelValues(strconv.Itoa(rw.status), r.Method, r.URL.Path).Inc()
			latency.WithLabelValues(strconv.Itoa(rw.status), r.Method, r.URL.Path).Observe(float64(time.Since(startTime).Nanoseconds()) / 1000000)
		}(startTime, rw)

		next.ServeHTTP(rw, r)

	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
