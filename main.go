package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"redirection/configparsing"
	"redirection/instrumentation"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

type PageData struct {
    Target  string
}


func getPort() string {
	envPort := os.Getenv("PORT")
	if len(envPort) != 0 { return fmt.Sprintf(":%s", envPort) }
	return ":80"
}




func main() {
	
	godotenv.Load()
	
	config := configparsing.ParseConfig()

	// Main mux, serving metrics and subMux
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	htmlTemplate := template.Must(template.ParseFiles("templates/index.html"))

	// Sub Mux for redirection requests, metrics logged with Prometheus
	subMux := http.NewServeMux()
	for _, config := range config {
		fmt.Println(fmt.Sprintf("Registering redirection for %s to %s", config.Path, config.Target))

		// TODO: check if WARN
		path := config.Path
		target := config.Target
		warning := config.Warning

		if (warning) {
			subMux.HandleFunc(path, func (w http.ResponseWriter, r *http.Request)  {
				templateData := PageData{
					Target:  target,
				}
				htmlTemplate.Execute(w, templateData)
			})

		} else {
			subMux.Handle(path, http.RedirectHandler(target, 307))
		}

		// Not Instrumented mux because otherwise multiple registration

	}

	// Instrumentation (Prometheus) and CORS
	instrumentedMux := instrumentation.MeasureResponseDuration(subMux)
	mux.Handle("/", instrumentedMux)
	corsHandler := cors.Default().Handler(mux)

	fmt.Println("Redirection service listening on port :7070")
	http.ListenAndServe(":7070", corsHandler)
	
}