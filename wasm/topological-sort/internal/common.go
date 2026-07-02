package internal

import "math"

type Orientation int

const (
	Horizontal Orientation = iota
	Vertical
)

var (
	WindowHeight           int
	WindowWidth            int
	GraphWidth             int
	GraphHeight            int
	TopologicalOrientation Orientation
)

const numVertices int = 8
const Delta int = 60
const Margin int = 100
const TopologicalReserved int = Margin * 3

const pointDiameter float32 = 30.0
const pointOutline float32 = 2.0
const highlightPointOutline float32 = 4.0

const arrowStroke float32 = 2.0
const highlightArrowStroke = 3.0

func quantized(i int) uint {
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
