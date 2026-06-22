package internal

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// Colors
var purple = color.RGBA{R: 225, G: 100, B: 225, A: 0xFF}
var orange = color.RGBA{R: 255, G: 150, B: 50, A: 0xFF}
var darkOrange = color.RGBA{R: 200, G: 100, B: 25, A: 0xFF}
var gray = color.RGBA{R: 150, G: 150, B: 150, A: 0xFF}

func renderer() fyne.WidgetRenderer {
	rect := canvas.NewRectangle(color.RGBA{R: 75, G: 200, B: 100, A: 255})
	rect.SetMinSize(fyne.NewSize(300, 300))
	return widget.NewSimpleRenderer(rect)
}

func offsetColor(orig color.RGBA, offset int) color.RGBA {
	colorOffset := uint8(offset) * 10
	colorRed := orig.R
	if colorRed > colorOffset {
		colorRed -= colorOffset
	}
	colorGreen := orig.G
	if colorGreen > colorOffset {
		colorGreen -= colorOffset
	}
	colorBlue := orig.B
	if colorBlue > colorOffset {
		colorBlue -= colorOffset
	}
	return color.RGBA{
		R: colorRed,
		G: colorGreen,
		B: colorBlue,
		A: orig.A,
	}
}
