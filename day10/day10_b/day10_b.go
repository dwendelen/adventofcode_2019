package main

import (
	"../../util"
	"fmt"
	"log"
	"math"
	"sort"
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

	if len(best.inLineOfSight) < 200 {
		log.Fatal("We need at least 200 asteroids within LOS")
	}

	sort.Sort(best)

	nb200 := best.inLineOfSight[199]

	fmt.Println(nb200.x, nb200.y)
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

func (a *asteroid) Len() int {
	return len(a.inLineOfSight)
}

func (a *asteroid) Less(i, j int) bool {
	a1 := a.inLineOfSight[i]
	a2 := a.inLineOfSight[j]

	x1 := a1.x - a.x
	y1 := a1.y - a.y
	x2 := a2.x - a.x
	y2 := a2.y - a.y

	angle1 := math.Atan2(float64(x1), float64(-y1))
	if angle1 < 0 {
		angle1 += 2 * math.Pi
	}
	angle2 := math.Atan2(float64(x2), float64(-y2))
	if angle2 < 0 {
		angle2 += 2 * math.Pi
	}

	return angle1 < angle2
}

func (a *asteroid) Swap(i, j int) {
	tmp := a.inLineOfSight[i]
	a.inLineOfSight[i] = a.inLineOfSight[j]
	a.inLineOfSight[j] = tmp
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
