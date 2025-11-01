package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/adrg/frontmatter"
)

const (
	contentDir   = "./content/posts"
	markdownName = "index.md"
	locationJson = "static/map/locations.json"
)

type locationData struct {
	Title  string    `toml:"title"`
	Coords []float32 `toml:"coords"`
}

func main() {
	// Iterate over all content markdown
	entries, err := os.ReadDir(contentDir)
	if err != nil {
		log.Fatalf("Cannot find %s {%v}", contentDir, err)
	}

	locations := []locationData{}

	for _, e := range entries {
		// content/post dirs only
		if !e.IsDir() {
			continue
		}
		fmt.Printf("Dir : %s\n", e.Name())
		indexPath := fmt.Sprintf("%s/%s/%s", contentDir, e.Name(), markdownName)
		if _, err := os.Stat(indexPath); err != nil {
			log.Printf("Could not find: %s", indexPath)
			continue
		}

		// Parse file
		if location, err := parseMarkdown(indexPath); err == nil {
			locations = append(locations, location)
		} else {
			log.Fatalf("Error parsing %s", indexPath)
		}
	}

	// Write to json
	if err = writeLocationJson(locations); err != nil {
		log.Fatalf("Error writing locations.json: (%v)", err)
	}
}

func parseMarkdown(indexPath string) (locationData, error) {
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return locationData{}, err
	}

	var matter locationData

	_, err = frontmatter.Parse(strings.NewReader(string(data)), &matter)
	if err != nil || matter.Title == "" || len(matter.Coords) != 2 {
		return locationData{}, err
	}

	fmt.Printf("    Title: %s\n", matter.Title)
	fmt.Printf("    Lat:   %.4f\n", matter.Coords[0])
	fmt.Printf("    Long:  %.4f\n", matter.Coords[1])

	return matter, nil
}

func writeLocationJson(locations []locationData) error {
	f, _ := os.Create(locationJson)
	defer f.Close()

	writer := bufio.NewWriter(f)
	writer.WriteString("[\n")
	for index, location := range locations {
		if index != 0 {
			writer.WriteString(",\n")
		}
		// Skip locations not available
		if len(location.Coords) != 2 {
			continue
		}
		writer.WriteString("  { \"name\": \"" + location.Title + "\", \"coords\": [")
		writer.WriteString(fmt.Sprintf("%.4f, %.4f", location.Coords[1], location.Coords[0]) + "] }")
	}
	writer.WriteString("\n]\n")
	writer.Flush()

	return nil
}
