package main

import (
	"../../util"
	"fmt"
	"log"
	"strconv"
)

func main() {
	program := load("day07/input.txt")

	combis := allCombinations()
	bestValue := 0
	for _, combi := range combis {
		value := 0
		for i := 0; i < 5; i++ {
			output := run(program, []int{combi[i], value})
			value = output[0]
		}
		if value > bestValue {
			bestValue = value
		}
	}

	fmt.Println(bestValue)
}

func allCombinations() [][]int {
	result := make([][]int, 0, 128)
	for i0 := 0; i0 < 5; i0++ {
		for i1 := 0; i1 < 5; i1++ {
			if i0 == i1 {
				continue
			}
			for i2 := 0; i2 < 5; i2++ {
				if i2 == i0 || i2 == i1 {
					continue
				}
				for i3 := 0; i3 < 5; i3++ {
					if i3 == i0 || i3 == i1 || i3 == i2 {
						continue
					}
					for i4 := 0; i4 < 5; i4++ {
						if i4 == i0 || i4 == i1 || i4 == i2 || i4 == i3 {
							continue
						}
						result = append(result, []int{i0, i1, i2, i3, i4})
					}
				}
			}
		}
	}
	return result
}

func load(file string) []int {
	values := make([]int, 0, 1024)

	util.ReadCommaSeparated(file, func(in string) error {
		asInt, err := strconv.Atoi(in)
		if err != nil {
			return err
		}
		values = append(values, asInt)
		return nil
	})

	return values
}

func run(program []int, input []int) []int {
	values := make([]int, len(program))
	copy(values, program)

	pos := 0
	output := make([]int, 0)

	for {
		mode3, mode2, mode1, code := parse(values[pos])

		switch code {
		case 99: // halt
			return output
		case 1: // add
			val1 := read(values, values[pos+1], mode1)
			val2 := read(values, values[pos+2], mode2)
			write(values, values[pos+3], mode3, val1+val2)
			pos += 4
			break
		case 2: // multiply
			val1 := read(values, values[pos+1], mode1)
			val2 := read(values, values[pos+2], mode2)
			write(values, values[pos+3], mode3, val1*val2)
			pos += 4
			break
		case 3: // input
			val := input[0]
			input = input[1:]
			write(values, values[pos+1], mode1, val)
			pos += 2
			break
		case 4: // output
			val := read(values, values[pos+1], mode1)
			output = append(output, val)
			pos += 2
			break
		case 5: // jump-if-true
			cond := read(values, values[pos+1], mode1)
			addr := read(values, values[pos+2], mode2)
			if cond != 0 {
				pos = addr
			} else {
				pos += 3
			}
			break
		case 6: // jump-if-false
			cond := read(values, values[pos+1], mode1)
			addr := read(values, values[pos+2], mode2)
			if cond == 0 {
				pos = addr
			} else {
				pos += 3
			}
			break
		case 7: // less than
			val1 := read(values, values[pos+1], mode1)
			val2 := read(values, values[pos+2], mode2)
			var res int
			if val1 < val2 {
				res = 1
			} else {
				res = 0
			}
			write(values, values[pos+3], mode3, res)
			pos += 4
			break
		case 8: // equals
			val1 := read(values, values[pos+1], mode1)
			val2 := read(values, values[pos+2], mode2)
			var res int
			if val1 == val2 {
				res = 1
			} else {
				res = 0
			}
			write(values, values[pos+3], mode3, res)
			pos += 4
			break
		default:
			log.Fatal("Unknown op code " + strconv.Itoa(code))
		}
	}
}

func read(values []int, nb int, mode int) int {
	switch mode {
	case 0:
		return values[nb]
	case 1:
		return nb
	default:
		log.Fatal("Illegal mode " + strconv.Itoa(mode))
		return 0
	}
}

func write(values []int, nb int, mode int, val int) {
	switch mode {
	case 0:
		values[nb] = val
	case 1:
		log.Fatal("You can not use immediate mode with a store")
	default:
		log.Fatal("Illegal mode " + strconv.Itoa(mode))
	}
}

func parse(instruction int) (mode3, mode2, mode1, code int) {
	if instruction >= 100000 {
		log.Fatal("Illegal instruction " + strconv.Itoa(instruction))
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
