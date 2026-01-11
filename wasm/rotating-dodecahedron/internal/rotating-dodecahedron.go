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

var dodecahedronPoints [20]coordinates3D
var dodecahedronLines [30]line3D
var centroid coordinates3D

func DrawDodecahedron(fContainer *fyne.Container) {
	// Clear the previous drawing
	fContainer.RemoveAll()

	for _, cLine := range dodecahedronLines {
		fLine := canvas.Line{
			Position1:   fyne.Position{cLine.a.x, cLine.a.y},
			Position2:   fyne.Position{cLine.b.x, cLine.b.y},
			StrokeColor: colornames.Deeppink,
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
	// Rotate each point of the dodecahedron
	for i := 0; i < len(dodecahedronPoints); i++ {
		// Translate to the origin
		p := coordinates3D{
			x: dodecahedronPoints[i].x - centroid.x,
			y: dodecahedronPoints[i].y - centroid.y,
			z: dodecahedronPoints[i].z - centroid.z,
		}
		rotatePoint(&p, 0.02, 0.01, 0.03)
		dodecahedronPoints[i].x = p.x + centroid.x
		dodecahedronPoints[i].y = p.y + centroid.y
		dodecahedronPoints[i].z = p.z + centroid.z
	}
}
