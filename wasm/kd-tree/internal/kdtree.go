package internal

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type KdTree struct {
	root *kdNode
	size int
}

type kdNode struct {
	point *Point
	left  *kdNode
	right *kdNode
}

func (t *KdTree) Populate(set map[int64]Point) {
	for _, p := range set {
		t.Insert(&p)
	}
}

func (t *KdTree) Insert(p *Point) {
	t.root = t.insert(t.root, p, true)
}

func (t *KdTree) insert(n *kdNode, p *Point, xcoord bool) *kdNode {
	if n == nil {
		t.size++
		return &kdNode{p, nil, nil}
	}

	if n.point == p {
		return n
	}

	if (xcoord && p.x < n.point.x) || (!xcoord && p.y < n.point.y) {
		n.left = t.insert(n.left, p, !xcoord)
	} else {
		n.right = t.insert(n.right, p, !xcoord)
	}

	return n
}

func (t *KdTree) Draw(c *fyne.Container) {
	r := Rectangle{
		topLeft:     Point{},
		bottomRight: Point{x: float32(WindowWidth), y: float32(WindowHeight)},
	}
	t.draw(c, t.root, &r, true)
}

func (t *KdTree) draw(c *fyne.Container, n *kdNode, r *Rectangle, xcoord bool) {
	if n == nil {
		return
	}

	if xcoord {
		c.Add(&canvas.Line{
			Position1:   fyne.Position{n.point.x, r.topLeft.y},
			Position2:   fyne.Position{n.point.x, r.bottomRight.y},
			StrokeColor: color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff},
			StrokeWidth: 1,
		})
	} else {
		c.Add(&canvas.Line{
			Position1:   fyne.Position{r.topLeft.x, n.point.y},
			Position2:   fyne.Position{r.bottomRight.x, n.point.y},
			StrokeColor: color.RGBA{R: 0x00, G: 0xff, B: 0x00, A: 0xff},
			StrokeWidth: 1,
		})
	}

	// Draw left tree
	if xcoord {
		r := Rectangle{
			topLeft:     Point{x: r.topLeft.x, y: r.topLeft.y},
			bottomRight: Point{x: n.point.x, y: r.bottomRight.y},
		}
		t.draw(c, n.left, &r, !xcoord)
	} else {
		r := Rectangle{
			topLeft:     Point{x: r.topLeft.x, y: r.topLeft.y},
			bottomRight: Point{x: r.bottomRight.x, y: n.point.y},
		}
		t.draw(c, n.left, &r, !xcoord)
	}

	// Draw the right tree
	if xcoord {
		r := Rectangle{
			topLeft:     Point{x: n.point.x, y: r.topLeft.y},
			bottomRight: Point{x: r.bottomRight.x, y: r.bottomRight.y},
		}
		t.draw(c, n.right, &r, !xcoord)
	} else {
		r := Rectangle{
			topLeft:     Point{x: r.topLeft.x, y: n.point.y},
			bottomRight: Point{x: r.bottomRight.x, y: r.bottomRight.y},
		}
		t.draw(c, n.right, &r, !xcoord)
	}
}

func (t *KdTree) RangePoints(r *Rectangle) []*Point {
	pts := []*Point{}

	t.rangePoints(t.root, pts, r,
		&Rectangle{Point{0.0, 0.0}, Point{float32(WindowWidth), float32(WindowHeight)}}, true)

	return pts
}

func (t *KdTree) rangePoints(n *kdNode, rangePoints []*Point, targetRect *Rectangle, currentRect *Rectangle, xcoord bool) {
	if n == nil {
		return
	}

	if !targetRect.intersects(currentRect) {
		return
	}

	// TODO
}

type distancePair struct {
	point    *Point
	distance float32
}

func (t *KdTree) Nearest(p *Point) *Point {
	if p == nil {
		return nil
	}

	currentNearest := distancePair{
		t.root.point,
		t.root.point.distanceSquaredTo(p),
	}
	currentRectangle := Rectangle{
		Point{0.0, 0.0},
		Point{float32(WindowWidth), float32(WindowHeight)},
	}

	t.nearest(t.root, p, &currentNearest, &currentRectangle, true)

	return currentNearest.point
}

