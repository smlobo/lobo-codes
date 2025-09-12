package internal

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// Custom widget that handles taps and drags
type TouchTracker struct {
	widget.BaseWidget
	label *widget.Label
}

func NewTouchTracker(label *widget.Label) *TouchTracker {
	t := &TouchTracker{
		label: label,
	}
	t.ExtendBaseWidget(t)
	return t
}

// Implement fyne.Tappable
func (t *TouchTracker) Tapped(ev *fyne.PointEvent) {
	t.label.SetText(fmt.Sprintf("Tapped at: x=%.1f, y=%.1f\n", ev.Position.X, ev.Position.Y))
}

// Implement fyne.Draggable
func (t *TouchTracker) Dragged(ev *fyne.DragEvent) {
	t.label.SetText(fmt.Sprintf("Dragging: ΔX=%.1f, ΔY=%.1f\n", ev.Dragged.DX, ev.Dragged.DY))
}

func (t *TouchTracker) DragEnd() {
	t.label.SetText("Drag ended")
}

// Renderer
func (t *TouchTracker) CreateRenderer() fyne.WidgetRenderer {
	return renderer()
}
