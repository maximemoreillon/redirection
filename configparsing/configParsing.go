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
	Warn bool
}

type Config []ConfigItem

func parseConfigFile () Config {
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

	return config
}

func ParseConfig () Config {

	envTargetUrl := os.Getenv("REDIRECTION_TARGET_URL")
	
	if envTargetUrl != "" {
		envWarn := os.Getenv("REDIRECTION_WARN")
		fmt.Printf("Using configuration from env, target URL is %s \n", envTargetUrl)
		config := ConfigItem{Path: "/", Target: envTargetUrl, Warn: envWarn != ""}
		return Config{config}
	} else {
		fmt.Printf("Using configuration from config.yml \n")
		return parseConfigFile()
	}

}