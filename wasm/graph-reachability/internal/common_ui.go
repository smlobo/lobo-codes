package internal

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// Colors
var green = color.RGBA{R: 50, G: 200, B: 75, A: 0xFF}
var red = color.RGBA{R: 220, G: 50, B: 50, A: 0xFF}

func renderer() fyne.WidgetRenderer {
	rect := canvas.NewRectangle(color.RGBA{R: 200, G: 180, B: 225, A: 255})
	rect.SetMinSize(fyne.NewSize(300, 300))
	return widget.NewSimpleRenderer(rect)
}
