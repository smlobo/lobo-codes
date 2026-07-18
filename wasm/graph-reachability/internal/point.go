package internal

import (
	"fmt"
)

type QuantizedPoint struct {
	x int
	y int
}

func (p *QuantizedPoint) String() string {
	return fmt.Sprintf("(%d, %d)", p.x, p.y)
}
