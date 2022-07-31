package internal

import (
	"fmt"
	"fyne.io/fyne/v2"
	"time"
)

var NumBodies int
var Radius float64
var Bodies []body
var Scale float64

func PrintUniverse() {
	fmt.Printf("# bodies: %d\n", NumBodies)
	fmt.Printf("Radius: %e\n", Radius)
	for index, body := range Bodies {
		fmt.Printf("  [%d] %s\n", index, body)
	}
}

func Simulate(c *fyne.Container, totalTime float64, timeInterval float64) {
	currentTime := 0.0
	for currentTime < totalTime {
		drawBodies(c)

		updateBodiesVelocities(timeInterval)
		updateBodiesPositions(timeInterval)

		// Sleep before drawing the next step
		time.Sleep(10 * time.Millisecond)

		currentTime += timeInterval
	}
}
func drawBodies(c *fyne.Container) {
	for i := 0; i < len(Bodies); i++ {
		Bodies[i].draw(c)
		c.Refresh()
	}
}

func updateBodiesPositions(t float64) {
	for i := 0; i < len(Bodies); i++ {
		(&Bodies[i]).updatePositions(t)
	}
}

func updateBodiesVelocities(t float64) {
	for i := 0; i < len(Bodies); i++ {
		Bodies[i].xForce = 0.0
		Bodies[i].yForce = 0.0
		for j := 0; j < len(Bodies); j++ {
			if i == j {
				continue
			}
			xForce, yForce := (&Bodies[i]).force(&Bodies[j])
			Bodies[i].xForce += xForce
			Bodies[i].yForce += yForce
		}
	}
	for i := 0; i < len(Bodies); i++ {
		xAccel := Bodies[i].xForce / Bodies[i].mass
		yAccel := Bodies[i].yForce / Bodies[i].mass
		Bodies[i].xVelocity += t * xAccel
		Bodies[i].yVelocity += t * yAccel
	}
}
