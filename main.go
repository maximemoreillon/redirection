package main

import (
	"fmt"
	"net/http"
	"redirection/configparsing"
	"redirection/instrumentation"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)



func main() {
	
	// Main mux, serving metrics and subMux
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	// Sub Mux for redirection requests, metrics logged with Prometheus
	subMux := http.NewServeMux()
	wrappedMux := instrumentation.MeasureResponseDuration(subMux)
	
	
	config := configparsing.ParseConfigFile()

	for _, config := range config {
		fmt.Println(fmt.Sprintf("Registering redirection for %s to %s", config.Path, config.Target))

		path := config.Path
		target := config.Target

		subMux.Handle(path, http.RedirectHandler(target, 307))
	}

	fmt.Println("Redirection service listening on port :7070")
	mux.Handle("/", wrappedMux)
	corsHandler := cors.Default().Handler(mux)
	http.ListenAndServe(":7070", corsHandler)
	
}