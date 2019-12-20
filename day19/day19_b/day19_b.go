package main

import (
	"../../util"
	"fmt"
	"log"
	"strconv"
)

type XY struct {
	x, y int64
}

func main() {
	program := load("day19/input.txt")

	var x128 int64
	for x128 = 0; !readBeam(program, x128, 128); x128++ {
	}
	x128-- // margin

	if readBeam(program, x128, 128) {
		log.Fatalln("Should not happen")
	}
	if testLowerLeft(program, x128+1, 128) {
		log.Fatalln("Should not happen")
	}
	if !testLowerLeft(program, findFirstX(program, x128, 128<<8), 128<<8) {
		log.Fatalln("Should not happen")
	}

	lowerY := int64(128)
	upperY := int64(128 << 8)
	for {
		nextY := (upperY + lowerY) / 2
		if nextY == lowerY {
			x := findFirstX(program, x128, upperY)
			y := upperY - 99
			fmt.Println(x, y, 10000*x+y)
			return
		}
		x := findFirstX(program, x128, nextY)
		if testLowerLeft(program, x, nextY) {
			upperY = nextY
		} else {
			lowerY = nextY
		}
	}
}

func testLowerLeft(program map[int64]int64, x, y int64) bool {
	return readBeam(program, x, y) &&
		readBeam(program, x+99, y) &&
		readBeam(program, x, y-99) &&
		readBeam(program, x+99, y-99)
}

func findFirstX(program map[int64]int64, x128, y int64) int64 {
	var x int64
	for x = (y * x128) / 128; !readBeam(program, x, y); x++ {
	}
	return x
}

func readBeam(program map[int64]int64, x int64, y int64) bool {
	input := make(chan int64)
	output := make(chan int64)

	go run(program, input, output)

	input <- x
	input <- y
	outp := <-output
	_, more := <-output
	if more {
		log.Fatalln("Should not happen")
	}

	switch outp {
	case 0:
		return false
	case 1:
		return true
	default:
		log.Fatalln("Should not happen")
		return false
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

func run(program map[int64]int64, input chan int64, output chan int64) {
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
			close(output)
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