func (t *KdTree) nearest(n *kdNode, p *Point, currentNearest *distancePair, r *Rectangle, xcoord bool) {
	// Kd-Tree leaf
	if n == nil {
		return
	}

	// Rectangle further than shortest distance, ignore Node and subtree
	if currentNearest.distance < r.distanceSquaredNearest(p) {
		return
	}

	// Check current Node distance
	currentDistance := n.point.distanceSquaredTo(p)
	if currentDistance < currentNearest.distance {
		currentNearest.distance = currentDistance
		currentNearest.point = n.point
	}

	var lesserRectangle, greaterRectangle Rectangle
	if xcoord {
		lesserRectangle = Rectangle{
			r.topLeft,
			Point{n.point.x, r.bottomRight.y},
		}
		greaterRectangle = Rectangle{
			Point{n.point.x, r.topLeft.y},
			r.bottomRight,
		}
	} else {
		lesserRectangle = Rectangle{
			r.topLeft,
			Point{r.bottomRight.x, n.point.y},
		}
		greaterRectangle = Rectangle{
			Point{r.topLeft.x, n.point.y},
			r.bottomRight,
		}
	}

	if lesserRectangle.distanceSquaredNearest(p) < greaterRectangle.distanceSquaredNearest(p) {
		t.nearest(n.left, p, currentNearest, &lesserRectangle, !xcoord)
		t.nearest(n.right, p, currentNearest, &greaterRectangle, !xcoord)
	} else {
		t.nearest(n.right, p, currentNearest, &greaterRectangle, !xcoord)
		t.nearest(n.left, p, currentNearest, &lesserRectangle, !xcoord)
	}
}

func (t *KdTree) Farthest(p *Point) *Point {
	if p == nil {
		return nil
	}

	currentFarthest := distancePair{
		t.root.point,
		t.root.point.distanceSquaredTo(p),
	}
	currentRectangle := Rectangle{
		Point{0.0, 0.0},
		Point{float32(WindowWidth), float32(WindowHeight)},
	}

	t.farthest(t.root, p, &currentFarthest, &currentRectangle, true)

	return currentFarthest.point
}

func (t *KdTree) farthest(n *kdNode, p *Point, currentFarthest *distancePair, r *Rectangle, xcoord bool) {
	// Kd-Tree leaf
	if n == nil {
		return
	}

	// Rectangle closer than longest distance, ignore Node and subtree
	if currentFarthest.distance > r.distanceSquaredFarthest(p) {
		return
	}

	// Check current Node distance
	currentDistance := n.point.distanceSquaredTo(p)
	if currentDistance > currentFarthest.distance {
		currentFarthest.distance = currentDistance
		currentFarthest.point = n.point
	}

	var lesserRectangle, greaterRectangle Rectangle
	if xcoord {
		lesserRectangle = Rectangle{
			r.topLeft,
			Point{n.point.x, r.bottomRight.y},
		}
		greaterRectangle = Rectangle{
			Point{n.point.x, r.topLeft.y},
			r.bottomRight,
		}
	} else {
		lesserRectangle = Rectangle{
			r.topLeft,
			Point{r.bottomRight.x, n.point.y},
		}
		greaterRectangle = Rectangle{
			Point{r.topLeft.x, n.point.y},
			r.bottomRight,
		}
	}

	if lesserRectangle.distanceSquaredFarthest(p) > greaterRectangle.distanceSquaredFarthest(p) {
		t.farthest(n.left, p, currentFarthest, &lesserRectangle, !xcoord)
		t.farthest(n.right, p, currentFarthest, &greaterRectangle, !xcoord)
	} else {
		t.farthest(n.right, p, currentFarthest, &greaterRectangle, !xcoord)
		t.farthest(n.left, p, currentFarthest, &lesserRectangle, !xcoord)
	}
}
