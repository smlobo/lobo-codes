package internal

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type MouseTracker struct {
	widget.BaseWidget
	label     *widget.Label
	graph     *DirectedAcyclicGraph
	previous  QuantizedPoint
	container *fyne.Container
}

func NewMouseTracker(label *widget.Label, container *fyne.Container) *MouseTracker {
	m := &MouseTracker{
		label:     label,
		container: container,
		previous:  QuantizedPoint{},
	}
	m.ExtendBaseWidget(m)
	return m
}

func (m *MouseTracker) SetGraph(g *DirectedAcyclicGraph) {
	m.graph = g
}

// MouseMoved Implement desktop.Mouseable to get mouse move events
func (m *MouseTracker) MouseMoved(ev *desktop.MouseEvent) {
	m.label.SetText(fmt.Sprintf("%.1f, %.1f", ev.Position.X, ev.Position.Y))
}

func (m *MouseTracker) MouseIn(ev *desktop.MouseEvent) {
	// Optional: handle mouse enter
}

func (m *MouseTracker) MouseOut() {
	// Optional: handle mouse leave
	m.label.SetText("Mouse out")
}

// CreateRenderer Needed to render something visible
func (m *MouseTracker) CreateRenderer() fyne.WidgetRenderer {
	return renderer()
}
