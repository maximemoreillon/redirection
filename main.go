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

func redirectWithEnv(mux *http.ServeMux, targetUrl string) {

	showWarning := os.Getenv("REDIRECTION_WARNING")


	if showWarning == "" {
		fmt.Printf("Show warning is NOT SET \n")
		mux.Handle("/", http.RedirectHandler(targetUrl, http.StatusTemporaryRedirect))

	} else {
		fmt.Printf("Show warning is set\n")

		warningHandler := http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
			index(targetUrl + r.URL.String()).Render(context.Background(), w)
		})

		mux.Handle("/", warningHandler)

	}
}

func redirectWithYamlConfig(mux *http.ServeMux, targetUrl string) {
	config := configparsing.ParseConfigFile()

	for _, config := range config {
		fmt.Printf("Registering redirection for %s to %s\n", config.Path, config.Target)

		path := config.Path
		target := config.Target
		warn := config.Warn

		if warn {
			warningHandler := http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
				index(target + r.URL.String()).Render(context.Background(), w)
			})
			mux.Handle(path, warningHandler)
	
		} else {
			
			mux.Handle("/", http.RedirectHandler(targetUrl, http.StatusTemporaryRedirect))
			
		}
	}
}


func main() {

	fmt.Println("Redirection service")

	mux := http.NewServeMux()
	redirectionMux := http.NewServeMux()

	targetUrl := os.Getenv("REDIRECTION_TARGET_URL")
	

	if targetUrl != "" {
		fmt.Printf("Using configuration from env, target URL is %s \n", targetUrl)
		redirectWithEnv(redirectionMux, targetUrl)
	} else {
		fmt.Printf("Using configuration from config.yml \n")
		redirectWithYamlConfig(redirectionMux, targetUrl)
	}

	if os.Getenv("REDIRECTION_EXPORT_METRICS") != "" {
		fmt.Printf("Exporting Prometheus metrics\n")
		mux.Handle("/metrics", promhttp.Handler())
		instrumentedMux := instrumentation.MeasureResponseDuration(redirectionMux)
		mux.Handle("/", instrumentedMux)
	} else {
		mux.Handle("/", redirectionMux)
	}

	
	muxWithCors := cors.Default().Handler(mux)

	port := getPort()
	fmt.Printf("[HTTP] Server listening on port %s\n",port)
	http.ListenAndServe(port, muxWithCors)
	
}