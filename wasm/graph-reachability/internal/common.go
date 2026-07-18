package internal

import "math"

var (
	WindowHeight int
	WindowWidth  int
)
var (
	MaxY int
	MaxX int
)

const numVertices int = 8
const Delta = int(pointDiameter)
const Margin int = 100

const pointDiameter float32 = 30.0
const pointOutline float32 = 2.0

func quantized(i int) int {
	return i / Delta
}

func weight(p, q *Vertex) uint64 {
	xDiff := math.Abs(float64(p.x - q.x))
	yDiff := math.Abs(float64(p.y - q.y))
	return uint64(math.Sqrt(xDiff*xDiff + yDiff*yDiff))
}

func euclidean(x, y int) float64 {
	return math.Sqrt(float64(x*x + y*y))
}
