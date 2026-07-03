package internal

import (
	"fmt"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type TopologicalSort struct {
	graph           *DirectedAcyclicGraph
	orderedVertices []*Vertex
}

func NewTopologicalSort(g *DirectedAcyclicGraph) *TopologicalSort {
	t := TopologicalSort{
		graph:           g,
		orderedVertices: []*Vertex{},
	}

	g.resetVisited()

	postOrder := []*Vertex{}
	for _, v := range g.vertices {
		t.dfs(v, &postOrder)
	}
	fmt.Printf("post order: %d\n", len(postOrder))

	// Reverse post order
	fmt.Printf("Topological order: ")
	for i := len(postOrder) - 1; i >= 0; i-- {
		t.orderedVertices = append(t.orderedVertices, postOrder[i])
		fmt.Printf("%d, ", postOrder[i].id)
	}
	fmt.Printf("\n")

	return &t
}

func (t *TopologicalSort) dfs(v *Vertex, postOrder *[]*Vertex) {
	if v.visited {
		return
	}
	v.visited = true
	for _, e := range v.outgoing {
		t.dfs(e.to, postOrder)
	}
	*postOrder = append(*postOrder, v)
}

func (t *TopologicalSort) draw(c *fyne.Container) {
	axisLength := 0
	fixedCoord := 0
	if TopologicalOrientation == Vertical {
		axisLength = GraphHeight - Margin
		fixedCoord = GraphWidth + TopologicalReserved/2
	} else {
		axisLength = GraphWidth - Margin
		fixedCoord = GraphHeight + TopologicalReserved/2
	}
	segmentLength := axisLength / (len(t.orderedVertices) + 1)
	segmentSlice := []int{}
	fixedSlice := []int{}
	for i := 1; i <= len(t.orderedVertices); i++ {
		segmentSlice = append(segmentSlice, i*segmentLength)
		fixedSlice = append(fixedSlice, fixedCoord)
	}
	var xCoords, yCoords *[]int
	if TopologicalOrientation == Vertical {
		xCoords = &fixedSlice
		yCoords = &segmentSlice
	} else {
		xCoords = &segmentSlice
		yCoords = &fixedSlice
	}

	sortedPositions := map[*Vertex]fyne.Position{}
	sortedIndexes := map[*Vertex]int{}
	for i, v := range t.orderedVertices {
		sortedPositions[v] = fyne.NewPos(float32((*xCoords)[i]), float32((*yCoords)[i]))
		sortedIndexes[v] = i
	}

	// Draw the topological edges first so vertices and labels remain readable.
	curvedEdgeCount := map[*Vertex]int{}
	for edge := range t.graph.edges {
		from, okFrom := sortedPositions[edge.from]
		to, okTo := sortedPositions[edge.to]
		if !okFrom || !okTo {
			continue
		}

		side := 1
		if absInt(sortedIndexes[edge.to]-sortedIndexes[edge.from]) != 1 {
			if curvedEdgeCount[edge.from]%2 == 1 {
				side = -1
			}
			curvedEdgeCount[edge.from]++
		}

		drawTopologicalEdge(c, edge, from, to, sortedIndexes[edge.from], sortedIndexes[edge.to], side)
	}

	// Draw Topological Sorted Vertices
	for i := 0; i < len(t.orderedVertices); i++ {
		x := (*xCoords)[i]
		y := (*yCoords)[i]
		v := t.orderedVertices[i]
		circle := canvas.Circle{
			Position1:   fyne.NewPos(float32(x)-pointDiameter/2, float32(y)-pointDiameter/2),
			Position2:   fyne.NewPos(float32(x)+pointDiameter/2, float32(y)+pointDiameter/2),
			FillColor:   v.color,
			StrokeColor: v.outlineColor,
			StrokeWidth: v.outline,
		}
		c.Add(&circle)

		label := canvas.NewText(fmt.Sprintf("%d", v.id), v.textColor)
		label.TextSize = 14
		label.Alignment = fyne.TextAlignCenter

		size := label.MinSize()
		label.Move(fyne.NewPos(
			float32(x)-size.Width/2,
			float32(y)-size.Height/2,
		))
		label.Resize(size)

		c.Add(label)
	}
}

func drawTopologicalEdge(c *fyne.Container, edge *Edge, from fyne.Position, to fyne.Position, fromIndex int, toIndex int, side int) {
	const (
		arrowAngle       = math.Pi / 6
		arrowLength      = 12
		curveSegments    = 20
		minControlOffset = float64(pointDiameter)
		distanceBendRate = 0.35
		vertexClearance  = float64(pointDiameter) / 2
	)

	dx := float64(to.X - from.X)
	dy := float64(to.Y - from.Y)
	distance := math.Hypot(dx, dy)
	if distance == 0 {
		return
	}

	if absInt(toIndex-fromIndex) == 1 {
		drawTopologicalStraightEdge(c, edge, from, to)
		return
	}

	maxControlOffset := float64(TopologicalReserved) - float64(pointDiameter)
	offset := math.Max(minControlOffset, distance*distanceBendRate)
	offset = math.Min(offset, maxControlOffset)
	midX := float64(from.X+to.X) / 2
	midY := float64(from.Y+to.Y) / 2

	controlX := midX
	controlY := midY
	if TopologicalOrientation == Vertical {
		controlX += float64(side) * offset
	} else {
		controlY += float64(side) * offset
	}

	points := make([]fyne.Position, 0, curveSegments+1)
	for i := 0; i <= curveSegments; i++ {
		t := float64(i) / curveSegments
		x := quadraticBezier(float64(from.X), controlX, float64(to.X), t)
		y := quadraticBezier(float64(from.Y), controlY, float64(to.Y), t)

		if math.Hypot(x-float64(from.X), y-float64(from.Y)) < vertexClearance ||
			math.Hypot(x-float64(to.X), y-float64(to.Y)) < vertexClearance {
			continue
		}
		points = append(points, fyne.NewPos(float32(x), float32(y)))
	}

	if len(points) < 2 {
		return
	}

	for i := 0; i < len(points)-1; i++ {
		line := canvas.NewLine(edge.color)
		line.StrokeWidth = edge.stroke
		line.Position1 = points[i]
		line.Position2 = points[i+1]
		c.Add(line)
	}

	end := points[len(points)-1]
	beforeEnd := points[len(points)-2]
	theta := math.Atan2(float64(end.Y-beforeEnd.Y), float64(end.X-beforeEnd.X))

	leftHead := canvas.NewLine(edge.color)
	leftHead.StrokeWidth = edge.stroke
	leftHead.Position1 = end
	leftHead.Position2 = fyne.NewPos(
		end.X-float32(arrowLength*math.Cos(theta-arrowAngle)),
		end.Y-float32(arrowLength*math.Sin(theta-arrowAngle)),
	)
	c.Add(leftHead)

	rightHead := canvas.NewLine(edge.color)
	rightHead.StrokeWidth = edge.stroke
	rightHead.Position1 = end
	rightHead.Position2 = fyne.NewPos(
		end.X-float32(arrowLength*math.Cos(theta+arrowAngle)),
		end.Y-float32(arrowLength*math.Sin(theta+arrowAngle)),
	)
	c.Add(rightHead)
}

func quadraticBezier(start, control, end, t float64) float64 {
	return (1-t)*(1-t)*start + 2*(1-t)*t*control + t*t*end
}

func drawTopologicalStraightEdge(c *fyne.Container, edge *Edge, from fyne.Position, to fyne.Position) {
	const (
		arrowAngle  = math.Pi / 6
		arrowLength = 12
	)

	dx := float64(to.X - from.X)
	dy := float64(to.Y - from.Y)
	distance := math.Hypot(dx, dy)
	if distance == 0 {
		return
	}

	vertexRadius := float64(pointDiameter) / 2
	ux := dx / distance
	uy := dy / distance

	start := fyne.NewPos(
		from.X+float32(ux*vertexRadius),
		from.Y+float32(uy*vertexRadius),
	)
	end := fyne.NewPos(
		to.X-float32(ux*vertexRadius),
		to.Y-float32(uy*vertexRadius),
	)

	line := canvas.NewLine(edge.color)
	line.StrokeWidth = edge.stroke
	line.Position1 = start
	line.Position2 = end
	c.Add(line)

	theta := math.Atan2(dy, dx)

	leftHead := canvas.NewLine(edge.color)
	leftHead.StrokeWidth = edge.stroke
	leftHead.Position1 = end
	leftHead.Position2 = fyne.NewPos(
		end.X-float32(arrowLength*math.Cos(theta-arrowAngle)),
		end.Y-float32(arrowLength*math.Sin(theta-arrowAngle)),
	)
	c.Add(leftHead)

	rightHead := canvas.NewLine(edge.color)
	rightHead.StrokeWidth = edge.stroke
	rightHead.Position1 = end
	rightHead.Position2 = fyne.NewPos(
		end.X-float32(arrowLength*math.Cos(theta+arrowAngle)),
		end.Y-float32(arrowLength*math.Sin(theta+arrowAngle)),
	)
	c.Add(rightHead)
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
