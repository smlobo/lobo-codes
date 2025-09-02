package internal

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type MouseTracker struct {
	widget.BaseWidget
	label        *widget.Label
	pointTree    *KdTree
	container    *fyne.Container
	previousLine *canvas.Line
}

func NewMouseTracker(label *widget.Label, container *fyne.Container) *MouseTracker {
	m := &MouseTracker{
		label:     label,
		container: container,
	}
	m.ExtendBaseWidget(m)
	return m
}

func (m *MouseTracker) SetPointTree(t *KdTree) {
	m.pointTree = t
}

// Implement desktop.Mouseable to get mouse move events
func (m *MouseTracker) MouseMoved(ev *desktop.MouseEvent) {
	m.label.SetText(fmt.Sprintf("Mouse at: %.0f, %.0f", ev.Position.X, ev.Position.Y))

	// Search for the nearest Point
	nearest := m.pointTree.Nearest(&Point{ev.Position.X, ev.Position.Y})

	// Draw a line to the nearest point
	if nearest == nil {
		return
	}

	m.container.Remove(m.previousLine)

	m.previousLine = &canvas.Line{
		Position1:   fyne.Position{ev.Position.X, ev.Position.Y},
		Position2:   fyne.Position{nearest.x, nearest.y},
		StrokeColor: color.RGBA{R: 0x2b, G: 0xbb, B: 0x21, A: 0xFF},
		StrokeWidth: 2,
	}
	m.container.Add(m.previousLine)
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
	rect := canvas.NewRectangle(color.RGBA{R: 251, G: 233, B: 183, A: 255})
	rect.SetMinSize(fyne.NewSize(300, 300))
	return widget.NewSimpleRenderer(rect)
}
