package main

import (
	"fmt"
	"strconv"
	"../../util"
)

func main() {
	fmt.Println(fuel(12), fuel(14), fuel(1969), fuel(100756))

	sum := 0
	util.ReadLines("day01/input.txt", func(in string) error {
		mass, err := strconv.Atoi(in)
		if err != nil {
			return err
		}
		sum += fuel(mass)
		return nil
	})

	fmt.Println(sum)
}

func fuel(mass int) int {
	return (mass / 3) - 2
}

