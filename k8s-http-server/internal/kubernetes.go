package internal

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

func InitKubernetesInfo() {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Printf("Error getting k8s cluster config: %s", err.Error())
		return
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Error getting k8s clientset: %s", err.Error())
	}

	serverVersion, err := clientset.DiscoveryClient.ServerVersion()
	log.Printf("Server version: %s", serverVersion.GitVersion)
}
