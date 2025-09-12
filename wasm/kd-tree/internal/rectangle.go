package internal

import "math"

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

// Distance to the nearest point in the rectangle
func (r *Rectangle) distanceSquaredNearest(p *Point) float32 {
	var dx, dy float32
	if p.x >= r.topLeft.x && p.x <= r.bottomRight.x {
		dx = 0.0
	} else {
		dxLeft := math.Abs(float64(p.x - r.topLeft.x))
		dxRight := math.Abs(float64(p.x - r.bottomRight.x))
		dx = float32(math.Min(dxLeft, dxRight))
	}
	if p.y >= r.topLeft.y && p.y <= r.bottomRight.y {
		dx = 0.0
	} else {
		dyLeft := math.Abs(float64(p.y - r.topLeft.y))
		dyRight := math.Abs(float64(p.y - r.bottomRight.y))
		dy = float32(math.Min(dyLeft, dyRight))
	}

	return dx*dx + dy*dy
}

// Distance to the farthest point in the rectangle
func (r *Rectangle) distanceSquaredFarthest(p *Point) float32 {
	dxLeft := math.Abs(float64(p.x - r.topLeft.x))
	dxRight := math.Abs(float64(p.x - r.bottomRight.x))
	dx := math.Max(dxLeft, dxRight)

	dyLeft := math.Abs(float64(p.y - r.topLeft.y))
	dyRight := math.Abs(float64(p.y - r.bottomRight.y))
	dy := math.Max(dyLeft, dyRight)

	return float32(dx*dx + dy*dy)
}
