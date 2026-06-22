package internal

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Vertex struct {
	id       uint
	x, y     int64
	qPoint   QuantizedPoint
	color    color.Color
	outgoing []*Edge
	incoming []*Edge
}

func (v *Vertex) String() string {
	s := fmt.Sprintf("{%d} (%d, %d) / %v\n", v.id, v.x, v.y, v.qPoint)
	s += fmt.Sprintf("  Outgoing: %v\n", v.outgoing)
	s += fmt.Sprintf("  Incoming: %v\n", v.incoming)
	return s
}

func (v *Vertex) dfs(visitedSet *map[*Vertex]struct{}) {
	// Terminating condition; if in set do not recurse
	if _, ok := (*visitedSet)[v]; ok {
		return
	}

	// Add to set
	(*visitedSet)[v] = struct{}{}

	for _, edge := range v.outgoing {
		edge.to.dfs(visitedSet)
	}
}

func (v *Vertex) draw(c *fyne.Container) {
	// Create a filled circle
	circle := canvas.Circle{
		Position1:   fyne.NewPos(float32(v.x)-pointDiameter/2, float32(v.y)-pointDiameter/2),
		Position2:   fyne.NewPos(float32(v.x)+pointDiameter/2, float32(v.y)+pointDiameter/2),
		FillColor:   v.color,
		StrokeColor: color.Color(color.RGBA{R: 6, G: 57, B: 112, A: 0xFF}),
		StrokeWidth: pointOutline,
	}
	c.Add(&circle)

	label := canvas.NewText(fmt.Sprintf("%d", v.id), color.White)
	label.TextSize = 14
	label.Alignment = fyne.TextAlignCenter

	size := label.MinSize()
	label.Move(fyne.NewPos(
		float32(v.x)-size.Width/2,
		float32(v.y)-size.Height/2,
	))
	label.Resize(size)

	c.Add(label)
}
