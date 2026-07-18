//go:build js

package main

import (
	"fmt"
	"graph-reachability/internal"
	"syscall/js"

	"fyne.io/fyne/v2/theme"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())
	w := a.NewWindow("Directed Graph Reachability")

	fContainer := container.NewWithoutLayout()

	mouseLabel := widget.NewLabel("Move the mouse to a graph vertex")
	mouseTracker := internal.NewMouseTracker(mouseLabel, fContainer)
	touchLabel := widget.NewLabel("Touch a graph vertex")
	touchTracker := internal.NewTouchTracker(touchLabel, fContainer)

	// Put the mouseTracker & touchTracker behind everything and expand to fill
	content := container.NewStack(mouseTracker, touchTracker, fContainer, container.NewVBox(mouseLabel, touchLabel))
	w.SetContent(content)

	internal.WindowWidth = js.Global().Get("innerWidth").Int()
	internal.WindowHeight = js.Global().Get("innerHeight").Int()
	internal.MaxX = (internal.WindowWidth - internal.Margin) / internal.Delta
	internal.MaxY = (internal.WindowHeight - internal.Margin) / internal.Delta
	fmt.Printf("windowWidth: %d, windowHeight: %d\n", internal.WindowWidth, internal.WindowHeight)
	fmt.Printf("maxX: %d, maxY: %d\n", internal.MaxX, internal.MaxY)

	// Generate a random directed graph
	g := internal.NewDirectedGraph()

	// Draw the graph
	g.Draw(fContainer)

	mouseTracker.SetGraph(g)
	touchTracker.SetGraph(g)

	w.ShowAndRun()
}
