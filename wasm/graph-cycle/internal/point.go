package internal

import (
	"fmt"
)

type QuantizedPoint struct {
	x, y uint
}

func (p *QuantizedPoint) String() string {
	return fmt.Sprintf("(%d, %d)", p.x, p.y)
}
