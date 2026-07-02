package internal

import (
	"fmt"
	"image/color"
	"math/rand/v2"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type DirectedAcyclicGraph struct {
	vertices map[QuantizedPoint]*Vertex
	edges    map[*Edge]struct{}
	nextID   uint
}

func NewDirectedAcyclicGraph() *DirectedAcyclicGraph {
	g := DirectedAcyclicGraph{
		vertices: map[QuantizedPoint]*Vertex{},
		edges:    map[*Edge]struct{}{},
		nextID:   0,
	}

	// Generate n random vertices
	for len(g.vertices) < numVertices {
		x := rand.IntN(GraphWidth-Margin) + Margin/2
		y := rand.IntN(GraphHeight-Margin) + Margin/2
		point := QuantizedPoint{quantized(x), quantized(y)}
		if existingVertex, ok := g.vertices[point]; !ok {
			g.vertices[point] = &Vertex{
				id:           g.nextID,
				x:            x,
				y:            y,
				qPoint:       point,
				color:        purple,
				outlineColor: darkPurple,
				outline:      pointOutline,
				textColor:    color.White,
				outgoing:     []*Edge{},
				incoming:     []*Edge{},
				visited:      false,
				cycleSlice:   [][]*Edge{},
				nextCycle:    0,
				timer:        nil,
			}
			fmt.Printf("[%d] New %v\n", len(g.vertices), g.vertices[point])
			g.nextID++
		} else {
			fmt.Printf("[%d] Existing %v : %v\n", len(g.vertices), point, existingVertex)
		}
	}

	// Slice of Vertex keys for ordered iteration
	quantizedPoints := make([]QuantizedPoint, 0, len(g.vertices))
	for p := range g.vertices {
		quantizedPoints = append(quantizedPoints, p)
	}
	sort.Slice(quantizedPoints, func(i, j int) bool {
		return euclidean(quantizedPoints[i].x, quantizedPoints[i].y) <
			euclidean(quantizedPoints[j].x, quantizedPoints[j].y)
	})

	// Create a complete graph, randomly choosing from/to edges
	completeEdges := []*Edge{}
	for i := 0; i < len(quantizedPoints); i++ {
		for j := i + 1; j < len(quantizedPoints); j++ {
			pi := quantizedPoints[i]
			pj := quantizedPoints[j]
			var fromV, toV *Vertex
			if rand.Int()%2 == 0 {
				fromV = g.vertices[pi]
				toV = g.vertices[pj]
			} else {
				fromV = g.vertices[pj]
				toV = g.vertices[pi]
			}
			completeEdges = append(completeEdges, &Edge{
				from:   fromV,
				to:     toV,
				weight: weight(fromV, toV),
				color:  color.Black,
				stroke: arrowStroke,
			})
		}
	}
	sort.Slice(completeEdges, func(i, j int) bool {
		return completeEdges[i].weight < completeEdges[j].weight
	})
	fmt.Printf("Complete Edges: %d\n", len(completeEdges))
	for i, edge := range completeEdges {
		fmt.Printf("  [%d] %v\n", i, edge)
		intersects := false
		for existingEdge := range g.edges {
			if existingEdge.intersects(edge) {
				intersects = true
				fmt.Printf("    Intersecting Edge %v\n", existingEdge)
				break
			}
		}
		if intersects {
			continue
		}

		g.addEdge(edge)
		fmt.Printf("    Potential Edge %v\n", edge)
		if !g.cycle(edge.from) {
			fmt.Printf("    Added Edge %v\n", edge)
			continue
		}

		// Try with the reverse edge
		g.removeEdge(edge)
		reverseEdge := edge.reverse()
		fmt.Printf("    Potential Reverse Edge %v\n", reverseEdge)

		g.addEdge(reverseEdge)
		fmt.Printf("    Potential Edge %v\n", reverseEdge)
		if g.cycle(reverseEdge.from) {
			g.removeEdge(reverseEdge)
		} else {
			fmt.Printf("    Added Edge %v\n", reverseEdge)
		}
	}

	return &g
}

func (g *DirectedAcyclicGraph) addEdge(edge *Edge) {
	g.edges[edge] = struct{}{}
	edge.from.outgoing = append(edge.from.outgoing, edge)
	edge.to.incoming = append(edge.to.incoming, edge)
}

func (g *DirectedAcyclicGraph) removeEdge(edge *Edge) {
	delete(g.edges, edge)
	edge.from.removeOutgoing(edge)
	edge.to.removeIncoming(edge)
}

func (g *DirectedAcyclicGraph) Draw(c *fyne.Container) {
	// Draw a line to divide the space into the original graph and topological sorted
	separator := canvas.NewLine(color.RGBA{R: 255, G: 20, B: 10, A: 255})
	separator.StrokeWidth = 2
	if TopologicalOrientation == Vertical {
		separator.Position1 = fyne.NewPos(float32(GraphWidth), float32(Margin))
		separator.Position2 = fyne.NewPos(float32(GraphWidth), float32(GraphHeight-Margin))
	} else {
		separator.Position1 = fyne.NewPos(float32(Margin), float32(GraphHeight))
		separator.Position2 = fyne.NewPos(float32(GraphWidth-Margin), float32(GraphHeight))
	}
	c.Add(separator)

	for edge := range g.edges {
		edge.draw(c)
	}
	for _, vertex := range g.vertices {
		vertex.draw(c)
	}

	// Compute the topological sort
	tSort := NewTopologicalSort(g)
	tSort.draw(c)
}

func (g *DirectedAcyclicGraph) resetColors() {
	for _, v := range g.vertices {
		v.color = purple
		v.outlineColor = darkPurple
		v.outline = pointOutline
		v.textColor = color.White
	}
	for e, _ := range g.edges {
		e.color = color.Black
		e.stroke = arrowStroke
	}
}

func (g *DirectedAcyclicGraph) highlightVertex(point QuantizedPoint) {
	g.resetColors()
	if vertex, ok := g.vertices[point]; ok {
		vertex.color = blue
		vertex.outlineColor = darkBlue
		vertex.outline = highlightPointOutline
		vertex.textColor = color.Black
		for _, outEdge := range vertex.outgoing {
			outEdge.color = darkGreen
			outEdge.stroke = highlightArrowStroke
		}
		for _, inEdge := range vertex.incoming {
			inEdge.color = darkRed
			inEdge.stroke = highlightArrowStroke
		}
	}
}

func (g *DirectedAcyclicGraph) resetVisited() {
	for _, v := range g.vertices {
		v.visited = false
	}
}

func (g *DirectedAcyclicGraph) cycle(vertex *Vertex) bool {
	for _, edge := range vertex.outgoing {
		toVertex := edge.to
		g.resetVisited()
		toVertex.cycleDFS()
		if vertex.visited {
			return true
		}
	}
	return false
}

func (g *DirectedAcyclicGraph) Redraw(c *fyne.Container) {
	c.Objects = nil
	g.Draw(c)
	c.Refresh()
}
