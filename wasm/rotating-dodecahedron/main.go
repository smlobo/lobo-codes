//go:build js

package main

import (
	"fyne.io/fyne/v2/widget"
	"rotating-dodecahedron/internal"
	"syscall/js"
	"time"
	"unicode"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func main() {
	a := app.New()
	w := a.NewWindow("Rotating Dodecahedron")

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
					internal.DrawDodecahedron(fContainer)
					internal.Rotate()
				}
			}
		}
	}()

	w.ShowAndRun()
}
