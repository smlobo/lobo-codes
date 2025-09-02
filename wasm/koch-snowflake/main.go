//go:build js

package main

import (
	"fmt"
	"koch-snowflake/internal"
	"math"
	"syscall/js"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())
	w := a.NewWindow("Koch Snowflake")

	fContainer := container.NewWithoutLayout()
	//fContainer := container.NewVBox()
	w.SetContent(fContainer)

	// Not accurate sizes - need to get from browser
	canvasSize := w.Canvas().Size()
	fmt.Printf("CanvasSize: %v\n", canvasSize)
	containerSize := fContainer.Size()
	fmt.Printf("containerSize: %v\n", containerSize)

	windowWidth := js.Global().Get("innerWidth").Int()
	windowHeight := js.Global().Get("innerHeight").Int()
	dimension := min(windowWidth, windowHeight)
	fmt.Printf("windowWidth: %d, windowHeight: %d, dimension: %d\n", windowWidth, windowHeight, dimension)

	go func() {
		// Infinite loop
		for {
			// Draw each step (in the background)
			for i := 1; i <= 5; i++ {
				// Clear the previous drawing
				fContainer.RemoveAll()
				fContainer.Refresh()

				// The drawing square
				squareSize := float32(dimension - 100)

				// Offset of drawing square in the windows
				offsetX := (float32(windowWidth) - squareSize) / 2
				offsetY := (float32(windowHeight) - squareSize) / 2

				// Start the recursion with the 3 sides of the equilateral triangle
				triangleSide := float32(dimension) / 1.3

				// Triangle height
				triangleHeight := float32(math.Sqrt(math.Pow(float64(triangleSide), 2) -
					math.Pow(float64(triangleSide/2), 2)))

				//fmt.Printf("Sq: %.1f, Ts: %.1f, Th: %.1f, 1.3Th: %.1f\n", squareSize, triangleSide,
				//	triangleHeight, 1.3*triangleHeight)

				// Bottom left start position
				bottomLeftX := squareSize/2 - triangleSide/2 + offsetX
				bottomLeftY := squareSize - (squareSize-1.3*triangleHeight)/2 - .3*triangleHeight + offsetY

				// The 3 points
				bottomLeft := fyne.Position{bottomLeftX, bottomLeftY}
				bottomRight := fyne.Position{bottomLeft.X + triangleSide, bottomLeft.Y}
				topPosition := fyne.Position{(bottomRight.X-bottomLeft.X)/2 + bottomLeft.X, bottomLeft.Y - triangleHeight}

				internal.KockSnowflake(fContainer, i-1, bottomRight, bottomLeft)
				internal.KockSnowflake(fContainer, i-1, bottomLeft, topPosition)
				internal.KockSnowflake(fContainer, i-1, topPosition, bottomRight)

				// Show the current drawing
				fContainer.Show()
				fContainer.Refresh()

				// Sleep before drawing the next step
				time.Sleep(time.Second)
			}
		}
	}()

	w.ShowAndRun()
}
