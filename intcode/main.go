package main

import (
	"../util"
	"fmt"
	"log"
	"strconv"
)

func main() {
	program := load("day09/input.txt")
	//program[6] = program[6] + 5
	printProg(program)
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

func printProg(program []int) {
	pos := 0

	for pos < len(program) {
		fmt.Printf("% 4d ", pos)
		fmt.Printf("%05d ", program[pos])

		mode3, mode2, mode1, code := parse(program[pos])

		switch code {
		case 99: // halt
			fmt.Printf("HLT")
		case 1: // add
			fmt.Printf("ADD ")
			printArg(program[pos+1], mode1)
			printArg(program[pos+2], mode2)
			printArg(program[pos+3], mode3)
			pos += 3
		case 2: // multiply
			fmt.Printf("MUL ")
			printArg(program[pos+1], mode1)
			printArg(program[pos+2], mode2)
			printArg(program[pos+3], mode3)
			pos += 3
		case 3: // input
			fmt.Printf("INP ")
			printArg(program[pos+1], mode1)
			pos += 1
		case 4: // output
			fmt.Printf("OUT ")
			printArg(program[pos+1], mode1)
			pos += 1
		case 5: // jump-if-true
			fmt.Printf("JPT ")
			printArg(program[pos+1], mode1)
			printArg(program[pos+2], mode2)
			pos += 2
		case 6: // jump-if-false
			fmt.Printf("JPF ")
			printArg(program[pos+1], mode1)
			printArg(program[pos+2], mode2)
			pos += 2
		case 7: // less than
			fmt.Printf("LES ")
			printArg(program[pos+1], mode1)
			printArg(program[pos+2], mode2)
			printArg(program[pos+3], mode3)
			pos += 3
		case 8: // equals
			fmt.Printf("EQL ")
			printArg(program[pos+1], mode1)
			printArg(program[pos+2], mode2)
			printArg(program[pos+3], mode3)
			pos += 3
		case 9: // equals
			fmt.Printf("REL ")
			printArg(program[pos+1], mode1)
			pos += 1
		default:
		}
		fmt.Printf("\n")
		pos++
	}
}
func printArg(val int, mode int) {
	var str string
	switch mode {
	case 0:
		str = "*"
	case 1:
	case 2:
		str = "%"
	}

	str += strconv.Itoa(val)
	fmt.Printf("%12s ", str)
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
