package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"koch-snowflake/internal"
	"math"
	"os"
	"strconv"
	"time"
)

func main() {
	var (
		depth int
		err   error
	)

	// Usage
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <depth>\n", os.Args[0])
		os.Exit(1)
	}

	if depth, err = strconv.Atoi(os.Args[1]); err != nil {
		fmt.Printf("Usage: %s <depth> (where depth is an int)\n", os.Args[0])
		os.Exit(1)
	}
	fmt.Printf("Koch Snowflake depth: %d\n", depth)

	// Initialize fyne
	fApp := app.New()
	fWindow := fApp.NewWindow("Koch Snowflake " + os.Args[1])
	fWindow.Resize(fyne.NewSize(800, 800))
	fContainer := container.NewWithoutLayout()
	fWindow.SetContent(fContainer)

	go func() {
		// Draw each step (in the background)
		for i := 1; i <= depth; i++ {
			// Clear the previous drawing
			fContainer.RemoveAll()
			fContainer.Refresh()

			// Start the recursion with the 3 sides of the equilateral triangle
			triangleSide := float32(600)

			// The 3 points
			bottomLeft := fyne.Position{100, 575}
			bottomRight := fyne.Position{bottomLeft.X + triangleSide, bottomLeft.Y}

			// Triangle height
			h := math.Sqrt(math.Pow(float64(triangleSide), 2) - math.Pow(float64(triangleSide/2), 2))
			topPosition := fyne.Position{(bottomRight.X-bottomLeft.X)/2 + bottomLeft.X, bottomLeft.Y - float32(h)}

			internal.KockSnowflake(fContainer, i-1, bottomRight, bottomLeft)
			internal.KockSnowflake(fContainer, i-1, bottomLeft, topPosition)
			internal.KockSnowflake(fContainer, i-1, topPosition, bottomRight)

			// Show the current drawing
			fContainer.Show()
			fContainer.Refresh()

			// Sleep before drawing the next step
			time.Sleep(time.Second)
		}
	}()

	fWindow.ShowAndRun()
}
