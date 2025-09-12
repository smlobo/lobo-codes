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

	mouseLabel := widget.NewLabel("Move the mouse inside the window")
	mouseTracker := internal.NewMouseTracker(mouseLabel, fContainer)
	touchLabel := widget.NewLabel("Touch inside the window")
	touchTracker := internal.NewTouchTracker(touchLabel)

	// Put the mouseTracker & touchTracker behind everything and expand to fill
	content := container.NewMax(mouseTracker, touchTracker, fContainer, container.NewVBox(mouseLabel, touchLabel))
	w.SetContent(content)

	internal.WindowWidth = int64(js.Global().Get("innerWidth").Int())
	internal.WindowHeight = int64(js.Global().Get("innerHeight").Int())
	fmt.Printf("windowWidth: %d, windowHeight: %d\n", internal.WindowWidth, internal.WindowHeight)

	// Generate random points
	pointSet := internal.GenerateRandomPoints()

	// Draw all points
	internal.DrawPoints(fContainer, pointSet)

	kdTree := internal.KdTree{}
	mouseTracker.SetPointTree(&kdTree)
	kdTree.Populate(pointSet)
	kdTree.Draw(fContainer)

	w.ShowAndRun()
}
