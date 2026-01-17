package internal

import (
	"math"

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

func (t *TouchTracker) Tapped(ev *fyne.PointEvent) {
	if !Paused {
		Paused = true
		Rotation = ZeroRotationAngle
	} else {
		Paused = false
		Rotation = DefaultRotationAngle
	}
}

func (t *TouchTracker) Dragged(ev *fyne.DragEvent) {
	if ev.Dragged.IsZero() {
		return
	}
	Paused = true
	Rotation = ZeroRotationAngle
	Rotation.y = float32(math.Atan(float64(-ev.Dragged.DX / 20.0)))
	Rotation.x = float32(math.Atan(float64(-ev.Dragged.DY / 20.0)))
	//fmt.Printf("Dragged: %v -> %v\n", ev, Rotation)
}

func (t *TouchTracker) DragEnd() {
	Paused = true
	Rotation = ZeroRotationAngle
	//fmt.Printf("Drag End\n")
}

// Renderer
func (t *TouchTracker) CreateRenderer() fyne.WidgetRenderer {
	return renderer()
}
