//go:build js

package main

import (
	"fmt"
	"kd-tree/internal"
	"syscall/js"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())
	w := a.NewWindow("Kd Tree Nearest Neighbor")

	fContainer := container.NewWithoutLayout()

	label := widget.NewLabel("Move the mouse inside the window")
	tracker := internal.NewMouseTracker(label, fContainer)

	// Put the tracker behind everything and expand to fill
	content := container.NewMax(tracker, fContainer, container.NewVBox(label))
	w.SetContent(content)

	internal.WindowWidth = int64(js.Global().Get("innerWidth").Int())
	internal.WindowHeight = int64(js.Global().Get("innerHeight").Int())
	fmt.Printf("windowWidth: %d, windowHeight: %d\n", internal.WindowWidth, internal.WindowHeight)

	// Generate random points
	pointSet := internal.GenerateRandomPoints()

	// Draw all points
	internal.DrawPoints(fContainer, pointSet)

	kdTree := internal.KdTree{}
	tracker.SetPointTree(&kdTree)
	kdTree.Populate(pointSet)

	w.ShowAndRun()
}
