package internal

import (
	"fmt"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type coordinates3D struct {
	x, y, z float32
}

type line3D struct {
	a, b *coordinates3D
}

type rotationAngle3D struct {
	x, y, z float32
}

var cubePoints [8]coordinates3D
var cubeLines [12]line3D
var cubeFaces [6][4]*coordinates3D
var centroid coordinates3D

var Rotation rotationAngle3D
var DefaultRotationAngle rotationAngle3D
var ZeroRotationAngle rotationAngle3D

const transparency = 100
const delta = 0.00001

func generateLinePoints(p, q fyne.Position) []fyne.Position {
	linePositions := []fyne.Position{}

	dX := p.X - q.X
	dY := p.Y - q.Y
	if math.Abs(float64(dX)) > math.Abs(float64(dY)) {
		origin := p
		end := q
		if p.X > q.X {
			origin = q
			end = p
		}
		sideDX := end.X - origin.X
		sideDY := end.Y - origin.Y
		ratio := sideDY / sideDX
		for i := 0; i < int(sideDX); i++ {
			linePositions = append(linePositions, fyne.Position{
				X: origin.X + float32(i),
				Y: origin.Y + ratio*float32(i),
			})
		}
	} else {
		origin := p
		end := q
		if p.Y > q.Y {
			origin = q
			end = p
		}
		sideDX := end.X - origin.X
		sideDY := end.Y - origin.Y
		ratio := sideDX / sideDY
		for i := 0; i < int(sideDY); i++ {
			linePositions = append(linePositions, fyne.Position{
				X: origin.X + ratio*float32(i),
				Y: origin.Y + float32(i),
			})
		}
	}
	return linePositions
}

func colorFace(fContainer *fyne.Container, corners [4]fyne.Position, color color.Color) {
	// Generate all points on opposite lines
	side1Positions := generateLinePoints(corners[0], corners[1])
	side2Positions := generateLinePoints(corners[2], corners[3])
	if len(side1Positions) != len(side2Positions) {
		fmt.Printf("Rhombus sides unequal: %d != %d\n", len(side1Positions), len(side2Positions))
	}
	numSidePositions := len(side1Positions)
	if len(side2Positions) < len(side1Positions) {
		numSidePositions = len(side2Positions)
	}
	//fmt.Printf("line positions: %d\n", len(side1Positions))

	// Draw lines to color
	for i := 0; i < numSidePositions; i += 5 {
		fLine := canvas.Line{
			Position1:   side1Positions[i],
			Position2:   side2Positions[i],
			StrokeColor: color,
			StrokeWidth: 2,
		}
		fContainer.Add(&fLine)
	}
}

func getFaceFynePositions(vertices [4]*coordinates3D) [4]fyne.Position {
	return [4]fyne.Position{
		{vertices[0].x, vertices[0].y},
		{vertices[1].x, vertices[1].y},
		{vertices[2].x, vertices[2].y},
		{vertices[3].x, vertices[3].y},
	}
}

func DrawCube(fContainer *fyne.Container) {
	// 0 Rotation angle
	if Rotation.x < delta && Rotation.y < delta && Rotation.z < delta {
		return
	}

	// Clear the previous drawing
	fContainer.RemoveAll()

	//for _, cLine := range cubeLines {
	//	fLine := canvas.Line{
	//		Position1:   fyne.Position{cLine.a.x, cLine.a.y},
	//		Position2:   fyne.Position{cLine.b.x, cLine.b.y},
	//		StrokeColor: colornames.Red,
	//		StrokeWidth: 2,
	//	}
	//	fContainer.Add(&fLine)
	//}

	// Color the faces
	colorFace(fContainer, getFaceFynePositions(cubeFaces[0]), color.NRGBA{R: 0, G: 255, B: 0, A: transparency})
	colorFace(fContainer, getFaceFynePositions(cubeFaces[1]), color.NRGBA{R: 255, G: 255, B: 0, A: transparency})
	colorFace(fContainer, getFaceFynePositions(cubeFaces[2]), color.NRGBA{R: 255, G: 0, B: 0, A: transparency})
	colorFace(fContainer, getFaceFynePositions(cubeFaces[3]), color.NRGBA{R: 0, G: 255, B: 255, A: transparency})
	colorFace(fContainer, getFaceFynePositions(cubeFaces[4]), color.NRGBA{R: 0, G: 0, B: 255, A: transparency})
	colorFace(fContainer, getFaceFynePositions(cubeFaces[5]), color.NRGBA{R: 255, G: 0, B: 255, A: transparency})

	// Show the current drawing
	fContainer.Show()
	fyne.Do(func() {
		fContainer.Refresh()
	})
}

func rotatePoint(p *coordinates3D, r rotationAngle3D) {
	fixed := coordinates3D{p.x, p.y, p.z}
	p.y = float32(math.Cos(float64(r.x)))*fixed.y - float32(math.Sin(float64(r.x)))*fixed.z
	p.z = float32(math.Sin(float64(r.x)))*fixed.y + float32(math.Cos(float64(r.x)))*fixed.z

	fixed = coordinates3D{p.x, p.y, p.z}
	p.x = float32(math.Cos(float64(r.y)))*fixed.x + float32(math.Sin(float64(r.y)))*fixed.z
	p.z = -float32(math.Sin(float64(r.y)))*fixed.x + float32(math.Cos(float64(r.y)))*fixed.z

	fixed = coordinates3D{p.x, p.y, p.z}
	p.x = float32(math.Cos(float64(r.z)))*fixed.x - float32(math.Sin(float64(r.z)))*fixed.y
	p.y = float32(math.Sin(float64(r.z)))*fixed.x + float32(math.Cos(float64(r.z)))*fixed.y
}

func Rotate() {
	// 0 Rotation angle
	if Rotation.x < delta && Rotation.y < delta && Rotation.z < delta {
		return
	}

	// Rotate each point of the cube
	for i := 0; i < len(cubePoints); i++ {
		// Translate to the origin
		p := coordinates3D{
			x: cubePoints[i].x - centroid.x,
			y: cubePoints[i].y - centroid.y,
			z: cubePoints[i].z - centroid.z,
		}
		rotatePoint(&p, Rotation)
		cubePoints[i].x = p.x + centroid.x
		cubePoints[i].y = p.y + centroid.y
		cubePoints[i].z = p.z + centroid.z
	}
}
