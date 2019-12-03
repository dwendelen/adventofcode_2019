package main

import (
	"../../util"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Coord struct {
	x     int
	y     int
	steps int
}

type Segment struct {
	from        Coord
	to          Coord
	orientation Orientation
}
type Orientation string

const (
	HORIZONTAL = "-"
	VERTICAL   = "|"
)

func main() {
	var segments1 []Segment
	var segments2 []Segment

	wire := 0

	util.ReadLines("day03/input.txt", func(in string) error {
		path := strings.Split(in, ",")

		segments := make([]Segment, 0, 512)

		point1 := Coord{0, 0, 0}
		for _, instr := range path {
			length, err := strconv.Atoi(instr[1:])
			if err != nil {
				return err
			}

			var point2 Coord
			var orientation Orientation
			switch instr[0] {
			case 'U':
				point2.x = point1.x
				point2.y = point1.y + length
				point2.steps = point1.steps + length
				orientation = VERTICAL
				break
			case 'D':
				point2.x = point1.x
				point2.y = point1.y - length
				point2.steps = point1.steps + length
				orientation = VERTICAL
				break
			case 'L':
				point2.x = point1.x - length
				point2.y = point1.y
				point2.steps = point1.steps + length
				orientation = HORIZONTAL
				break
			case 'R':
				point2.x = point1.x + length
				point2.y = point1.y
				point2.steps = point1.steps + length
				orientation = HORIZONTAL
				break
			default:
				return errors.New("Unknown command " + instr[0:1])
			}

			segments = append(segments, Segment{point1, point2, orientation})
			point1 = point2
		}

		switch wire {
		case 0:
			segments1 = segments
			break
		case 1:
			segments2 = segments
			break
		default:
			return errors.New("More then 2 wires")
		}

		wire++
		return nil
	})

	closest := 0xefffffff
	for _, segment1 := range segments1 {
		for _, segment2 := range segments2 {
			coll := collision(segment1, segment2)
			if coll != 0 {
				if coll < closest {
					closest = coll
				}
			}
		}
	}

	fmt.Println(closest)
}

func collision(seg1, seg2 Segment) int {
	if seg1.orientation == seg2.orientation {
		return 0
	}

	var hor, ver Segment
	if seg1.orientation == HORIZONTAL {
		hor, ver = seg1, seg2
	} else {
		hor, ver = seg2, seg1
	}

	if between(hor.from.y, ver.from.y, ver.to.y) &&
		between(ver.from.x, hor.from.x, hor.to.x) {

		steps1 := hor.from.steps + abs(hor.from.x-ver.from.x)
		steps2 := ver.from.steps + abs(ver.from.y-hor.from.y)

		return steps1 + steps2
	} else {
		return 0
	}
}

func between(toTest, a, b int) bool {
	return toTest <= max(a, b) && toTest >= min(a, b)
}

func max(a, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	} else {
		return a
	}
}
