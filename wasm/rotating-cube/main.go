//go:build js

package main

import (
	"fyne.io/fyne/v2/widget"
	"rotating-cube/internal"
	"syscall/js"
	"time"
	"unicode"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())
	w := a.NewWindow("Rotating Cube")

	fContainer := container.NewWithoutLayout()

	instrLabel := widget.NewLabel("Tap/Press any key to toggle start/stop")
	touchTracker := internal.NewTouchTracker()
	content := container.NewStack(touchTracker, fContainer, container.NewVBox(instrLabel))
	w.SetContent(content)

	w.Canvas().SetOnTypedRune(func(r rune) {
		if unicode.IsSpace(r) || unicode.IsLetter(r) {
			if !internal.Paused {
				internal.Paused = true
			} else {
				internal.Paused = false
			}
		}
	})

	internal.Initialize(float32(js.Global().Get("innerWidth").Int()),
		float32(js.Global().Get("innerHeight").Int()))

	ticker := time.NewTicker(time.Millisecond * 25)

	go func() {
		for {
			select {
			case <-ticker.C:
				if !internal.Paused {
					internal.DrawCube(fContainer)
					internal.Rotate()
				}
			}
		}
	}()

	w.ShowAndRun()
}
