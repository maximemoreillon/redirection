package configparsing

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)


type ConfigItem struct {
	Path string
	Target string
	Warning bool
}

type Config []ConfigItem


func ParseConfig () Config {

	var config Config

	// parse from env
	envTarget := os.Getenv("TARGET")
	envWarning := os.Getenv("WARNING")

	if(len(envTarget) != 0) {
		fmt.Println("Configuration via environment variable available")
		configItem := ConfigItem{Path: "/", Target: envTarget, Warning: len(envWarning) != 0}
		
		config = append(config, configItem)
	}

	// Parse from file
	filename, _ := filepath.Abs("./config/config.yml")
	_, err := os.Stat(filename)

	if err == nil {
		fmt.Println("Configuration file available")
		yamlFile, err := os.ReadFile(filename)
		if err != nil {
			panic(err)
		}

			configFromFile := Config{}
			err = yaml.Unmarshal(yamlFile, &configFromFile)
			if err != nil {
					panic(err)
			}

			for _, configItem := range configFromFile {
				config = append(config, configItem)
			}

	} else {
			fmt.Println("No configuration file available")
	}


	return config
}
	