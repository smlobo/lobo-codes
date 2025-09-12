package internal

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

func renderer() fyne.WidgetRenderer {
	rect := canvas.NewRectangle(color.RGBA{R: 251, G: 233, B: 183, A: 255})
	rect.SetMinSize(fyne.NewSize(300, 300))
	return widget.NewSimpleRenderer(rect)
}
