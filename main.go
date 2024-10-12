package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"redirection/configparsing"

	"github.com/rs/cors"
)

func redirectWithEnv(mux *http.ServeMux, targetUrl string) {

	showWarning := os.Getenv("REDIRECTION_WARNING")


	if showWarning == "" {
		fmt.Printf("Show warning is NOT SET \n")
		mux.Handle("/", http.RedirectHandler(targetUrl, http.StatusTemporaryRedirect))

	} else {
		fmt.Printf("Show warning is set to %s\n", showWarning)

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

	targetUrl := os.Getenv("REDIRECTION_TARGET_URL")
	

	if targetUrl != "" {
		fmt.Printf("Using configuration from env, target URL is %s \n", targetUrl)
		redirectWithEnv(mux, targetUrl)
	} else {
		fmt.Printf("Using configuration from config.yml \n")
		redirectWithYamlConfig(mux, targetUrl)
	}

	muxWithCors := cors.Default().Handler(mux)

	fmt.Println("[HTTP] Server listening on port :7070")
	http.ListenAndServe(":7070", muxWithCors)
	
}