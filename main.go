package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"strconv"
	"path/filepath"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
)

type Config []struct {
	Path string
	Target string
}

// TODO: file splitting
//statusRecorder to record the status code from the ResponseWriter
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}


// func (rec *statusRecorder) WriteHeader(statusCode int) {
// 	rec.statusCode = statusCode
// 	rec.ResponseWriter.WriteHeader(statusCode)
// }

func measureResponseDuration(next http.Handler) http.Handler {
	buckets := []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}

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
		route := r.URL.Path
		responseTimeHistogram.WithLabelValues(route, r.Method, statusCode).Observe(duration.Seconds())
	})
}


func main() {
	
	// Main mux, serving metrics and subMux
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	// Sub Mux for redirection requests, metrics logged with Prometheus
	subMux := http.NewServeMux()
	wrappedMux := measureResponseDuration(subMux)
	
	// Parse Yaml
	filename, _ := filepath.Abs("./config/config.yml")
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	
	config := Config{}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
			panic(err)
	}

	for _, config := range config {
		fmt.Println(fmt.Sprintf("Registering redirection for %s to %s", config.Path, config.Target))

		path := config.Path
		target := config.Target

		subMux.Handle(path, http.RedirectHandler(target, 307))
	}

	fmt.Println("Redirection service listening on port :7070")
	mux.Handle("/", wrappedMux)
	http.ListenAndServe(":7070", mux)
	
}