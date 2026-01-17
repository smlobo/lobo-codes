package internal

import (
	"fmt"
	"math"
)

const margin = 10

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
	centerX := windowWidth / 2.0
	centerY := windowHeight / 2.0
	origin := coordinates3D{centerX, centerY, 0}
	cubePoints = [8]coordinates3D{
		{origin.x - halfSide, origin.y - halfSide, origin.z - halfSide},
		{origin.x + halfSide, origin.y - halfSide, origin.z - halfSide},
		{origin.x - halfSide, origin.y + halfSide, origin.z - halfSide},
		{origin.x + halfSide, origin.y + halfSide, origin.z - halfSide},
		{origin.x - halfSide, origin.y - halfSide, origin.z + halfSide},
		{origin.x + halfSide, origin.y - halfSide, origin.z + halfSide},
		{origin.x - halfSide, origin.y + halfSide, origin.z + halfSide},
		{origin.x + halfSide, origin.y + halfSide, origin.z + halfSide},
	}
	cubeLines = [12]line3D{
		{&cubePoints[0], &cubePoints[1]},
		{&cubePoints[0], &cubePoints[2]},
		{&cubePoints[1], &cubePoints[3]},
		{&cubePoints[2], &cubePoints[3]},
		{&cubePoints[0], &cubePoints[4]},
		{&cubePoints[1], &cubePoints[5]},
		{&cubePoints[2], &cubePoints[6]},
		{&cubePoints[3], &cubePoints[7]},
		{&cubePoints[4], &cubePoints[5]},
		{&cubePoints[4], &cubePoints[6]},
		{&cubePoints[5], &cubePoints[7]},
		{&cubePoints[6], &cubePoints[7]},
	}
	cubeFaces = [6][4]*coordinates3D{
		{&cubePoints[0], &cubePoints[1], &cubePoints[2], &cubePoints[3]},
		{&cubePoints[0], &cubePoints[1], &cubePoints[4], &cubePoints[5]},
		{&cubePoints[0], &cubePoints[2], &cubePoints[4], &cubePoints[6]},
		{&cubePoints[1], &cubePoints[3], &cubePoints[5], &cubePoints[7]},
		{&cubePoints[2], &cubePoints[3], &cubePoints[6], &cubePoints[7]},
		{&cubePoints[4], &cubePoints[5], &cubePoints[6], &cubePoints[7]},
	}
	for _, cubePoint := range cubePoints {
		centroid.x += cubePoint.x
		centroid.y += cubePoint.y
		centroid.z += cubePoint.z
	}
	centroid.x /= float32(len(cubePoints))
	centroid.y /= float32(len(cubePoints))
	centroid.z /= float32(len(cubePoints))

	DefaultRotationAngle = rotationAngle3D{
		x: 0.02,
		y: 0.01,
		z: 0.03,
	}
	ZeroRotationAngle = rotationAngle3D{
		x: 0.0,
		y: 0.0,
		z: 0.0,
	}
	Rotation = DefaultRotationAngle
}
