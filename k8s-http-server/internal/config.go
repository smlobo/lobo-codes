package internal

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var Config map[string]string

func SetupConfig() {
	// Open the Config file
	configFile, err := os.Open("Config/k8s-http-server.yaml")
	if err != nil {
		fmt.Println("Error opening Config file:", err)
		os.Exit(1)
	}
	defer configFile.Close()

	// Unmarshal into the map
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&Config)
	if err != nil {
		fmt.Println("Error unmarshalling Config file:", err)
		os.Exit(1)
	}

	// Override with environment variables
	for variable, _ := range Config {
		envVarSet := os.Getenv(variable)
		if len(envVarSet) != 0 {
			Config[variable] = envVarSet
			log.Println("From environment: [", variable, "] =", envVarSet)
		}
	}
}
