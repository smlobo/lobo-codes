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
	x, y         int64
	qPoint       QuantizedPoint
	color        color.Color
	outlineColor color.Color
	outline      float32
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

func (v *Vertex) doColor() bool {
	if len(v.cycleSlice) < 2 {
		return false
	}
	doColor := true
	if v.timer != nil {
		select {
		case <-v.timer.C:
			doColor = true
		default:
			doColor = false
		}
	}
	if !doColor {
		return false
	}
	v.timer = time.NewTimer(time.Millisecond)
	return doColor
}

func (v *Vertex) colorCycle() {
	if len(v.cycleSlice) == 0 {
		return
	}
	if v.nextCycle == len(v.cycleSlice) {
		v.nextCycle = 0
	}
	fmt.Printf("Coloring [%d] cycle: %v\n", v.nextCycle, v.cycleSlice[v.nextCycle])
	for _, cEdge := range v.cycleSlice[v.nextCycle] {
		cEdge.color = offsetColor(darkOrange, v.nextCycle)
		cEdge.stroke = highlightArrowStroke
		cEdge.from.color = offsetColor(orange, v.nextCycle)
		cEdge.from.outlineColor = offsetColor(darkOrange, v.nextCycle)
		cEdge.from.outline = highlightPointOutline
	}
	v.nextCycle++

	//for _, cycle := range v.cycleSlice {
	//	for _, cEdge := range cycle {
	//		cEdge.color = offsetColor(gray, v.nextCycle)
	//		cEdge.from.color = offsetColor(orange, v.nextCycle)
	//	}
	//}
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
