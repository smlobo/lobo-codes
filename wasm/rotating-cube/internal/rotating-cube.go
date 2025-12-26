package internal

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"golang.org/x/image/colornames"
)

type coordinates3D struct {
	x, y, z float32
}

type line3D struct {
	a, b *coordinates3D
}

var cubePoints [8]coordinates3D
var cubeLines [12]line3D
var centroid coordinates3D

func DrawCube(fContainer *fyne.Container) {
	// Clear the previous drawing
	fContainer.RemoveAll()

	for _, cLine := range cubeLines {
		fLine := canvas.Line{
			Position1:   fyne.Position{cLine.a.x, cLine.a.y},
			Position2:   fyne.Position{cLine.b.x, cLine.b.y},
			StrokeColor: colornames.Blue,
			StrokeWidth: 2,
		}
		fContainer.Add(&fLine)
	}

	// Show the current drawing
	fContainer.Show()
	fyne.Do(func() {
		fContainer.Refresh()
	})
}

func rotatePoint(p *coordinates3D, x, y, z float32) {
	fixed := coordinates3D{p.x, p.y, p.z}
	p.y = float32(math.Cos(float64(x)))*fixed.y - float32(math.Sin(float64(x)))*fixed.z
	p.z = float32(math.Sin(float64(x)))*fixed.y + float32(math.Cos(float64(x)))*fixed.z

	fixed = coordinates3D{p.x, p.y, p.z}
	p.x = float32(math.Cos(float64(y)))*fixed.x + float32(math.Sin(float64(y)))*fixed.z
	p.z = -float32(math.Sin(float64(y)))*fixed.x + float32(math.Cos(float64(y)))*fixed.z

	fixed = coordinates3D{p.x, p.y, p.z}
	p.x = float32(math.Cos(float64(z)))*fixed.x - float32(math.Sin(float64(z)))*fixed.y
	p.y = float32(math.Sin(float64(z)))*fixed.x + float32(math.Cos(float64(z)))*fixed.y
}

func Rotate() {
	// Rotate each point of the cube
	for i := 0; i < len(cubePoints); i++ {
		// Translate to the origin
		p := coordinates3D{
			x: cubePoints[i].x - centroid.x,
			y: cubePoints[i].y - centroid.y,
			z: cubePoints[i].z - centroid.z,
		}
		rotatePoint(&p, 0.02, 0.01, 0.03)
		cubePoints[i].x = p.x + centroid.x
		cubePoints[i].y = p.y + centroid.y
		cubePoints[i].z = p.z + centroid.z
	}
}
