package internal

import (
	"fmt"
	"math/rand/v2"
	"sort"

	"fyne.io/fyne/v2"
)

type DirectedGraph struct {
	vertices map[QuantizedPoint]*Vertex
	edges    map[EdgeId]*Edge
	nextID   uint
}

func NewDirectedGraph() *DirectedGraph {
	g := DirectedGraph{
		vertices: map[QuantizedPoint]*Vertex{},
		edges:    map[EdgeId]*Edge{},
		nextID:   0,
	}

	MaxX = WindowWidth / Delta
	MaxY = WindowHeight / Delta

	// Generate n random vertices
	for len(g.vertices) < numVertices {
		x := rand.IntN(WindowWidth-Margin) + Margin/2
		y := rand.IntN(WindowHeight-Margin) + Margin/2
		point := QuantizedPoint{quantized(x), quantized(y)}
		if existingQPoint, ok := g.quantizedPointOrNeighbor(point); !ok {
			g.vertices[point] = &Vertex{
				id:       g.nextID,
				x:        x,
				y:        y,
				qPoint:   point,
				color:    green,
				outgoing: []*Edge{},
				incoming: []*Edge{},
			}
			fmt.Printf("[%d] New %v\n", len(g.vertices), g.vertices[point])
			g.nextID++
		} else {
			fmt.Printf("[%d] Existing %v : %v\n", len(g.vertices), point, existingQPoint)
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
				from:  fromV,
				to:    toV,
				color: black,
			})
		}
	}
	sort.Slice(completeEdges, func(i, j int) bool {
		return completeEdges[i].weight() < completeEdges[j].weight()
	})
	fmt.Printf("Complete Edges: %d\n", len(completeEdges))
	for i, edge := range completeEdges {
		fmt.Printf("  [%d] %v\n", i, &edge)
		intersects := false
		for _, existingEdge := range g.edges {
			if existingEdge.intersects(&edge) {
				intersects = true
				fmt.Printf("    Intersecting Edge %v\n", existingEdge)
				break
			}
		}
		if !intersects {
			eId := EdgeId{
				from: edge.from.id,
				to:   edge.to.id,
			}
			g.edges[eId] = &edge
			edge.from.outgoing = append(edge.from.outgoing, &edge)
			edge.to.incoming = append(edge.to.incoming, &edge)
			fmt.Printf("    Created Edge %v\n", &edge)
		}
	}

	return &g
}

func (g *DirectedGraph) quantizedPointOrNeighbor(p QuantizedPoint) (QuantizedPoint, bool) {
	if _, ok := g.vertices[p]; ok {
		return p, true
	}
	neighbors := []QuantizedPoint{
		{p.x - 1, p.y - 1},
		{p.x - 1, p.y},
		{p.x - 1, p.y + 1},
		{p.x, p.y - 1},
		{p.x, p.y + 1},
		{p.x + 1, p.y - 1},
		{p.x + 1, p.y},
		{p.x + 1, p.y + 1},
	}
	for _, neighbor := range neighbors {
		if _, ok := g.vertices[neighbor]; ok {
			return neighbor, true
		}
	}
	return QuantizedPoint{}, false
}

func (g *DirectedGraph) Draw(c *fyne.Container) {
	for _, edge := range g.edges {
		edge.draw(c)
	}
	for _, vertex := range g.vertices {
		vertex.draw(c)
	}
}

func (g *DirectedGraph) HighlightReachableAt(point QuantizedPoint) {
	for _, v := range g.vertices {
		v.color = green
	}
	for _, e := range g.edges {
		e.color = black
	}
	if vertex, ok := g.vertices[point]; ok {
		reachableSlice := g.Reachable(vertex)
		for v, _ := range reachableSlice {
			v.color = red
			for _, e := range v.outgoing {
				e.color = darkRed
			}
		}
	}
}

func (g *DirectedGraph) edgeExists(p, q *Vertex) bool {
	either := EdgeId{
		from: p.id,
		to:   q.id,
	}
	other := EdgeId{
		from: q.id,
		to:   p.id,
	}
	if _, ok := g.edges[either]; ok {
		return true
	}
	if _, ok := g.edges[other]; ok {
		return true
	}
	return false
}

