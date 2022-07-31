package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"n-body/internal"
	"os"
	"strconv"
)

const size = 800

func main() {
	// Usage
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <simulation-time> <time-interval>\n", os.Args[0])
		os.Exit(1)
	}

	// Command line args
	totalTime, _ := strconv.ParseFloat(os.Args[1], 64)
	timeInterval, _ := strconv.ParseFloat(os.Args[2], 64)
	fmt.Printf("Time: %e, interval: %e\n", totalTime, timeInterval)

	// Read the n-body file
	internal.ParseUniverse(os.Stdin)
	//internal.PrintUniverse()

	internal.Scale = size / (2 * internal.Radius)

	// Initialize fyne
	fApp := app.New()
	fWindow := fApp.NewWindow("Universe")
	windowSize := fyne.NewSize(size, size)
	fWindow.Resize(windowSize)
	fContainer := container.NewWithoutLayout()
	fWindow.SetContent(fContainer)

	go func() {
		// Draw the background
		backgroundImage := canvas.NewImageFromFile("nbody/starfield.jpg")
		backgroundImage.Resize(windowSize)
		fContainer.Add(backgroundImage)

		internal.Simulate(fContainer, totalTime, timeInterval)
		internal.PrintUniverse()
	}()

	fWindow.ShowAndRun()
}
