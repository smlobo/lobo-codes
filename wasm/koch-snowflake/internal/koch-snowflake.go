package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"image/color"
	"math"
	"math/rand"
	"time"
)

func randomColor() color.Color {
	const maxRGB = 256
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	hColor := color.RGBA{
		R: uint8(random.Intn(maxRGB)),
		G: uint8(random.Intn(maxRGB)),
		B: uint8(random.Intn(maxRGB)),
		A: maxRGB - 1,
	}
	return hColor
}

func kockSnowflake(c *fyne.Container, depth int, position1, position2 fyne.Position) {
	// Draw when we get to the innermost recursion
	if depth == 0 {
		sColor := randomColor()
		line := canvas.Line{
			Position1:   position1,
			Position2:   position2,
			StrokeColor: sColor,
			StrokeWidth: 2,
		}
		c.Add(&line)

		return
	}

	// Line length
	length := math.Sqrt(math.Pow(float64(position2.X-position1.X), 2) + math.Pow(float64(position2.Y-position1.Y), 2))

	// X third positions
	totalX := position2.X - position1.X
	thirdX := totalX / 3
	oneThirdX := position1.X + thirdX
	twoThirdX := position1.X + thirdX*2

	// Y third positions
	totalY := position2.Y - position1.Y
	thirdY := totalY / 3
	oneThirdY := position1.Y + thirdY
	twoThirdY := position1.Y + thirdY*2

	// Angle
	//angle := math.Atan(float64(totalY / totalX))
	angle := math.Atan2(float64(totalY), float64(totalX))
	perpendicular := math.Pi/2 + angle

	// Mid point
	midX := position1.X + totalX/2
	midY := position1.Y + totalY/2

	// Equilateral triangle height
	tLength := length / 3
	tHeight := math.Sqrt(math.Pow(tLength, 2) - math.Pow(tLength/2, 2))

	// X & Y offsets
	xDiff := tHeight * math.Cos(perpendicular)
	yDiff := tHeight * math.Sin(perpendicular)
	tipX := midX - float32(xDiff)
	tipY := midY - float32(yDiff)

	// The 3 points of the equilateral triangle
	oneThirdPosition := fyne.Position{oneThirdX, oneThirdY}
	twoThirdPosition := fyne.Position{twoThirdX, twoThirdY}
	tipPosition := fyne.Position{tipX, tipY}

	// Draw 2 lines on this line - both ends
	kockSnowflake(c, depth-1, position1, oneThirdPosition)
	kockSnowflake(c, depth-1, twoThirdPosition, position2)

	// Draw the angled lines
	kockSnowflake(c, depth-1, oneThirdPosition, tipPosition)
	kockSnowflake(c, depth-1, tipPosition, twoThirdPosition)
}
