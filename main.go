package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
)


type Route struct {
	Path string "yaml:Path"
	Target string "yaml:Target"
}

type Config []Route


func main() {

	http.Handle("/metrics", promhttp.Handler())
	
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
		http.HandleFunc(config.Path, func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, config.Target, 301)
		})
	}

	// TODO: create handlers in a for loop
	

	http.ListenAndServe(":7070", nil)
}