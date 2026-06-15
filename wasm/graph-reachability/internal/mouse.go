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
	graph     *DirectedGraph
	container *fyne.Container
}

func NewMouseTracker(label *widget.Label, container *fyne.Container) *MouseTracker {
	m := &MouseTracker{
		label:     label,
		container: container,
	}
	m.ExtendBaseWidget(m)
	return m
}

func (m *MouseTracker) SetGraph(g *DirectedGraph) {
	m.graph = g
}

// Implement desktop.Mouseable to get mouse move events
func (m *MouseTracker) MouseMoved(ev *desktop.MouseEvent) {
	m.label.SetText(fmt.Sprintf("%.1f, %.1f", ev.Position.X, ev.Position.Y))

	point := QuantizedPoint{
		x: quantized(int64(ev.Position.X)),
		y: quantized(int64(ev.Position.Y)),
	}
	m.graph.HighlightReachableAt(point)
	m.graph.Redraw(m.container)
}

func (m *MouseTracker) MouseIn(ev *desktop.MouseEvent) {
	// Optional: handle mouse enter
}

func (m *MouseTracker) MouseOut() {
	// Optional: handle mouse leave
	m.label.SetText("Mouse out")
}

// Needed to render something visible
func (m *MouseTracker) CreateRenderer() fyne.WidgetRenderer {
	//rect := canvas.NewRectangle(color.NRGBA{R: 200, G: 200, B: 255, A: 255})
	return renderer()
}
