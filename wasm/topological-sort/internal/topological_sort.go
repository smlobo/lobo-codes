package internal

import (
	"fmt"

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
