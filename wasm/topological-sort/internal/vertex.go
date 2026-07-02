package internal

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Vertex struct {
	id           uint
	x, y         int
	qPoint       QuantizedPoint
	color        color.Color
	outlineColor color.Color
	outline      float32
	textColor    color.Color
	outgoing     []*Edge
	incoming     []*Edge
	visited      bool
	cycleSlice   [][]*Edge
	nextCycle    int
	timer        *time.Timer
}

func (v *Vertex) String() string {
	s := fmt.Sprintf("{%d} (%d, %d) / %v\n", v.id, v.x, v.y, v.qPoint)
	s += fmt.Sprintf("  Outgoing: %v\n", v.outgoing)
	s += fmt.Sprintf("  Incoming: %v\n", v.incoming)
	return s
}

func (v *Vertex) removeOutgoing(edge *Edge) {
	// Get index
	eIndex := 0
	for i, e := range v.outgoing {
		if e == edge {
			eIndex = i
			break
		}
	}
	v.outgoing = append(v.outgoing[:eIndex], v.outgoing[eIndex+1:]...)
}

func (v *Vertex) removeIncoming(edge *Edge) {
	// Get index
	eIndex := 0
	for i, e := range v.incoming {
		if e == edge {
			eIndex = i
			break
		}
	}
	v.incoming = append(v.incoming[:eIndex], v.incoming[eIndex+1:]...)
}

func (v *Vertex) cycleDFS() {
	// End condition
	if v.visited {
		return
	}
	v.visited = true
	for _, edge := range v.outgoing {
		toVertex := edge.to
		toVertex.cycleDFS()
	}
}

func (v *Vertex) draw(c *fyne.Container) {
	// Create a filled circle
	circle := canvas.Circle{
		Position1:   fyne.NewPos(float32(v.x)-pointDiameter/2, float32(v.y)-pointDiameter/2),
		Position2:   fyne.NewPos(float32(v.x)+pointDiameter/2, float32(v.y)+pointDiameter/2),
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
		float32(v.x)-size.Width/2,
		float32(v.y)-size.Height/2,
	))
	label.Resize(size)

	c.Add(label)
}
