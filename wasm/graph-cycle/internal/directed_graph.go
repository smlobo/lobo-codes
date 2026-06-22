package internal

import (
	"fmt"
	"image/color"
	"math/rand/v2"
	"sort"

	"fyne.io/fyne/v2"
)

type DirectedGraph struct {
	vertices map[QuantizedPoint]*Vertex
	edges    map[*Edge]struct{}
	nextID   uint
}

func NewDirectedGraph() *DirectedGraph {
	g := DirectedGraph{
		vertices: map[QuantizedPoint]*Vertex{},
		edges:    map[*Edge]struct{}{},
		nextID:   0,
	}

	MaxX = WindowWidth / Delta
	MaxY = WindowHeight / Delta

	// Generate n random vertices
	for len(g.vertices) < numVertices {
		x := rand.Int64N(WindowWidth-Margin) + Margin/2
		y := rand.Int64N(WindowHeight-Margin) + Margin/2
		point := QuantizedPoint{quantized(x), quantized(y)}
		if existingVertex, ok := g.vertices[point]; !ok {
			g.vertices[point] = &Vertex{
				id:           g.nextID,
				x:            x,
				y:            y,
				qPoint:       point,
				color:        purple,
				outlineColor: purple,
				outline:      pointOutline,
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
	completeEdges := []Edge{}
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
			completeEdges = append(completeEdges, Edge{
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
		fmt.Printf("  [%d] %v\n", i, &edge)
		intersects := false
		for existingEdge := range g.edges {
			if existingEdge.intersects(&edge) {
				intersects = true
				fmt.Printf("    Intersecting Edge %v\n", &existingEdge)
				break
			}
		}
		if !intersects {
			g.edges[&edge] = struct{}{}
			edge.from.outgoing = append(edge.from.outgoing, &edge)
			edge.from.incoming = append(edge.from.incoming, &edge)
			fmt.Printf("    Created Edge %v\n", &edge)
		}
	}

	return &g
}

func (g *DirectedGraph) Draw(c *fyne.Container) {
	for edge := range g.edges {
		edge.draw(c)
	}
	for _, vertex := range g.vertices {
		vertex.draw(c)
	}
}

func (g *DirectedGraph) resetColors() {
	for _, v := range g.vertices {
		v.color = purple
		v.outlineColor = purple
		v.outline = pointOutline
	}
	for e, _ := range g.edges {
		e.color = color.Black
		e.stroke = arrowStroke
	}
}

func (g *DirectedGraph) HighlightCycleAt(point QuantizedPoint) bool {
	g.resetColors()
	if vertex, ok := g.vertices[point]; ok {
		// Record the cycle slice in the vertex
		g.Cycle(vertex)
		// Color the 1st cycle
		vertex.colorCycle()
		if len(vertex.cycleSlice) > 1 {
			return true
		}
	}
	return false
}

func (g *DirectedGraph) Cycle(vertex *Vertex) {
	// Reset visited before DFS
	for _, v := range g.vertices {
		v.visited = false
	}

	vertex.cycleSlice = [][]*Edge{}
	vertex.nextCycle = 0
	vertex.timer = nil
	vertexStack := []*Edge{}
	for _, edge := range vertex.outgoing {
		edge.dfs(vertex, &vertex.cycleSlice, &vertexStack)
	}
}

func (g *DirectedGraph) HighlightNextCycleAt(point QuantizedPoint) bool {
	if vertex, ok := g.vertices[point]; ok && vertex.doColor() {
		g.resetColors()
		vertex.colorCycle()
		return true
	}
	return false
}

func (g *DirectedGraph) Redraw(c *fyne.Container) {
	c.Objects = nil
	g.Draw(c)
	c.Refresh()
}
