package internal

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"math"
)

const G = 6.67e-11

type body struct {
	xCoord    float64
	yCoord    float64
	xVelocity float64
	yVelocity float64
	mass      float64
	imageFile string
	xForce    float64
	yForce    float64
	image     *canvas.Image
}

func (b body) String() string {
	return fmt.Sprintf("<%e> <%e> <%e> <%e> <%e> %s", b.xCoord, b.yCoord, b.xVelocity,
		b.yVelocity, b.mass, b.imageFile)
}

func (b *body) draw(c *fyne.Container) {
	if b.image == nil {
		b.image = canvas.NewImageFromFile(fmt.Sprintf("nbody/%s", b.imageFile))
		b.image.FillMode = canvas.ImageFillOriginal
		//b.image.SetMinSize(fyne.NewSize(10, 10))
		b.image.Resize(fyne.NewSize(10, 10))
		//fmt.Printf("Created body: %v\n", b.circle)
		c.Add(b.image)
	}

	canvasX := 400 + b.xCoord*Scale
	canvasY := 400 - b.yCoord*Scale

	b.image.Move(fyne.NewPos(float32(canvasX), float32(canvasY)))
}

func (b *body) updatePositions(t float64) {
	b.xCoord = b.xCoord + b.xVelocity*t
	b.yCoord = b.yCoord + b.yVelocity*t
}

func (b *body) force(other *body) (xF, yF float64) {
	deltaX := other.xCoord - b.xCoord
	deltaY := other.yCoord - b.yCoord
	rSq := deltaX*deltaX + deltaY*deltaY
	r := math.Sqrt(rSq)
	f := G * b.mass * other.mass / rSq
	xF = f * deltaX / r
	yF = f * deltaY / r
	return
}
