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
	initialValues := make([]int, 0, 256)

	util.ReadCommaSeparated("day02/input.txt", func(in string) error {
		asInt, err := strconv.Atoi(in)
		if err != nil {
			return err
		}
		initialValues = append(initialValues, asInt)
		return nil
	})

	for i1 := 0; i1 < 100; i1++ {
		for i2 := 0; i2 < 100; i2++ {
			values := make([]int, len(initialValues))
			copy(values, initialValues)

			values[1] = i1
			values[2] = i2

			res, newVals := run(values)
			if res == HALT && newVals[0] == 19690720 {
				fmt.Println("i1:", i1, "i2:", i2)
				return
			}
		}
	}


	fmt.Println("Nothing found!")
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