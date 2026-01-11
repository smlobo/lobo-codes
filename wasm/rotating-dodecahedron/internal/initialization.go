package internal

import (
	"fmt"
	"math"
)

const margin = 10
const goldenRatio = 1.618

var windowWidth, windowHeight, dimension float32
var Paused bool

func Initialize(width, height float32) {
	windowWidth = width
	windowHeight = height
	dimension = min(windowWidth, windowHeight) - margin

	// Cube side is constrained by the max diagonal
	side := dimension / float32(math.Sqrt(3.0))
	fmt.Printf("windowWidth: %.2f, windowHeight: %.2f, dimension: %.2f, side: %.2f\n",
		windowWidth, windowHeight, dimension, side)

	// Cube
	halfSide := side / 2
	phi := halfSide * goldenRatio
	rPhi := halfSide / goldenRatio
	centerX := windowWidth / 2.0
	centerY := windowHeight / 2.0
	origin := coordinates3D{centerX, centerY, 0}
	dodecahedronPoints = [20]coordinates3D{
		{origin.x - halfSide, origin.y - halfSide, origin.z - halfSide},
		{origin.x - halfSide, origin.y + halfSide, origin.z - halfSide},
		{origin.x + halfSide, origin.y - halfSide, origin.z - halfSide},
		{origin.x + halfSide, origin.y + halfSide, origin.z - halfSide},
		{origin.x - halfSide, origin.y - halfSide, origin.z + halfSide},
		{origin.x - halfSide, origin.y + halfSide, origin.z + halfSide},
		{origin.x + halfSide, origin.y - halfSide, origin.z + halfSide},
		{origin.x + halfSide, origin.y + halfSide, origin.z + halfSide},
		{origin.x, origin.y - phi, origin.z - rPhi},
		{origin.x, origin.y + phi, origin.z - rPhi},
		{origin.x, origin.y - phi, origin.z + rPhi},
		{origin.x, origin.y + phi, origin.z + rPhi},
		{origin.x - rPhi, origin.y, origin.z - phi},
		{origin.x + rPhi, origin.y, origin.z - phi},
		{origin.x - rPhi, origin.y, origin.z + phi},
		{origin.x + rPhi, origin.y, origin.z + phi},
		{origin.x - phi, origin.y - rPhi, origin.z},
		{origin.x + phi, origin.y - rPhi, origin.z},
		{origin.x - phi, origin.y + rPhi, origin.z},
		{origin.x + phi, origin.y + rPhi, origin.z},
	}
	dodecahedronLines = [30]line3D{
		{&dodecahedronPoints[8], &dodecahedronPoints[0]},
		{&dodecahedronPoints[8], &dodecahedronPoints[2]},
		{&dodecahedronPoints[8], &dodecahedronPoints[10]},
		{&dodecahedronPoints[10], &dodecahedronPoints[4]},
		{&dodecahedronPoints[10], &dodecahedronPoints[6]},
		{&dodecahedronPoints[9], &dodecahedronPoints[1]},
		{&dodecahedronPoints[9], &dodecahedronPoints[3]},
		{&dodecahedronPoints[9], &dodecahedronPoints[11]},
		{&dodecahedronPoints[11], &dodecahedronPoints[5]},
		{&dodecahedronPoints[11], &dodecahedronPoints[7]},
		{&dodecahedronPoints[12], &dodecahedronPoints[0]},
		{&dodecahedronPoints[12], &dodecahedronPoints[1]},
		{&dodecahedronPoints[12], &dodecahedronPoints[13]},
		{&dodecahedronPoints[13], &dodecahedronPoints[2]},
		{&dodecahedronPoints[13], &dodecahedronPoints[3]},
		{&dodecahedronPoints[14], &dodecahedronPoints[4]},
		{&dodecahedronPoints[14], &dodecahedronPoints[5]},
		{&dodecahedronPoints[14], &dodecahedronPoints[15]},
		{&dodecahedronPoints[15], &dodecahedronPoints[6]},
		{&dodecahedronPoints[15], &dodecahedronPoints[7]},
		{&dodecahedronPoints[16], &dodecahedronPoints[0]},
		{&dodecahedronPoints[16], &dodecahedronPoints[4]},
		{&dodecahedronPoints[16], &dodecahedronPoints[18]},
		{&dodecahedronPoints[18], &dodecahedronPoints[1]},
		{&dodecahedronPoints[18], &dodecahedronPoints[5]},
		{&dodecahedronPoints[17], &dodecahedronPoints[2]},
		{&dodecahedronPoints[17], &dodecahedronPoints[6]},
		{&dodecahedronPoints[17], &dodecahedronPoints[19]},
		{&dodecahedronPoints[19], &dodecahedronPoints[3]},
		{&dodecahedronPoints[19], &dodecahedronPoints[7]},
	}
	for _, dodecahedronPoint := range dodecahedronPoints {
		centroid.x += dodecahedronPoint.x
		centroid.y += dodecahedronPoint.y
		centroid.z += dodecahedronPoint.z
	}
	centroid.x /= float32(len(dodecahedronPoints))
	centroid.y /= float32(len(dodecahedronPoints))
	centroid.z /= float32(len(dodecahedronPoints))
}
