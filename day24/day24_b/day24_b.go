package main

import (
	"../../util"
	"fmt"
	"log"
	"math"
)

func main() {
	file := "day24/input.txt"
	times := 200
	//file := "day24/day24_b/example.txt"
	//times := 10

	grid := readInitial(file)

	for i := 0; i < times; i++ {
		stepAndSwap(&grid)
	}

	res := 0
	for _, cell := range grid {
		if cell {
			res++
		}
	}
	fmt.Println(res)
}

func stepAndSwap(grid *map[XYD]bool) {
	buffer := make(map[XYD]bool)

	for xyd, cell := range *grid {
		neighbours := 0

		switch xyd.y {
		case 0:
			visit(*grid, buffer, 2, 1, xyd.depth-1, cell, &neighbours)
			visit(*grid, buffer, xyd.x, 1, xyd.depth, cell, &neighbours)
		case 1:
			visit(*grid, buffer, xyd.x, 0, xyd.depth, cell, &neighbours)
			if xyd.x == 2 {
				visit(*grid, buffer, 0, 0, xyd.depth+1, cell, &neighbours)
				visit(*grid, buffer, 1, 0, xyd.depth+1, cell, &neighbours)
				visit(*grid, buffer, 2, 0, xyd.depth+1, cell, &neighbours)
				visit(*grid, buffer, 3, 0, xyd.depth+1, cell, &neighbours)
				visit(*grid, buffer, 4, 0, xyd.depth+1, cell, &neighbours)
			} else {
				visit(*grid, buffer, xyd.x, 2, xyd.depth, cell, &neighbours)
			}
		case 2:
			visit(*grid, buffer, xyd.x, 1, xyd.depth, cell, &neighbours)
			visit(*grid, buffer, xyd.x, 3, xyd.depth, cell, &neighbours)
		case 3:
			if xyd.x == 2 {
				visit(*grid, buffer, 0, 4, xyd.depth+1, cell, &neighbours)
				visit(*grid, buffer, 1, 4, xyd.depth+1, cell, &neighbours)
				visit(*grid, buffer, 2, 4, xyd.depth+1, cell, &neighbours)
				visit(*grid, buffer, 3, 4, xyd.depth+1, cell, &neighbours)
				visit(*grid, buffer, 4, 4, xyd.depth+1, cell, &neighbours)
			} else {
				visit(*grid, buffer, xyd.x, 2, xyd.depth, cell, &neighbours)
			}
			visit(*grid, buffer, xyd.x, 4, xyd.depth, cell, &neighbours)

		case 4:
			visit(*grid, buffer, xyd.x, 3, xyd.depth, cell, &neighbours)
			visit(*grid, buffer, 2, 3, xyd.depth-1, cell, &neighbours)
		default:
			log.Fatalln("Error")
		}

		switch xyd.x {
		case 0:
			visit(*grid, buffer, 1, 2, xyd.depth-1, cell, &neighbours)
			visit(*grid, buffer, 1, xyd.y, xyd.depth, cell, &neighbours)
		case 1:
			visit(*grid, buffer, 0, xyd.y, xyd.depth, cell, &neighbours)
			if xyd.y == 2 {
				visit(*grid, buffer, 0, 0, xyd.depth+1, cell, &neighbours)
				visit(*grid, buffer, 0, 1, xyd.depth+1, cell, &neighbours)
				visit(*grid, buffer, 0, 2, xyd.depth+1, cell, &neighbours)
				visit(*grid, buffer, 0, 3, xyd.depth+1, cell, &neighbours)
				visit(*grid, buffer, 0, 4, xyd.depth+1, cell, &neighbours)
			} else {
				visit(*grid, buffer, 2, xyd.y, xyd.depth, cell, &neighbours)
			}
		case 2:
			visit(*grid, buffer, 1, xyd.y, xyd.depth, cell, &neighbours)
			visit(*grid, buffer, 3, xyd.y, xyd.depth, cell, &neighbours)
		case 3:
			if xyd.y == 2 {
				visit(*grid, buffer, 4, 0, xyd.depth+1, cell, &neighbours)
				visit(*grid, buffer, 4, 1, xyd.depth+1, cell, &neighbours)
				visit(*grid, buffer, 4, 2, xyd.depth+1, cell, &neighbours)
				visit(*grid, buffer, 4, 3, xyd.depth+1, cell, &neighbours)
				visit(*grid, buffer, 4, 4, xyd.depth+1, cell, &neighbours)
			} else {
				visit(*grid, buffer, 2, xyd.y, xyd.depth, cell, &neighbours)
			}
			visit(*grid, buffer, 4, xyd.y, xyd.depth, cell, &neighbours)
		case 4:
			visit(*grid, buffer, 3, xyd.y, xyd.depth, cell, &neighbours)
			visit(*grid, buffer, 3, 2, xyd.depth-1, cell, &neighbours)
		default:
			log.Fatalln("Error")
		}

		var newVal bool
		if cell {
			newVal = neighbours == 1
		} else {
			newVal = neighbours == 1 || neighbours == 2
		}
		buffer[xyd] = newVal
	}
	*grid = buffer
}

