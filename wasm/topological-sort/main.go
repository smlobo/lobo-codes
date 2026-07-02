//go:build js

package main

import (
	"fmt"
	"syscall/js"
	"topological-sort/internal"

	"fyne.io/fyne/v2/theme"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())
	w := a.NewWindow("Acyclic Directed Graph Topological Sort")

	fContainer := container.NewWithoutLayout()

	//mouseLabel := widget.NewLabel("Move the mouse to a graph vertex")
	//mouseTracker := internal.NewMouseTracker(mouseLabel, fContainer)
	touchLabel := widget.NewLabel("Click/touch a graph vertex")
	touchTracker := internal.NewTouchTracker(touchLabel, fContainer)

	// Put the mouseTracker & touchTracker behind everything and expand to fill
	//content := container.NewStack(mouseTracker, touchTracker, fContainer, container.NewVBox(mouseLabel, touchLabel))
	content := container.NewStack(touchTracker, fContainer, container.NewVBox(touchLabel))
	w.SetContent(content)

	internal.WindowWidth = js.Global().Get("innerWidth").Int()
	internal.WindowHeight = js.Global().Get("innerHeight").Int()
	if internal.WindowWidth > internal.WindowHeight {
		internal.TopologicalOrientation = internal.Vertical
		internal.GraphWidth = internal.WindowWidth - internal.TopologicalReserved
		internal.GraphHeight = internal.WindowHeight
	} else {
		internal.TopologicalOrientation = internal.Horizontal
		internal.GraphWidth = internal.WindowWidth
		internal.GraphHeight = internal.WindowHeight - internal.TopologicalReserved
	}
	fmt.Printf("windowWidth: %d, windowHeight: %d, graphWidth: %d, graphHeight: %d, topoRes: %d\n",
		internal.WindowWidth, internal.WindowHeight, internal.GraphWidth, internal.GraphHeight,
		internal.TopologicalReserved)

	// Generate a random directed graph
	g := internal.NewDirectedAcyclicGraph()

	// Draw the graph
	g.Draw(fContainer)

	//mouseTracker.SetGraph(g)
	touchTracker.SetGraph(g)

	w.ShowAndRun()
}
