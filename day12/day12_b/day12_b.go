package main

import (
	"../../util"
	"fmt"
	"log"
)

func main() {
	systems := input()
	var cycles [3]int64

	for systemIdx, system := range systems {
		var skip int64
		skip, cycles[systemIdx] = runSystem(system)
		if skip != 0 {
			log.Println("We assumed skip = 0", systemIdx)
		}
		fmt.Println(systemIdx, skip, cycles[systemIdx])
	}

	res := util.Lcm(util.Lcm(cycles[0], cycles[1]), cycles[2])
	fmt.Println(res)
}

func example() [][]*miniMoon {
	return [][]*miniMoon{
		{
			{-1, 0},
			{2, 0},
			{4, 0},
			{3, 0},
		},
		{
			{0, 0},
			{-10, 0},
			{-8, 0},
			{5, 0},
		},
		{
			{2, 0},
			{-7, 0},
			{8, 0},
			{-1, 0},
		},
	}
}

func input() [][]*miniMoon {
	return [][]*miniMoon{
		{
			{-3, 0},
			{3, 0},
			{-13, 0},
			{6, 0},
		},
		{
			{15, 0},
			{13, 0},
			{18, 0},
			{0, 0},
		},
		{
			{-11, 0},
			{-19, 0},
			{-2, 0},
			{-1, 0},
		},
	}
}

type moon4 struct {
	moon0, moon1, moon2, moon3 miniMoon
}

func runSystem(moons []*miniMoon) (int64, int64) {
	knownPositions := make(map[moon4]int64)
	i := int64(0)
	for {
		knownPositions[moon4{
			*moons[0],
			*moons[1],
			*moons[2],
			*moons[3],
		}] = i
		step(moons)
		i++
		j, found := knownPositions[moon4{
			*moons[0],
			*moons[1],
			*moons[2],
			*moons[3],
		}]
		if found {
			return j, i
		}
	}
}

func step(moons []*miniMoon) {
	for _, toUpdate := range moons {
		for _, pulling := range moons {
			if toUpdate == pulling {
				continue
			}
			toUpdate.applyGravity(pulling)
		}
	}
	for _, toUpdate := range moons {
		toUpdate.applyVelocity()
	}
}

type xyz struct {
	x, y, z int64
}

func (x *xyz) plusEq(other *xyz) {
	x.x += other.x
	x.y += other.y
	x.z += other.z
}

type miniMoon struct {
	pos, vel int64
}

func (m *miniMoon) applyGravity(other *miniMoon) {
	m.vel += gravity(m.pos, other.pos)
}

func (m *miniMoon) applyVelocity() {
	m.pos += m.vel
}

func (m *miniMoon) copy() *miniMoon {
	return &miniMoon{m.pos, m.vel}
}

func gravity(affected int64, causedBy int64) int64 {
	switch {
	case causedBy == affected:
		return 0
	case causedBy > affected:
		return 1
	case causedBy < affected:
		return -1
	}
	log.Fatal("Should not happen")
	return 0
}
