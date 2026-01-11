package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// Custom widget that handles taps and drags
type TouchTracker struct {
	widget.BaseWidget
}

func NewTouchTracker() *TouchTracker {
	t := &TouchTracker{}
	t.ExtendBaseWidget(t)
	return t
}

// Implement fyne.Tappable
func (t *TouchTracker) Tapped(ev *fyne.PointEvent) {
	if !Paused {
		Paused = true
	} else {
		Paused = false
	}
}

// Implement fyne.Draggable
func (t *TouchTracker) Dragged(ev *fyne.DragEvent) {
	// Optional: handle drag
}

func (t *TouchTracker) DragEnd() {
	// Optional: handle drag end
}

// Renderer
func (t *TouchTracker) CreateRenderer() fyne.WidgetRenderer {
	return renderer()
}