func (g *DirectedGraph) createEdgesForNewVertex(vertex *Vertex) {
	// Create a complete graph, randomly choosing from/to edges
	completeEdges := []Edge{}
	for _, v := range g.vertices {
		if v == vertex {
			continue
		}
		var fromV, toV *Vertex
		if rand.Int()%2 == 0 {
			fromV = vertex
			toV = v
		} else {
			fromV = v
			toV = vertex
		}
		completeEdges = append(completeEdges, Edge{
			from:  fromV,
			to:    toV,
			color: black,
		})
	}
	sort.Slice(completeEdges, func(i, j int) bool {
		return completeEdges[i].weight() < completeEdges[j].weight()
	})
	fmt.Printf("[Insert] Complete Edges: %d\n", len(completeEdges))
	newEdges := []*Edge{}
	for i, edge := range completeEdges {
		fmt.Printf("  [%d] %v\n", i, &edge)
		intersects := false
		for _, existingEdge := range g.edges {
			if existingEdge.intersects(&edge) {
				intersects = true
				fmt.Printf("    [Insert] Intersecting Edge %v\n", existingEdge)
				break
			}
		}
		if !intersects {
			eId := EdgeId{
				from: edge.from.id,
				to:   edge.to.id,
			}
			g.edges[eId] = &edge
			newEdges = append(newEdges, &edge)
			edge.from.outgoing = append(edge.from.outgoing, &edge)
			edge.to.incoming = append(edge.to.incoming, &edge)
			fmt.Printf("    [Insert] Created Edge %v\n", &edge)
		}
	}
}

func (g *DirectedGraph) createEdgesForVertices(vertices []*Vertex) {
	// Create a complete graph, randomly choosing from/to edges
	completeEdges := []Edge{}
	for i := 0; i < len(vertices); i++ {
		for j := i + 1; j < len(vertices); j++ {
			if g.edgeExists(vertices[i], vertices[j]) {
				continue
			}
			var fromV, toV *Vertex
			if rand.Int()%2 == 0 {
				fromV = vertices[i]
				toV = vertices[j]
			} else {
				fromV = vertices[j]
				toV = vertices[i]
			}
			completeEdges = append(completeEdges, Edge{
				from:  fromV,
				to:    toV,
				color: black,
			})
		}
	}
	sort.Slice(completeEdges, func(i, j int) bool {
		return completeEdges[i].weight() < completeEdges[j].weight()
	})
	fmt.Printf("[Remove] Complete Edges: %d\n", len(completeEdges))
	newEdges := []*Edge{}
	for i, edge := range completeEdges {
		fmt.Printf("  [%d] %v\n", i, &edge)
		intersects := false
		for _, existingEdge := range g.edges {
			if existingEdge.intersects(&edge) {
				intersects = true
				fmt.Printf("    [Remove] Intersecting Edge %v\n", existingEdge)
				break
			}
		}
		if !intersects {
			eId := EdgeId{
				from: edge.from.id,
				to:   edge.to.id,
			}
			g.edges[eId] = &edge
			newEdges = append(newEdges, &edge)
			edge.from.outgoing = append(edge.from.outgoing, &edge)
			edge.to.incoming = append(edge.to.incoming, &edge)
			fmt.Printf("    [Remove] Created Edge %v\n", &edge)
		}
	}
}

func (g *DirectedGraph) Update(point QuantizedPoint, x, y int) {
	// Remove vertex?
	if qPoint, ok := g.quantizedPointOrNeighbor(point); ok {
		vertex, _ := g.vertices[qPoint]
		fmt.Printf("Removing vertex %v\n", vertex)
		delete(g.vertices, qPoint)
		surrounding := []*Vertex{}
		for _, out := range vertex.outgoing {
			delete(g.edges, EdgeId{from: out.from.id, to: out.to.id})
			surrounding = append(surrounding, out.to)
			out.to.deleteIncoming(out)
		}
		for _, in := range vertex.incoming {
			delete(g.edges, EdgeId{from: in.from.id, to: in.to.id})
			surrounding = append(surrounding, in.from)
			in.from.deleteOutgoing(in)
		}
		// Create new random non-intersecting edges
		g.createEdgesForVertices(surrounding)
		return
	}

	// Insert vertex
	g.vertices[point] = &Vertex{
		id:       g.nextID,
		x:        x,
		y:        y,
		qPoint:   point,
		color:    green,
		outgoing: []*Edge{},
		incoming: []*Edge{},
	}
	fmt.Printf("[%d] [Insert] New %v\n", len(g.vertices), g.vertices[point])
	g.nextID++
	g.createEdgesForNewVertex(g.vertices[point])
}

func (g *DirectedGraph) Reachable(vertex *Vertex) map[*Vertex]struct{} {
	visitedSet := map[*Vertex]struct{}{}
	vertex.dfs(&visitedSet)
	return visitedSet
}

func (g *DirectedGraph) Redraw(c *fyne.Container) {
	c.Objects = nil
	g.Draw(c)
	c.Refresh()
}
