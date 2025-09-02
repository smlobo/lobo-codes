package internal

type Rectangle struct {
	topLeft, bottomRight Point
}

func (r *Rectangle) intersects(o *Rectangle) bool {
	if o == nil {
		return false
	}

	return r.topLeft.x <= o.bottomRight.x && r.bottomRight.x >= o.topLeft.x &&
		r.topLeft.y <= o.bottomRight.y && r.bottomRight.y >= o.topLeft.y
}

func (r *Rectangle) distanceSquaredTo(p *Point) float32 {
	var dx, dy float32
	if p.x < r.topLeft.x {
		dx = p.x - r.topLeft.x
	} else if p.x > r.bottomRight.x {
		dx = p.x - r.bottomRight.x
	}
	if p.y < r.topLeft.y {
		dy = p.y - r.topLeft.y
	} else if p.y > r.bottomRight.y {
		dy = p.y - r.bottomRight.y
	}
	return dx*dx + dy*dy
}
