package main

import (
	"../../util"
	"container/heap"
	"fmt"
	"math"
)

func main() {
	file := "day20/input.txt"

	lines := util.ReadAllLines(file)
	start, outerToInner, end := processLabels(lines)
	grid := loadGrid(lines)

	routes := loadRoutes(grid, start, outerToInner, end)

	paths := &PathHeap{make([]Path, 0)}
	heap.Init(paths)
	heap.Push(paths, Path{start, 0, 0})

	var res int64
	for {
		path := heap.Pop(paths).(Path)
		if path.currentLevel == 0 && path.current == end {
			res = path.distanceUntilNow //To undo the one from the jump
			break
		}

		if path.currentLevel == 0 {
			for to, route := range routes[path.current] {
				if to == end || route.dLevel > 0 {
					heap.Push(paths, Path{route.to,
						path.currentLevel + route.dLevel,
						path.distanceUntilNow + route.distance,
					})
				}
			}
		} else {
			for to, route := range routes[path.current] {
				if to != end {
					heap.Push(paths, Path{route.to,
						path.currentLevel + route.dLevel,
						path.distanceUntilNow + route.distance,
					})
				}
			}
		}

	}
	fmt.Println(res)
}

func loadRoutes(grid map[XY]*Cell, start XY, outerToInner map[XY]XY, end XY) map[XY]map[XY]Route {
	routes := make(map[XY]map[XY]Route)
	loadRoute(grid, outerToInner, end, start, routes)
	for outer, inner := range outerToInner {
		loadRoute(grid, outerToInner, end, outer, routes)
		loadRoute(grid, outerToInner, end, inner, routes)
	}
	return routes
}

func loadRoute(grid map[XY]*Cell, outerToInner map[XY]XY, end XY, from XY, routes map[XY]map[XY]Route) {
	resetGrid(grid)
	flood(grid, from.x, from.y, 0)

	newRoutes := make(map[XY]Route)
	cell := grid[end]
	if cell.distance != math.MaxInt64 {
		newRoutes[end] = Route{cell.distance, end, 0}
	}
	for outer, inner := range outerToInner {
		cell = grid[outer]
		if cell.distance != math.MaxInt64 {
			newRoutes[outer] = Route{cell.distance + 1, inner, -1}
		}

		cell = grid[inner]
		if cell.distance != math.MaxInt64 {
			newRoutes[inner] = Route{cell.distance + 1, outer, 1}
		}
	}
	delete(newRoutes, from)
	routes[from] = newRoutes
}

func flood(grid map[XY]*Cell, x int64, y int64, distance int64) {
	cell, found := grid[XY{x, y}]
	if !found {
		return
	}

	if distance >= cell.distance {
		return
	}

	cell.distance = distance
	flood(grid, x-1, y, distance+1)
	flood(grid, x+1, y, distance+1)
	flood(grid, x, y-1, distance+1)
	flood(grid, x, y+1, distance+1)
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

func resetGrid(grid map[XY]*Cell) {
	for _, cell := range grid {
		cell.distance = math.MaxInt64
	}
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
		maybeAddLabel(labels, lines, x, 0, 0, -1, true)
		if thickness-1 <= x && x <= sizeX-thickness {
			maybeAddLabel(labels, lines, x, thickness-1, 0, 1, false)
			maybeAddLabel(labels, lines, x, sizeY-thickness, 0, -1, false)
		}
		maybeAddLabel(labels, lines, x, sizeY-1, 0, 1, true)
	}

	for y := int64(0); y < sizeY; y++ {
		maybeAddLabel(labels, lines, 0, y, -1, 0, true)
		if thickness-1 <= y && y <= sizeY-thickness {
			maybeAddLabel(labels, lines, thickness-1, y, 1, 0, false)
			maybeAddLabel(labels, lines, sizeX-thickness, y, -1, 0, false)
		}
		maybeAddLabel(labels, lines, sizeX-1, y, 1, 0, true)
	}

	var start, end XY
	outerToInner := make(map[XY]XY)
	for lbl, arr := range labels {
		switch lbl {
		case "AA":
			start = arr[0]
		case "ZZ":
			end = arr[0]
		default:
			outerToInner[arr[0]] = arr[1]
		}
	}
	return start, outerToInner, end
}

func maybeAddLabel(labels map[string][]XY, lines []string, x int64, y int64, dx int64, dy int64, outer bool) {
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
		arr = make([]XY, 2)
		labels[label] = arr
	}
	if outer {
		arr[0] = XY{x, y}
	} else {
		arr[1] = XY{x, y}
	}
}

type XY struct {
	x, y int64
}

type Cell struct {
	distance int64
}

type Route struct {
	distance int64
	to       XY
	dLevel   int64
}

type PathHeap struct {
	heap []Path
}

func (p *PathHeap) Len() int {
	return len(p.heap)
}

func (p *PathHeap) Less(i, j int) bool {
	return p.heap[i].distanceUntilNow < p.heap[j].distanceUntilNow
}

func (p *PathHeap) Swap(i, j int) {
	tmp := p.heap[i]
	p.heap[i] = p.heap[j]
	p.heap[j] = tmp
}

func (p *PathHeap) Push(e interface{}) {
	p.heap = append(p.heap, e.(Path))
}

func (p *PathHeap) Pop() interface{} {
	heapp := p.heap
	idx := len(heapp) - 1
	p.heap = heapp[:idx]
	return heapp[idx]
}

type Path struct {
	current          XY
	currentLevel     int64
	distanceUntilNow int64
}
