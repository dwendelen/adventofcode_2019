package main

import (
	"../../util"
	"fmt"
	"log"
	"math"
	"strconv"
)

type KIND int

const (
	START KIND = iota
	EMPTY
	WALL
	OXYGEN
)

const (
	NORTH int64 = 1
	SOUTH       = 2
	WEST        = 3
	EAST        = 4
)

const (
	FB_WALL   int64 = 0
	FB_EMPTY        = 1
	FB_OXYGEN       = 2
)

type Cell struct {
	kind             KIND
	distanceToOxygen int64
}

func main() {
	program := load("day15/input.txt")

	input := make(chan int64, 1)
	output := make(chan int64)
	halt := make(chan bool)

	go run(program, input, output, halt)

	grid := make(map[XY]*Cell)
	grid[XY{0, 0}] = &Cell{START, math.MaxInt64}

	pos := XY{0, 0}
	oxygen := XY{0, 0}

	explore(grid, pos, input, output, &oxygen)

	flood(grid, oxygen, 0)

	maxDistanceFromOxygen := int64(0)
	for _, cell := range grid {
		if cell.kind == EMPTY {
			if cell.distanceToOxygen > maxDistanceFromOxygen {
				maxDistanceFromOxygen = cell.distanceToOxygen
			}
		}
	}
	fmt.Println(maxDistanceFromOxygen)
}

func explore(grid map[XY]*Cell, pos XY, input, output chan int64, oxygen *XY) {
	explore2(grid, XY{pos.x, pos.y + 1}, NORTH, SOUTH, input, output, oxygen)
	explore2(grid, XY{pos.x + 1, pos.y}, EAST, WEST, input, output, oxygen)
	explore2(grid, XY{pos.x, pos.y - 1}, SOUTH, NORTH, input, output, oxygen)
	explore2(grid, XY{pos.x - 1, pos.y}, WEST, EAST, input, output, oxygen)
}

func explore2(grid map[XY]*Cell, pos XY, cmd, cmd2 int64, input, output chan int64, oxygen *XY) {
	_, known := grid[pos]
	if !known {
		input <- cmd
		feedback := <-output
		switch feedback {
		case FB_WALL:
			grid[pos] = &Cell{WALL, math.MaxInt64}
			return
		case FB_EMPTY:
			grid[pos] = &Cell{EMPTY, math.MaxInt64}
		case FB_OXYGEN:
			grid[pos] = &Cell{OXYGEN, math.MaxInt64}
			*oxygen = pos
		default:
			log.Fatalln("Should not happen")
		}

		explore(grid, pos, input, output, oxygen)
		input <- cmd2
		feedback = <-output
		if feedback != FB_EMPTY {
			log.Fatalln("Should not happen")
		}
	}
}

func flood(grid map[XY]*Cell, xy XY, val int64) {
	cell, known := grid[xy]
	if !known {
		return
	}

	if cell.kind == WALL {
		return
	}

	if val < cell.distanceToOxygen {
		cell.distanceToOxygen = val
		flood(grid, XY{xy.x + 1, xy.y}, val+1)
		flood(grid, XY{xy.x - 1, xy.y}, val+1)
		flood(grid, XY{xy.x, xy.y + 1}, val+1)
		flood(grid, XY{xy.x, xy.y - 1}, val+1)
	}
}

type XY struct{ x, y int64 }

func printGrid(grid map[XY]*Cell) {
	maxX := int64(0)
	maxY := int64(0)
	minX := int64(0)
	minY := int64(0)
	for xy, _ := range grid {
		maxX = util.Max(maxX, xy.x)
		minX = util.Min(minX, xy.x)
		maxY = util.Max(maxY, xy.y)
		minY = util.Min(minY, xy.y)
	}

	for y := maxY; y >= minY; y-- {
		for x := minX; x <= maxX; x++ {
			cell, found := grid[XY{x, y}]
			if !found {
				fmt.Print("?")
			} else {
				switch cell.kind {
				case START:
					fmt.Print("S")
				case EMPTY:
					fmt.Print(" ")
				case WALL:
					fmt.Print("#")
				case OXYGEN:
					fmt.Print("X")
				default:
					log.Fatalln("Should not happen")
				}
			}
		}
		fmt.Println()
	}
}

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

func run(program map[int64]int64, input chan int64, output chan int64, halt chan bool) {
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
			halt <- true
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
			val := <-input
			write(values, offset, readValue(values, pos+1), mode1, val)
			pos += 2
		case 4: // output
			val := read(values, offset, readValue(values, pos+1), mode1)
			output <- val
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
