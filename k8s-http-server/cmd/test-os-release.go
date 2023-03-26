package main

import (
	"fmt"
	"k8s-http-server/internal"
)

func main() {
	internal.InitOsRelease("os-release")
	fmt.Println("Name: " + internal.OsRelease.Name)
	fmt.Println("Id: " + internal.OsRelease.Id)
	fmt.Println("Version Id: " + internal.OsRelease.VersionId)
}
