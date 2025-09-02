package internal

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
	if currentNearest.distance < r.distanceSquaredTo(p) {
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

	if lesserRectangle.distanceSquaredTo(p) < r.distanceSquaredTo(p) {
		t.nearest(n.left, p, currentNearest, &lesserRectangle, !xcoord)
		t.nearest(n.right, p, currentNearest, &greaterRectangle, !xcoord)
	} else {
		t.nearest(n.right, p, currentNearest, &greaterRectangle, !xcoord)
		t.nearest(n.left, p, currentNearest, &lesserRectangle, !xcoord)
	}
}
