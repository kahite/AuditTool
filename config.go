package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var configFileName = "config.yml"

// MizarConfig Structure of the mizar configuration
type MizarConfig struct {
	Host     string
	User     string
	Password string
}

// Queries Structure for the queries
type Queries struct {
	Count []string `yaml:",flow"`
}

// ConfigParameter Structure of the config file
type ConfigParameter struct {
	Mizar   MizarConfig
	Queries Queries
}

func getConfigFileContent() []byte {
	content, err := ioutil.ReadFile(configFileName)

	if err != nil {
		log.Fatal(err)
	}

	return content
}

func readConf() ConfigParameter {
	config := ConfigParameter{}

	content := getConfigFileContent()

	err := yaml.Unmarshal(content, &config)

	if err != nil {
		log.Fatal(err)
	}

	return config
}
