package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"h-tree/internal"
	"os"
	"strconv"
)

func main() {
	var (
		depth int
		err   error
	)

	const windowSize = 800

	// Usage
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <depth>\n", os.Args[0])
		os.Exit(1)
	}

	if depth, err = strconv.Atoi(os.Args[1]); err != nil {
		fmt.Printf("Usage: %s <depth> (where depth is an int)\n", os.Args[0])
		os.Exit(1)
	}
	fmt.Printf("H-Tree depth: %d\n", depth)

	// Initialize fyne
	fApp := app.New()
	fWindow := fApp.NewWindow("H-Tree")
	fWindow.Resize(fyne.NewSize(windowSize, windowSize))
	fContainer := container.NewWithoutLayout()
	fWindow.SetContent(fContainer)

	// Draw the h-tree (in the background)
	go internal.HTree(fContainer, depth, windowSize/2, windowSize/2, windowSize/2)

	fWindow.ShowAndRun()
}
