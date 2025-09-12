package internal

import (
	"fmt"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type MouseTracker struct {
	widget.BaseWidget
	label            *widget.Label
	pointTree        *KdTree
	container        *fyne.Container
	previousNearest  *canvas.Line
	previousFarthest *canvas.Line
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
	var nearestTxt, farthestTxt string

	// Search for the nearest Point
	nearest := m.pointTree.Nearest(&Point{ev.Position.X, ev.Position.Y})

	// Draw a line to the nearest point
	if nearest != nil {
		m.container.Remove(m.previousNearest)

		m.previousNearest = &canvas.Line{
			Position1:   fyne.Position{ev.Position.X, ev.Position.Y},
			Position2:   fyne.Position{nearest.x, nearest.y},
			StrokeColor: color.RGBA{R: 0x2b, G: 0xbb, B: 0x21, A: 0xff},
			StrokeWidth: 2,
		}
		m.container.Add(m.previousNearest)

		dx := ev.Position.X - nearest.x
		dy := ev.Position.Y - nearest.y
		nearestDistance := math.Sqrt(float64(dx*dx + dy*dy))
		nearestTxt = fmt.Sprintf("Nearest: %.2f\n", nearestDistance)
	}

	// Search for the farthest Point
	farthest := m.pointTree.Farthest(&Point{ev.Position.X, ev.Position.Y})

	// Draw a line to the farthest point
	if farthest != nil {
		m.container.Remove(m.previousFarthest)

		m.previousFarthest = &canvas.Line{
			Position1:   fyne.Position{ev.Position.X, ev.Position.Y},
			Position2:   fyne.Position{farthest.x, farthest.y},
			StrokeColor: color.RGBA{R: 0xbb, G: 0x2b, B: 0x21, A: 0xff},
			StrokeWidth: 2,
		}
		m.container.Add(m.previousFarthest)

		dx := ev.Position.X - farthest.x
		dy := ev.Position.Y - farthest.y
		farthestDistance := math.Sqrt(float64(dx*dx + dy*dy))
		farthestTxt = fmt.Sprintf("Farthest: %.2f", farthestDistance)
	}

	m.label.SetText(nearestTxt + farthestTxt)
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
