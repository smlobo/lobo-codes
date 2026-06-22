package internal

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// TouchTracker Custom widget that handles taps and drags
type TouchTracker struct {
	widget.BaseWidget
	label     *widget.Label
	graph     *DirectedGraph
	previous  QuantizedPoint
	container *fyne.Container
}

func NewTouchTracker(label *widget.Label, container *fyne.Container) *TouchTracker {
	t := &TouchTracker{
		label:     label,
		container: container,
	}
	t.ExtendBaseWidget(t)
	return t
}

func (m *TouchTracker) SetGraph(g *DirectedGraph) {
	m.graph = g
}

// Tapped Implement fyne.Tappable
func (t *TouchTracker) Tapped(ev *fyne.PointEvent) {
	point := QuantizedPoint{
		x: quantized(int64(ev.Position.X)),
		y: quantized(int64(ev.Position.Y)),
	}

	if point != t.previous {
		if t.graph.HighlightCycleAt(point) {
			t.label.SetText(fmt.Sprintf("Click/touch again for more cycles (%.1f, %.1f)",
				ev.Position.X, ev.Position.Y))
		} else {
			t.label.SetText(fmt.Sprintf("(%.1f, %.1f)", ev.Position.X, ev.Position.Y))
		}
		t.graph.Redraw(t.container)
		t.previous = point
	} else {
		if t.graph.HighlightNextCycleAt(point) {
			t.label.SetText(fmt.Sprintf("Click/touch again for more cycles (%.1f, %.1f)",
				ev.Position.X, ev.Position.Y))
			t.graph.Redraw(t.container)
		}
	}
}

// Dragged Implement fyne.Draggable
func (t *TouchTracker) Dragged(ev *fyne.DragEvent) {
	t.label.SetText(fmt.Sprintf("Dragging: ΔX=%.1f, ΔY=%.1f\n", ev.Dragged.DX, ev.Dragged.DY))
}

func (t *TouchTracker) DragEnd() {
	t.label.SetText("Drag ended")
}

// CreateRenderer Renderer
func (t *TouchTracker) CreateRenderer() fyne.WidgetRenderer {
	return renderer()
}
