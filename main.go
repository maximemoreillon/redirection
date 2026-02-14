package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"redirection/configparsing"
	"redirection/instrumentation"

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
		fmt.Printf("Registering redirection for %s to %s ", config.Path, config.Target)
		if config.Warn {
			fmt.Print("with warning page\n")
			warningHandler := http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
				index(config.Target + r.URL.String()).Render(context.Background(), w)
			})
			mux.Handle(config.Path, warningHandler)
		} else {
			fmt.Print("without warning page\n")
			mux.Handle("/", http.RedirectHandler(config.Target, http.StatusTemporaryRedirect))
		}
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