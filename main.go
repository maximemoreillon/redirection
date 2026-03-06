package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"redirection/configparsing"
	"redirection/instrumentation"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

func getPort() string {
	envPort := os.Getenv("PORT")
	if envPort != "" { 
		return fmt.Sprintf(":%s", envPort) 
	} else {
		return ":80"
	}
}


func registerConfigToMux(mux *http.ServeMux, config configparsing.Config) {

	for _, config := range config {
		fmt.Printf("Registering redirection for %s to %s \n", config.Path, config.Target)

		
		handler := http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {

			userAgent := r.Header.Get("User-Agent")
			isBrowser := strings.Contains(userAgent, "Mozilla") && (strings.Contains(userAgent, "Chrome") || strings.Contains(userAgent, "Safari") || strings.Contains(userAgent, "Firefox"))

			if config.Warn && isBrowser {
				index(config.Target + r.URL.String()).Render(context.Background(), w)
			} else {
				// TODO: figure out what redirect to use
				// 301 StatusMovedPermanently
				// 307 StatusTemporaryRedirect (preserves method)
				// 308 StatusPermanentRedirect (preserves method)
				http.Redirect(w, r, config.Target, http.StatusTemporaryRedirect)
			}
		})
		mux.Handle(config.Path, handler)
	}
}



func main() {

	fmt.Println("Redirection service")

	mux := http.NewServeMux()
	redirectionMux := http.NewServeMux()

	config := configparsing.ParseConfig()

	registerConfigToMux(redirectionMux, config)
	
	if os.Getenv("REDIRECTION_EXPORT_METRICS") != "" {
		fmt.Printf("Exporting Prometheus metrics\n")
		mux.Handle("/metrics", promhttp.Handler())
		mux.Handle("/", instrumentation.MeasureResponseDuration(redirectionMux))
	} else {
		mux.Handle("/", redirectionMux)
	}

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	muxWithCors := cors.Default().Handler(mux)

	port := getPort()
	fmt.Printf("[HTTP] Server listening on port %s\n",port)
	http.ListenAndServe(port, muxWithCors)
	
}