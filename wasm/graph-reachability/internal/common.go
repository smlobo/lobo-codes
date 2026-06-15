package internal

import "math"

var (
	WindowHeight int64
	WindowWidth  int64
)
var (
	MaxY int64
	MaxX int64
)

const numVertices int = 20
const Delta int64 = 60
const Margin int64 = 100

const pointDiameter float32 = 30.0
const pointOutline float32 = 2.0

func quantized(i int64) uint {
	return uint(i / Delta)
}

func weight(p, q *Vertex) uint64 {
	xDiff := math.Abs(float64(p.x - q.x))
	yDiff := math.Abs(float64(p.y - q.y))
	return uint64(math.Sqrt(xDiff*xDiff + yDiff*yDiff))
}

func euclidean(x, y uint) float64 {
	return math.Sqrt(float64(x*x + y*y))
}
