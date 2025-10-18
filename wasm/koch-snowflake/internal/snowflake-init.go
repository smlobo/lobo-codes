package internal

import (
	"fyne.io/fyne/v2"
	"math"
)

func InitSnowflake(fContainer *fyne.Container, step, windowWidth, windowHeight int) {
	dimension := min(windowWidth, windowHeight)

	// Clear the previous drawing
	fContainer.RemoveAll()

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

	kockSnowflake(fContainer, step, bottomRight, bottomLeft)
	kockSnowflake(fContainer, step, bottomLeft, topPosition)
	kockSnowflake(fContainer, step, topPosition, bottomRight)

	// Show the current drawing
	fContainer.Show()
	fyne.Do(func() {
		fContainer.Refresh()
	})
}
