package internal

import (
	"fmt"
	"os"
	"strings"
)

type GoVersionInfo string

var GoVersion GoVersionInfo

func InitGoVersion(fileName string) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading golang_version.txt file: %s; %s\n", fileName, err)
	} else {

	}
	GoVersion = GoVersionInfo(strings.TrimSpace(string(content)))
}
