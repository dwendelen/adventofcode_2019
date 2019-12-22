package main

import (
	"../../util"
	"fmt"
	"log"
	"regexp"
	"strconv"
)

var regex = regexp.MustCompile("(deal into new stack)|cut (-?\\d*)|deal with increment (\\d*)")

const file = "day22/input.txt"
const size = 10007

func main() {

	cards := make([]int64, size)
	other := make([]int64, size)

	for i := int64(0); i < size; i++ {
		cards[i] = i
	}

	util.ReadLines(file, func(in string) error {
		pieces := regex.FindStringSubmatch(in)
		switch {
		case pieces[1] != "":
			newStack(cards, other)
		case pieces[2] != "":
			asInt, err := strconv.Atoi(pieces[2])
			if err != nil {
				log.Fatalln("Should not happen")
			}
			cut(cards, asInt, other)
		case pieces[3] != "":
			asInt, err := strconv.Atoi(pieces[3])
			if err != nil {
				log.Fatalln("Should not happen")
			}
			increments(cards, asInt, other)
		default:
			log.Fatalln("Should not happen")
		}
		cards, other = other, cards

		for i := 0; i < size; i++ {
			if cards[i] == 2019 {
				fmt.Println(i)
				break
			}
		}

		return nil
	})

	for i := 0; i < size; i++ {
		if cards[i] == 2019 {
			fmt.Println(i)
			break
		}
	}
}

func newStack(cards []int64, other []int64) {
	for i, card := range cards {
		other[size-1-i] = card
	}
}

func cut(cards []int64, cutPoint int, other []int64) {
	if cutPoint < 0 {
		cutPoint += size
	}

	j := 0
	for i := cutPoint; i < size; i++ {
		other[j] = cards[i]
		j++
	}
	for i := 0; i < cutPoint; i++ {
		other[j] = cards[i]
		j++
	}
}

func increments(cards []int64, incr int, other []int64) {
	for i := 0; i < size; i++ {
		other[(incr*i)%size] = cards[i]
	}
}