func visit(grid map[XYD]bool, buffer map[XYD]bool, x, y, depth int, expand bool, acc *int) {
	key := XYD{x, y, depth}

	cell, found := grid[key]
	if !found {
		cell = false
	}
	if cell {
		*acc++
	}

	//if expand {
	_, found = buffer[key]
	if !found {
		buffer[key] = false
	}
	//}
}

func printGrid(grid map[XYD]bool) {
	minDepth := math.MaxInt32
	maxDepth := -math.MaxInt32
	for xyd, _ := range grid {
		minDepth = util.MinInt(minDepth, xyd.depth)
		maxDepth = util.MaxInt(maxDepth, xyd.depth)
	}

	for d := minDepth; d <= maxDepth; d++ {
		fmt.Println("Depth", d)
		for y := 0; y < 5; y++ {
			for x := 0; x < 5; x++ {
				cell, found := grid[XYD{x, y, d}]
				if !found {
					fmt.Print("?")
				} else if cell {
					fmt.Print("#")
				} else {
					fmt.Print(".")
				}
			}
			fmt.Println()
		}
		fmt.Println()
	}
}

func readInitial(file string) map[XYD]bool {
	grid := make(map[XYD]bool)
	y := 0
	util.ReadLines(file, func(in string) error {
		for x, cell := range in {
			var val bool
			switch cell {
			case '.':
				val = false
			case '#':
				val = true
			default:
				log.Fatalln("Should not happen")
			}
			grid[XYD{x, y, 0}] = val
		}
		y++
		return nil
	})

	grid[XYD{1, 2, -1}] = false
	grid[XYD{3, 2, -1}] = false
	grid[XYD{2, 1, -1}] = false
	grid[XYD{2, 3, -1}] = false

	grid[XYD{0, 0, 1}] = false
	grid[XYD{1, 0, 1}] = false
	grid[XYD{2, 0, 1}] = false
	grid[XYD{3, 0, 1}] = false
	grid[XYD{4, 0, 1}] = false

	grid[XYD{0, 1, 1}] = false
	grid[XYD{4, 1, 1}] = false
	grid[XYD{0, 2, 1}] = false
	grid[XYD{4, 2, 1}] = false
	grid[XYD{0, 3, 1}] = false
	grid[XYD{4, 3, 1}] = false

	grid[XYD{0, 4, 1}] = false
	grid[XYD{1, 4, 1}] = false
	grid[XYD{2, 4, 1}] = false
	grid[XYD{3, 4, 1}] = false
	grid[XYD{4, 4, 1}] = false

	delete(grid, XYD{2, 2, 0})
	return grid
}

type XYD struct {
	x, y, depth int
}
