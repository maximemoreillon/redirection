package instrumentation

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

//statusRecorder to record the status code from the ResponseWriter
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func MeasureResponseDuration(next http.Handler) http.Handler {
	buckets := []float64{.003, .03, .1, .3, 1.5, 10}

	responseTimeHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		// Namespace: "namespace",
		Name:      "http_request_duration_seconds",
		Help:      "Histogram of response time for handler in seconds",
		Buckets:   buckets,
	}, []string{"route", "method", "status_code"})

	prometheus.MustRegister(responseTimeHistogram)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		rec := statusRecorder{w, 200}
		
		next.ServeHTTP(&rec, r)

		duration := time.Since(start)
		statusCode := strconv.Itoa(rec.statusCode)
		// r.Pattern is set by the inner mux to the matched config path, which keeps
		// label cardinality bounded by the config rather than by requested URLs.
		route := r.Pattern
		if route == "" {
			route = "unmatched"
		}
		responseTimeHistogram.WithLabelValues(route, r.Method, statusCode).Observe(duration.Seconds())
	})
}