package configparsing

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)


type Config []struct {
	Path string
	Target string
	Warn bool
}


func ParseConfigFile () Config {
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