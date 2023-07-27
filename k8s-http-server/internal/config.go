package internal

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var config map[string]string

func SetupConfig() {
	// Open the config file
	configFile, err := os.Open("config/k8s-http-server.yaml")
	if err != nil {
		fmt.Println("Error opening config file:", err)
		os.Exit(1)
	}
	defer configFile.Close()

	// Unmarshal into the map
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error unmarshalling config file:", err)
		os.Exit(1)
	}

	// Override with environment variables
	for variable, _ := range config {
		envVarSet := os.Getenv(variable)
		if len(envVarSet) != 0 {
			config[variable] = envVarSet
			log.Println("From environment: [", variable, "] =", envVarSet)
		}
	}
}
