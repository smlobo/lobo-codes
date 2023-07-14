package main

import (
	"flag"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"n-body/internal"
	"os"
)

const size = 800

func main() {
	// Input arguments
	bgPtr := flag.Bool("bg", false, "star field background")
	timePtr := flag.Float64("time", 157788000.0, "total time")
	intervalPtr := flag.Float64("interval", 25000.0, "time interval")

	flag.Parse()

	// Command line args
	fmt.Printf("Time: %e, interval: %e\n", *timePtr, *intervalPtr)

	// Read the n-body file
	internal.ParseUniverse(os.Stdin)
	//internal.PrintUniverse()

	internal.Scale = size / (2 * internal.Radius)

	// Initialize fyne
	fApp := app.New()
	fWindow := fApp.NewWindow("N-body")
	windowSize := fyne.NewSize(size, size)
	fWindow.Resize(windowSize)
	fContainer := container.NewWithoutLayout()
	fWindow.SetContent(fContainer)

	go func() {
		// Draw the background
		if *bgPtr {
			fWindow.SetTitle("Universe")
			backgroundImage := canvas.NewImageFromFile("nbody/starfield.jpg")
			backgroundImage.Resize(windowSize)
			fContainer.Add(backgroundImage)
		}

		internal.Simulate(fContainer, *timePtr, *intervalPtr)
		internal.PrintUniverse()
	}()

	fWindow.ShowAndRun()
}
