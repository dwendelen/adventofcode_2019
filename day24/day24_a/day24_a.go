package main

import (
	"../../util"
	"fmt"
	"log"
)

func main() {
	file := "day24/input.txt"
	//file := "day24/day24_a/example.txt"

	grid := readInitial(file)
	buffer := make([][]bool, len(grid))
	for i := range buffer {
		buffer[i] = make([]bool, len(grid[i]))
	}

	ratings := make(map[int]bool)
	ratings[calcRating(grid)] = true
	for {

		stepAndSwap(&grid, &buffer)
		rat := calcRating(grid)
		_, found := ratings[rat]
		if found {
			fmt.Println(rat)
			break
		}
		ratings[rat] = true
	}
}

func stepAndSwap(from, to *[][]bool) {
	for y := 0; y < len(*from); y++ {
		row := (*from)[y]
		for x := 0; x < len(row); x++ {
			neighbours := nbNeighbours(*from, x, y)
			var newVal bool
			if read(*from, x, y) {
				newVal = neighbours == 1
			} else {
				newVal = neighbours == 1 || neighbours == 2
			}
			(*to)[y][x] = newVal
		}
	}
	*from, *to = *to, *from
}

func calcRating(grid [][]bool) int {
	rating := 0
	power := 1
	for _, row := range grid {
		for _, cell := range row {
			if cell {
				rating += power
			}
			power *= 2
		}
	}
	return rating
}

func nbNeighbours(grid [][]bool, x, y int) int {
	res := 0
	if read(grid, x+1, y) {
		res++
	}
	if read(grid, x-1, y) {
		res++
	}
	if read(grid, x, y+1) {
		res++
	}
	if read(grid, x, y-1) {
		res++
	}
	return res
}

func read(grid [][]bool, x, y int) bool {
	if y < 0 || y >= len(grid) {
		return false
	}

	row := grid[y]

	if x < 0 || x >= len(row) {
		return false
	}

	return row[x]
}

func printGrid(grid [][]bool) {
	for _, row := range grid {
		for _, cell := range row {
			if cell {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func readInitial(file string) [][]bool {
	rows := make([][]bool, 0)
	util.ReadLines(file, func(in string) error {
		row := make([]bool, len(in))
		for i, cell := range in {
			var val bool
			switch cell {
			case '.':
				val = false
			case '#':
				val = true
			default:
				log.Fatalln("Should not happen")
			}
			row[i] = val
		}
		rows = append(rows, row)
		return nil
	})
	return rows
}
