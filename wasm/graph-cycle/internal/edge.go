package internal

import (
	"fmt"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Edge struct {
	from, to *Vertex
	weight   uint64
	color    color.Color
	stroke   float32
}

func (e *Edge) String() string {
	return fmt.Sprintf("(%d -> %d; %d)", e.from.id, e.to.id, e.weight)
}

func (e *Edge) intersects(other *Edge) bool {
	const eps = 1e-9

	if other == nil {
		return false
	}

	// Edges that share a vertex are allowed to meet at that endpoint.
	if e.from == other.from || e.from == other.to || e.to == other.from || e.to == other.to {
		return false
	}

	det := float64(e.from.x-e.to.x)*float64(other.from.y-other.to.y) -
		float64(e.from.y-e.to.y)*float64(other.from.x-other.to.x)
	if math.Abs(det) < eps {
		return false
	}

	t1 := float64(e.from.x*e.to.y - e.from.y*e.to.x)
	t2 := float64(other.from.x*other.to.y - other.from.y*other.to.x)

	intersectX := (t1*float64(other.from.x-other.to.x) - float64(e.from.x-e.to.x)*t2) / det
	intersectY := (t1*float64(other.from.y-other.to.y) - float64(e.from.y-e.to.y)*t2) / det

	onSegment := func(edge *Edge, x float64, y float64) bool {
		return x >= math.Min(float64(edge.from.x), float64(edge.to.x))-eps &&
			x <= math.Max(float64(edge.from.x), float64(edge.to.x))+eps &&
			y >= math.Min(float64(edge.from.y), float64(edge.to.y))-eps &&
			y <= math.Max(float64(edge.from.y), float64(edge.to.y))+eps
	}

	return onSegment(e, intersectX, intersectY) && onSegment(other, intersectX, intersectY)
}

func (e *Edge) dfs(origin *Vertex, cycleSlice *[][]*Edge, vertexStack *[]*Edge) {
	// Terminating condition 1; return to origin -> record cycle
	if e.to == origin {
		cycle := []*Edge{}
		for _, cV := range *vertexStack {
			cycle = append(cycle, cV)
		}
		cycle = append(cycle, e)
		fmt.Printf("[%d] cycle: %v\n", len(*cycleSlice), cycle)
		*cycleSlice = append(*cycleSlice, cycle)
		return
	}
	// Terminating condition 2; other cycle
	if e.to.visited {
		return
	}

	// Mark visited
	e.to.visited = true

	// Push to stack
	*vertexStack = append(*vertexStack, e)

	// Recurse
	for _, edge := range e.to.outgoing {
		edge.dfs(origin, cycleSlice, vertexStack)
	}

	// Unmark
	e.to.visited = false

	// Pop from stack
	*vertexStack = (*vertexStack)[:len(*vertexStack)-1]
}

func (e *Edge) draw(c *fyne.Container) {
	const arrowAngle = math.Pi / 6
	const arrowLength float64 = 12

	fromX := float64(e.from.x)
	fromY := float64(e.from.y)
	toX := float64(e.to.x)
	toY := float64(e.to.y)

	dx := toX - fromX
	dy := toY - fromY
	length := math.Hypot(dx, dy)
	if length == 0 {
		return
	}

	vertexRadius := float64(pointDiameter) / 2
	ux := dx / length
	uy := dy / length

	startX := fromX + ux*vertexRadius
	startY := fromY + uy*vertexRadius
	endX := toX - ux*vertexRadius
	endY := toY - uy*vertexRadius

	shaft := canvas.NewLine(e.color)
	shaft.StrokeWidth = e.stroke
	shaft.Position1 = fyne.NewPos(float32(startX), float32(startY))
	shaft.Position2 = fyne.NewPos(float32(endX), float32(endY))
	c.Add(shaft)

	theta := math.Atan2(dy, dx)
	leftX := endX - arrowLength*math.Cos(theta-arrowAngle)
	leftY := endY - arrowLength*math.Sin(theta-arrowAngle)
	rightX := endX - arrowLength*math.Cos(theta+arrowAngle)
	rightY := endY - arrowLength*math.Sin(theta+arrowAngle)

	leftHead := canvas.NewLine(e.color)
	leftHead.StrokeWidth = e.stroke
	leftHead.Position1 = fyne.NewPos(float32(endX), float32(endY))
	leftHead.Position2 = fyne.NewPos(float32(leftX), float32(leftY))
	c.Add(leftHead)

	rightHead := canvas.NewLine(e.color)
	rightHead.StrokeWidth = e.stroke
	rightHead.Position1 = fyne.NewPos(float32(endX), float32(endY))
	rightHead.Position2 = fyne.NewPos(float32(rightX), float32(rightY))
	c.Add(rightHead)
}
