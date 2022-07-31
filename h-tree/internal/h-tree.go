package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"image/color"
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

func HTree(c *fyne.Container, depth int, size int, centerX int, centerY int) {
	// Base case
	if depth == 0 {
		return
	}

	hColor := randomColor()

	// Draw H to size and centered
	// Horizontal
	hLine := canvas.Line{
		Position1:   fyne.Position{float32(centerX - size/2), float32(centerY)},
		Position2:   fyne.Position{float32(centerX + size/2), float32(centerY)},
		StrokeColor: hColor,
		StrokeWidth: 2,
	}
	c.Add(&hLine)
	// Vertical left
	vLineLeft := canvas.Line{
		Position1:   fyne.Position{float32(centerX - size/2), float32(centerY - size/2)},
		Position2:   fyne.Position{float32(centerX - size/2), float32(centerY + size/2)},
		StrokeColor: hColor,
		StrokeWidth: 2,
	}
	c.Add(&vLineLeft)
	// Vertical right
	vLineRight := canvas.Line{
		Position1:   fyne.Position{float32(centerX + size/2), float32(centerY - size/2)},
		Position2:   fyne.Position{float32(centerX + size/2), float32(centerY + size/2)},
		StrokeColor: hColor,
		StrokeWidth: 2,
	}
	c.Add(&vLineRight)

	c.Show()
	c.Refresh()
	time.Sleep(time.Millisecond * 10)

	// Draw 4 H's at the 4 corners
	HTree(c, depth-1, size/2, int(vLineLeft.Position1.X), int(vLineLeft.Position1.Y))
	HTree(c, depth-1, size/2, int(vLineRight.Position1.X), int(vLineLeft.Position1.Y))
	HTree(c, depth-1, size/2, int(vLineLeft.Position2.X), int(vLineLeft.Position2.Y))
	HTree(c, depth-1, size/2, int(vLineRight.Position2.X), int(vLineRight.Position2.Y))

}
