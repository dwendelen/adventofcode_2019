package main

import (
	"../../util"
	"fmt"
	"strconv"
)

type result string
const HALT = "HALT"
const ERROR = "ERROR"

func main() {
	values := make([]int, 0, 256)

	util.ReadCommaSeparated("day02/input.txt", func(in string) error {
		asInt, err := strconv.Atoi(in)
		if err != nil {
			return err
		}
		values = append(values, asInt)
		return nil
	})

	values[1] = 12
	values[2] = 2

	res, newVals := run(values)

	fmt.Println(res, newVals)
}

func run(input []int) (result, []int) {
	values := make([]int, len(input))
	copy(values, input)

	pos := 0

	for true {
		code := values[pos]
		switch code {
		case 99:
			return HALT, values
		case 1:
			values[values[pos + 3]] = values[values[pos + 1]] + values[values[pos + 2]]
			pos += 4
			break
		case 2:
			values[values[pos + 3]] = values[values[pos + 1]] * values[values[pos + 2]]
			pos += 4
			break
		default:
			return ERROR, values
		}
	}
	return ERROR, values
}