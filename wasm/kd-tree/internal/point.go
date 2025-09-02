package internal

import (
	"fmt"
	"image/color"
	"math/rand/v2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Point struct {
	x, y float32
}

func (p *Point) String() string {
	return fmt.Sprintf("(%f, %f)", p.x, p.y)
}

func (p *Point) Key() int64 {
	return int64(p.x)/delta + int64(p.y)/delta*maxX
}

func (p *Point) draw(c *fyne.Container) {
	// Create a filled circle
	circle := canvas.Circle{
		Position1:   fyne.NewPos(p.x-pointDiameter/2, p.y-pointDiameter/2),
		Position2:   fyne.NewPos(p.x+pointDiameter/2, p.y+pointDiameter/2),
		FillColor:   color.Color(color.RGBA{R: 30, G: 129, B: 176, A: 0xFF}),
		StrokeColor: color.Color(color.RGBA{R: 6, G: 57, B: 112, A: 0xFF}),
		StrokeWidth: pointOutline,
	}
	c.Add(&circle)
}

func (p *Point) distanceSquaredTo(q *Point) float32 {
	dx := p.x - q.x
	dy := p.y - q.y
	return dx*dx + dy*dy
}

func GenerateRandomPoints() map[int64]Point {
	pointSet := map[int64]Point{}

	maxX = WindowWidth / delta
	maxY = WindowHeight / delta

	// Generate 20 random points
	for len(pointSet) < 20 {
		point := Point{float32(rand.Int64N(WindowWidth)), float32(rand.Int64N(WindowHeight))}
		if _, ok := pointSet[point.Key()]; !ok {
			fmt.Printf("[%d] New %v\n", len(pointSet), point)
			pointSet[point.Key()] = point
		} else {
			fmt.Printf("[%d] Existing %v\n", len(pointSet), point)
		}
	}

	return pointSet
}

func DrawPoints(c *fyne.Container, pointSet map[int64]Point) {
	for _, point := range pointSet {
		point.draw(c)
	}
	c.Refresh()
}
