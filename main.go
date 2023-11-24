package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
)

type Config []struct {
	Path string
	Target string
}


func main() {

	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())
	
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
		mux.Handle(config.Path, http.RedirectHandler(config.Target, 307))
	}	

	fmt.Println("Redirection service listening on port :7070")
	http.ListenAndServe(":7070", mux)
	
}