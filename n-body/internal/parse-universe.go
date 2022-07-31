package internal

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func ParseUniverse(file *os.File) {
	var (
		buf []byte
	)

	reader := bufio.NewReader(file)
	buf, _, _ = reader.ReadLine()
	NumBodies, _ = strconv.Atoi(strings.TrimSpace(string(buf)))
	buf, _, _ = reader.ReadLine()
	Radius, _ = strconv.ParseFloat(strings.TrimSpace(string(buf)), 64)
	Bodies = make([]body, NumBodies)
	for i := 0; i < NumBodies; i++ {
		buf, _, _ = reader.ReadLine()
		bodyInfo := strings.Fields(string(buf))
		xL, _ := strconv.ParseFloat(bodyInfo[0], 64)
		yL, _ := strconv.ParseFloat(bodyInfo[1], 64)
		xV, _ := strconv.ParseFloat(bodyInfo[2], 64)
		yV, _ := strconv.ParseFloat(bodyInfo[3], 64)
		m, _ := strconv.ParseFloat(bodyInfo[4], 64)
		Bodies[i] = body{
			xCoord:    xL,
			yCoord:    yL,
			xVelocity: xV,
			yVelocity: yV,
			mass:      m,
			imageFile: strings.Replace(strings.TrimSpace(bodyInfo[5]), ".gif", ".png", 1),
		}
	}

}
