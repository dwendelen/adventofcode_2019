package main

import (
	"../../util"
	"fmt"
	"log"
	"math"
	"strconv"
)

func main() {
	program := load("day13/input.txt")

	program[0] = 2
	score := int64(0)

	grid := make(map[XY]int64)

	run(program, getInput(grid), onOutput(grid, &score))

	for _, cell := range grid {
		if cell == 2 {
			log.Fatal("There should not be any blocks left")
		}
	}
	fmt.Println(score)
}

func onOutput(grid map[XY]int64, score *int64) func(int64) {
	outputs := make([]int64, 0, 3)
	return func(out int64) {
		outputs = append(outputs, out)
		if len(outputs) >= 3 {
			if outputs[0] == -1 && outputs[1] == 0 {
				*score = outputs[2]
			} else {
				grid[XY{outputs[0], outputs[1]}] = outputs[2]
			}
			outputs = outputs[3:]
		}
	}
}

func getInput(grid map[XY]int64) func() int64 {
	return func() int64 {
		minPad := int64(math.MaxInt64)
		maxPad := int64(math.MinInt64)
		ballFound := false
		ball := XY{0, 0}

		for xy, cell := range grid {
			switch cell {
			case 3:
				minPad = util.Min(minPad, xy.x)
				maxPad = util.Max(maxPad, xy.x)
			case 4:
				if ballFound {
					log.Fatal("Only one ball expected")
				}
				ballFound = true
				ball = xy
			}
		}
		if !ballFound {
			return 0
		}
		middle := (minPad + maxPad) / 2
		switch {
		case middle == ball.x:
			return 0
		case middle < ball.x:
			return 1
		case middle > ball.x:
			return -1
		}
		return 0
	}
}

func getColor(hull map[XY]int64, pos XY) int64 {
	color, found := hull[pos]
	if found {
		return color
	} else {
		return 0
	}
}

type XY struct{ x, y int64 }

func load(file string) map[int64]int64 {
	values := make(map[int64]int64)
	i := int64(0)

	util.ReadCommaSeparated(file, func(in string) error {
		asInt, err := strconv.Atoi(in)
		if err != nil {
			return err
		}
		values[i] = int64(asInt)
		i++
		return nil
	})

	return values
}

func run(program map[int64]int64, input func() int64, output func(int64)) {
	values := make(map[int64]int64)
	for k, v := range program {
		values[k] = v
	}

	pos := int64(0)
	offset := int64(0)

	for {
		mode3, mode2, mode1, code := parse(readValue(values, pos))

		switch code {
		case 99: // halt
			return
		case 1: // add
			val1 := read(values, offset, readValue(values, pos+1), mode1)
			val2 := read(values, offset, readValue(values, pos+2), mode2)
			write(values, offset, readValue(values, pos+3), mode3, val1+val2)
			pos += 4
		case 2: // multiply
			val1 := read(values, offset, readValue(values, pos+1), mode1)
			val2 := read(values, offset, readValue(values, pos+2), mode2)
			write(values, offset, readValue(values, pos+3), mode3, val1*val2)
			pos += 4
		case 3: // input
			val := input()
			write(values, offset, readValue(values, pos+1), mode1, val)
			pos += 2
		case 4: // output
			val := read(values, offset, readValue(values, pos+1), mode1)
			output(val)
			pos += 2
		case 5: // jump-if-true
			cond := read(values, offset, readValue(values, pos+1), mode1)
			addr := read(values, offset, readValue(values, pos+2), mode2)
			if cond != 0 {
				pos = addr
			} else {
				pos += 3
			}
		case 6: // jump-if-false
			cond := read(values, offset, readValue(values, pos+1), mode1)
			addr := read(values, offset, readValue(values, pos+2), mode2)
			if cond == 0 {
				pos = addr
			} else {
				pos += 3
			}
		case 7: // less than
			val1 := read(values, offset, readValue(values, pos+1), mode1)
			val2 := read(values, offset, readValue(values, pos+2), mode2)
			var res int64
			if val1 < val2 {
				res = 1
			} else {
				res = 0
			}
			write(values, offset, readValue(values, pos+3), mode3, res)
			pos += 4
		case 8: // equals
			val1 := read(values, offset, readValue(values, pos+1), mode1)
			val2 := read(values, offset, readValue(values, pos+2), mode2)
			var res int64
			if val1 == val2 {
				res = 1
			} else {
				res = 0
			}
			write(values, offset, readValue(values, pos+3), mode3, res)
			pos += 4
		case 9: //adjust offset
			val := read(values, offset, readValue(values, pos+1), mode1)
			offset += val
			pos += 2
		default:
			log.Fatal("Unknown op code " + string(code))
		}
	}
}

func readValue(values map[int64]int64, nb int64) int64 {
	res, ok := values[nb]
	if ok {
		return res
	} else {
		return 0
	}
}

func read(values map[int64]int64, offset int64, nb int64, mode int64) int64 {
	switch mode {
	case 0:
		return readValue(values, nb)
	case 1:
		return nb
	case 2:
		return readValue(values, nb+offset)
	default:
		log.Fatal("Illegal mode " + string(mode))
		return 0
	}
}

func write(values map[int64]int64, offset int64, nb int64, mode int64, val int64) {
	switch mode {
	case 0:
		values[nb] = val
	case 1:
		log.Fatal("You can not use immediate mode with a store")
	case 2:
		values[nb+offset] = val
	default:
		log.Fatal("Illegal mode " + string(mode))
	}
}

func parse(instruction int64) (mode3, mode2, mode1, code int64) {
	if instruction >= 100000 {
		log.Fatal("Illegal instruction " + string(instruction))
	}
	code = instruction % 100
	rest := instruction / 100
	mode1 = rest % 10
	rest = rest / 10
	mode2 = rest % 10
	rest = rest / 10
	mode3 = rest % 10

	return
}
