package internal

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type OsReleaseInfo struct {
	Name      string
	Id        string
	VersionId string
}

var OsRelease OsReleaseInfo

func InitOsRelease(fileName string) {
	readFile, err := os.Open(fileName)

	if err != nil {
		log.Printf("could not open os-release file: %s; %s", fileName, err.Error())
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		lineTokens := strings.Split(fileScanner.Text(), "=")
		if len(lineTokens) != 2 {
			continue
		}
		if lineTokens[0] == "NAME" {
			OsRelease.Name = lineTokens[1]
		} else if lineTokens[0] == "ID" {
			OsRelease.Id = lineTokens[1]
		} else if lineTokens[0] == "VERSION_ID" {
			OsRelease.VersionId = lineTokens[1]
		}
	}

	_ = readFile.Close()
}

func (osri OsReleaseInfo) String() string {
	return fmt.Sprintf("<%s,%s,%s>", osri.Name, osri.Id, osri.VersionId)
}
