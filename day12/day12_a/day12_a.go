package main

import (
	"../../util"
	"fmt"
	"log"
)

func main() {
	moons := input()
	for i := 0; i < 1000; i++ {
		step(moons)
	}
	energy := totalEnergy(moons)
	fmt.Println(energy)
}

func example() []*moon {
	return []*moon{
		{xyz{-1, 0, 2}, xyz{0, 0, 0}},
		{xyz{2, -10, -7}, xyz{0, 0, 0}},
		{xyz{4, -8, 8}, xyz{0, 0, 0}},
		{xyz{3, 5, -1}, xyz{0, 0, 0}},
	}
}

func input() []*moon {
	return []*moon{
		{xyz{-3, 15, -11}, xyz{0, 0, 0}},
		{xyz{3, 13, -19}, xyz{0, 0, 0}},
		{xyz{-13, 18, -2}, xyz{0, 0, 0}},
		{xyz{6, 0, -1}, xyz{0, 0, 0}},
	}
}

func step(moons []*moon) {
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

func totalEnergy(moons []*moon) int64 {
	sum := int64(0)
	for _, moon := range moons {
		sum += moon.energy()
	}
	return sum
}

type xyz struct {
	x, y, z int64
}

func (x *xyz) plusEq(other *xyz) {
	x.x += other.x
	x.y += other.y
	x.z += other.z
}

type moon struct {
	pos, vel xyz
}

func (m *moon) applyGravity(other *moon) {
	m.vel.x += gravity(m.pos.x, other.pos.x)
	m.vel.y += gravity(m.pos.y, other.pos.y)
	m.vel.z += gravity(m.pos.z, other.pos.z)
}

func (m *moon) applyVelocity() {
	m.pos.plusEq(&m.vel)
}

func (m *moon) energy() int64 {
	kin := util.Abs64(m.vel.x) +
		util.Abs64(m.vel.y) +
		util.Abs64(m.vel.z)
	pot := util.Abs64(m.pos.x) +
		util.Abs64(m.pos.y) +
		util.Abs64(m.pos.z)

	return kin * pot
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
