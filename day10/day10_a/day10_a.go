package main

import (
	"../../util"
	"fmt"
)

func main() {
	asteroids := loadAsteroids()

	best := asteroids[0]
	best.loadLineOfSight(asteroids)
	for _, candidate := range asteroids {
		candidate.loadLineOfSight(asteroids)
		if len(candidate.inLineOfSight) > len(best.inLineOfSight) {
			best = candidate
		}
	}
	fmt.Println(best.x, best.y, len(best.inLineOfSight))
}

func loadAsteroids() []*asteroid {
	asteroids := make([]*asteroid, 0)
	y := 0

	util.ReadLines("day10/input.txt", func(in string) error {
		for x, c := range in {
			if c == '#' {
				asteroids = append(asteroids, &asteroid{x, y, make([]*asteroid, 0)})
			}
		}
		y++
		return nil
	})
	return asteroids
}

type asteroid struct {
	x             int
	y             int
	inLineOfSight []*asteroid
}

func (a *asteroid) loadLineOfSight(asteroids []*asteroid) {
	if len(a.inLineOfSight) != 0 {
		return
	}

	los := make([]*asteroid, 0)
outer:
	for _, asteroid := range asteroids {
		if asteroid == a {
			continue
		}
		for i, inSight := range los {
			if onOneLine(a, asteroid, inSight) {
				if util.Abs(asteroid.x-a.x)+util.Abs(asteroid.y-a.y) < util.Abs(inSight.x-a.x)+util.Abs(inSight.y-a.y) {
					los[i] = asteroid
				}
				continue outer
			}
		}
		los = append(los, asteroid)
	}
	a.inLineOfSight = los
}

func onOneLine(ref, a1, a2 *asteroid) bool {
	x1 := a1.x - ref.x
	x2 := a2.x - ref.x
	y1 := a1.y - ref.y
	y2 := a2.y - ref.y

	if x1*x2 < 0 || y1*y2 < 0 { // y or x has different sign (relative to ref)
		return false
	}

	return x1*y2 == y1*x2
}
