//go:build js

package main

import (
	"fmt"
	"koch-snowflake/internal"
	"syscall/js"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())
	w := a.NewWindow("Koch Snowflake")

	fContainer := container.NewWithoutLayout()
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

	ticker := time.NewTicker(time.Second)

	go func() {
		// Step 0 -> 4, then back to 0 (0 == triangle)
		step := 0

		for {
			select {
			case <-ticker.C:
				internal.InitSnowflake(fContainer, step, windowWidth, windowHeight)
				step++
				if step == 5 {
					step = 0
				}
			}
		}
	}()

	w.ShowAndRun()
}
