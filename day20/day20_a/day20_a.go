package main

import (
	"../../util"
	"fmt"
	"log"
	"math"
)

func main() {
	file := "day20/input.txt"

	lines := util.ReadAllLines(file)
	start, transports, end := processLabels(lines)
	grid := loadGrid(lines)

	flood(grid, transports, start.x, start.y, 0)

	fmt.Println(grid[end].distance)
}

func flood(grid map[XY]*Cell, transports map[XY]XY, x int64, y int64, distance int64) {
	cell, found := grid[XY{x, y}]
	if !found {
		return
	}

	if distance >= cell.distance {
		return
	}

	cell.distance = distance
	flood(grid, transports, x-1, y, distance+1)
	flood(grid, transports, x+1, y, distance+1)
	flood(grid, transports, x, y-1, distance+1)
	flood(grid, transports, x, y+1, distance+1)

	dest, trans := transports[XY{x, y}]
	if trans {
		flood(grid, transports, dest.x, dest.y, distance+1)
	}
}

func loadGrid(lines []string) map[XY]*Cell {
	grid := make(map[XY]*Cell)

	sizeY := int64(len(lines)) - 4
	sizeX := int64(len(lines[2])) - 2

	for y := int64(0); y < sizeY; y++ {
		for x := int64(0); x < sizeX; x++ {
			if lines[2+y][2+x] == '.' {
				grid[XY{x, y}] = &Cell{math.MaxInt64}
			}
		}
	}

	return grid
}

func processLabels(lines []string) (XY, map[XY]XY, XY) {
	sizeY := int64(len(lines)) - 4
	sizeX := int64(len(lines[2])) - 2

	var thickness int64
	for x := int64(2); ; x++ {
		char := lines[2+sizeY/2][x]
		if char != '.' && char != '#' {
			thickness = x - 2
			break
		}
	}

	labels := make(map[string][]XY)

	for x := int64(0); x < sizeX; x++ {
		maybeAddLabel(labels, lines, x, 0, 0, -1)
		if thickness-1 <= x && x <= sizeX-thickness {
			maybeAddLabel(labels, lines, x, thickness-1, 0, 1)
			maybeAddLabel(labels, lines, x, sizeY-thickness, 0, -1)
		}
		maybeAddLabel(labels, lines, x, sizeY-1, 0, 1)
	}

	for y := int64(0); y < sizeY; y++ {
		maybeAddLabel(labels, lines, 0, y, -1, 0)
		if thickness-1 <= y && y <= sizeY-thickness {
			maybeAddLabel(labels, lines, thickness-1, y, 1, 0)
			maybeAddLabel(labels, lines, sizeX-thickness, y, -1, 0)
		}
		maybeAddLabel(labels, lines, sizeX-1, y, 1, 0)
	}

	var start, end XY
	transports := make(map[XY]XY)
	for lbl, arr := range labels {
		switch lbl {
		case "AA":
			if len(arr) != 1 {
				log.Fatalln("Should not happen, probably wrong thickness")
			}
			start = arr[0]
		case "ZZ":
			if len(arr) != 1 {
				log.Fatalln("Should not happen, probably wrong thickness")
			}
			end = arr[0]
		default:
			if len(arr) != 2 {
				log.Fatalln("Should not happen, probably wrong thickness")
			}
			transports[arr[0]] = arr[1]
			transports[arr[1]] = arr[0]
		}
	}
	return start, transports, end
}

func maybeAddLabel(labels map[string][]XY, lines []string, x int64, y int64, dx int64, dy int64) {
	if lines[y+2][x+2] != '.' {
		return
	}

	var label string
	if dx+dy < 0 {
		label = string([]byte{lines[y+2+dy+dy][x+2+dx+dx], lines[y+2+dy][x+2+dx]})
	} else {
		label = string([]byte{lines[y+2+dy][x+2+dx], lines[y+2+dy+dy][x+2+dx+dx]})
	}

	arr, found := labels[label]
	if !found {
		arr = make([]XY, 0)
	}
	labels[label] = append(arr, XY{x, y})
}

type XY struct {
	x, y int64
}

type Cell struct {
	distance int64
}
