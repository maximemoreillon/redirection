package main

import (
	"fmt"
	"net/http"
	"os"
	"redirection/configparsing"
	"redirection/instrumentation"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)


func getPort() string {

	envPort := os.Getenv("PORT")
	if len(envPort) != 0 { return fmt.Sprintf(":%s", envPort) }
	return ":80"
}

func main() {
	
	godotenv.Load()
	
	config := configparsing.ParseConfigFile()

	// Main mux, serving metrics and subMux
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	// Sub Mux for redirection requests, metrics logged with Prometheus
	subMux := http.NewServeMux()
	for _, config := range config {
		fmt.Println(fmt.Sprintf("Registering redirection for %s to %s", config.Path, config.Target))

		path := config.Path
		target := config.Target

		// Not Instrumented mux because otherwise multiple registration
		subMux.Handle(path, http.RedirectHandler(target, 307))
	}

	instrumentedMux := instrumentation.MeasureResponseDuration(subMux)
	mux.Handle("/", instrumentedMux)
	corsHandler := cors.Default().Handler(mux)

	http.ListenAndServe(":7070", corsHandler)
	fmt.Println("Redirection service listening on port :7070")
	
}