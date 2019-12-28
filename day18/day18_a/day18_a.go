package main

import (
	"../../util"
	"fmt"
	"log"
	"math"
)

func main() {
	//inp := "day18/day18_a/example4.txt"
	inp := "day18/input.txt"
	grid := loadGrid(inp)

	startNode, keys := loadNodes(grid)

	todo := make([]Fn, 0)
	toAdd := startNode.proposePath()
	todo = append(todo, toAdd...)

	for len(todo) != 0 {
		toRun := todo[len(todo)-1]
		todo = todo[:len(todo)-1]

		todo = append(todo, toRun()...)
	}

	allKeys := (int64(1) << len(keys)) - 1

	best := int64(math.MaxInt64)
	for _, key := range keys {
		for _, path := range key.paths {
			if path.provides == allKeys {
				best = util.Min(best, path.distance)
			}
		}
	}

	fmt.Println(best)
}

func loadNodes(grid [][]Cell) (Start, []Key) {
	var startNode Start
	nbKeys := int64(0)
	keys := make([]Key, 256)
	doors := make([]Door, 256)

	for y, cells := range grid {
		for x, cell := range cells {
			switch cell.(type) {
			case Cell_Key:
				key := cell.(Cell_Key).id
				if key >= nbKeys {
					nbKeys = key + 1
				}
				keys[key].init(key)
				flood(grid, x, y, &keys[key], keys, doors)
			case Cell_Door:
				door := cell.(Cell_Door).id
				doors[door].init(door)
				flood(grid, x, y, &doors[door], keys, doors)
			case Cell_Start:
				startNode.init()
				flood(grid, x, y, &startNode, keys, doors)
			}
		}
	}

	return startNode, keys[:nbKeys]
}

func flood(grid [][]Cell, x, y int, point StartingPoint, keys []Key, doors []Door) {
	next := make([]XYD, 0)
	visited := make(map[XY]bool)

	next = append(next, XYD{XY{x, y}, int64(0)})
	for len(next) != 0 {
		cell := next[0]
		next = next[1:]

		_, alreadyVisited := visited[cell.XY]
		if alreadyVisited {
			continue
		}

		visited[cell.XY] = true

		gridCell := grid[cell.y][cell.x]
		switch gridCell.(type) {
		case Cell_Start:
		case Cell_Empty:
		case Cell_Wall:
			continue
		case Cell_Key:
			cellKey := gridCell.(Cell_Key)
			key := &keys[cellKey.id]
			if key != point {
				point.addDirect(key, cell.distance)
				continue
			}
		case Cell_Door:
			cellDoor := gridCell.(Cell_Door)
			door := &doors[cellDoor.id]
			if door != point {
				point.addDirect(door, cell.distance)
				continue
			}
		default:
			log.Fatalln("Should not happen")
		}
		next = append(next, XYD{XY{cell.x + 1, cell.y}, cell.distance + 1})
		next = append(next, XYD{XY{cell.x - 1, cell.y}, cell.distance + 1})
		next = append(next, XYD{XY{cell.x, cell.y + 1}, cell.distance + 1})
		next = append(next, XYD{XY{cell.x, cell.y - 1}, cell.distance + 1})
	}

}

type XYD struct {
	XY
	distance int64
}

type StartingPoint interface {
	addDirect(node Node, dist int64)
}

type Start struct {
	direct map[Node]int64
}

func (s *Start) init() {
	s.direct = make(map[Node]int64)
}

func (s *Start) addDirect(node Node, dist int64) {
	s.direct[node] = dist
}

func (s *Start) proposePath() []Fn {
	fns := make([]Fn, 0)
	for node, dist := range s.direct {
		node2 := node
		dist2 := dist
		fns = append(fns, func() []Fn {
			return node2.proposePath(Path{0, dist2})
		})
	}
	return fns
}

type Fn func() []Fn
type Node interface {
	proposePath(path Path) []Fn
}

type NormalNode struct {
	paths  []Path //path from start, to this
	direct map[Node]int64
}

func (n *NormalNode) init() {
	n.paths = make([]Path, 0)
	n.direct = make(map[Node]int64)
}
func (n *NormalNode) addDirect(node Node, dist int64) {
	n.direct[node] = dist
}
func (n *NormalNode) updatePath(path Path) bool {
	offset := 0
	paths := n.paths
	for i, existingPath := range paths {
		if existingPath.distance <= path.distance &&
			existingPath.provides&path.provides == path.provides {
			// inferior
			return false
		}
		if existingPath.distance >= path.distance &&
			existingPath.provides&path.provides == existingPath.provides {
			// superior -> delete existing one
			paths[i-offset] = paths[len(paths)-1]
			offset++
			paths = paths[:len(paths)-1]
		}
	}
	n.paths = append(paths, path)
	return true
}

type Key struct {
	NormalNode
	id   int64
	leaf bool
}

func (k *Key) init(id int64) {
	k.NormalNode.init()
	k.id = id
	k.leaf = true
}

func (k *Key) proposePath(path Path) []Fn {
	fns := make([]Fn, 0)

	path.provides = path.provides | (1 << k.id)
	wasAdded := k.updatePath(path)
	if !wasAdded {
		return fns
	}

	for node, dist := range k.direct {
		node2 := node
		dist2 := dist
		provides := path.provides
		fns = append(fns, func() []Fn {
			return node2.proposePath(Path{provides, path.distance + dist2})
		})
	}
	return fns
}

type Door struct {
	NormalNode
	id int64
}

func (d *Door) init(id int64) {
	d.NormalNode.init()
	d.id = id
}

func (d *Door) proposePath(path Path) []Fn {
	wasAdded := d.updatePath(path)
	fns := make([]Fn, 0)

	if !wasAdded {
		return fns
	}

	mask := int64(1) << d.id
	if path.provides&mask != mask {
		return fns
	}

	for node, dist := range d.direct {
		node2 := node
		dist2 := dist
		fns = append(fns, func() []Fn {
			return node2.proposePath(Path{path.provides, path.distance + dist2})
		})
	}
	return fns
}

type Path struct {
	provides int64
	distance int64
}

//Cells
type Cell interface{}
type Cell_Empty struct{}
type Cell_Wall struct{}
type Cell_Key struct {
	id int64
}
type Cell_Door struct {
	id int64
}
type Cell_Start struct{}

func loadGrid(file string) [][]Cell {
	grid := make([][]Cell, 0)

	y := 0
	util.ReadLines(file, func(in string) error {
		row := make([]Cell, len(in))
		for x := 0; x < len(in); x++ {
			char := in[x]
			var cell Cell
			switch {
			case 'a' <= char && char <= 'z':
				cell = Cell_Key{int64(char - 'a')}
			case 'A' <= char && char <= 'Z':
				cell = Cell_Door{int64(char - 'A')}
			case char == '.':
				cell = Cell_Empty{}
			case char == '#':
				cell = Cell_Wall{}
			case char == '@':
				cell = Cell_Start{}
			default:
				log.Fatalln("Unknown char")
			}
			row[x] = cell
		}
		grid = append(grid, row)
		y++
		return nil
	})

	return grid
}

type XY struct {
	x, y int
}
