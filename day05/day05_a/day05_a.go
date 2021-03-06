package main

import (
	"../../util"
	"fmt"
	"log"
	"strconv"
)

func main() {
	program := load("day05/input.txt")

	output := run(program, []int{1})

	fmt.Println(output)
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
		case 99:
			return output
		case 1:
			val1 := read(values, values[pos+1], mode1)
			val2 := read(values, values[pos+2], mode2)
			write(values, values[pos+3], mode3, val1+val2)
			pos += 4
			break
		case 2:
			val1 := read(values, values[pos+1], mode1)
			val2 := read(values, values[pos+2], mode2)
			write(values, values[pos+3], mode3, val1*val2)
			pos += 4
			break
		case 3:
			val := input[0]
			input = input[1:]
			write(values, values[pos+1], mode1, val)
			pos += 2
			break
		case 4:
			val := read(values, values[pos+1], mode1)
			output = append(output, val)
			pos += 2
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
