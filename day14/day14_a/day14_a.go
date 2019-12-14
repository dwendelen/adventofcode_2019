package main

import (
	"../../util"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	reactions := loadReactions("day14/input.txt")
	ore := reactions["ORE"]
	ore.quantity = 1
	res := ore.amountOfReactionsNeededForOneFuel()
	fmt.Println(res)
}

type Reaction struct {
	name                  string
	quantity              int64
	isUsedBy              map[*Reaction]int64
	amountReactionsNeeded int64
}

func (p *Reaction) addSource(quantity int64, source *Reaction) {
	source.isUsedBy[p] = quantity
}

func (p *Reaction) amountOfReactionsNeededForOneFuel() int64 {
	if p.name == "FUEL" {
		return 1
	}

	if p.amountReactionsNeeded < 0 {
		total := int64(0)
		for usedBy, amount := range p.isUsedBy {
			total += amount * usedBy.amountOfReactionsNeededForOneFuel()
		}

		p.amountReactionsNeeded = util.Ceil64(total, p.quantity)
	}

	return p.amountReactionsNeeded
}

var regex1 = regexp.MustCompile("(\\d* \\w*(, \\d* \\w*)*) => (\\d* \\w*)")

func loadReactions(file string) map[string]*Reaction {
	reactions := make(map[string]*Reaction)
	util.ReadLines(file, func(in string) error {

		submatch := regex1.FindStringSubmatch(in)
		if submatch == nil {
			log.Fatalln("Not matching regex1", in)
		}

		nb, name := parsePiece(submatch[3])
		prod := getOrCreate(reactions, name)
		prod.quantity = nb
		prod.name = name

		pieces := strings.Split(submatch[1], ", ")

		for _, piece := range pieces {
			nb, name := parsePiece(piece)
			source := getOrCreate(reactions, name)
			prod.addSource(nb, source)
		}

		return nil
	})
	return reactions
}

func getOrCreate(reactions map[string]*Reaction, name string) *Reaction {
	prod, found := reactions[name]
	if !found {
		prod = &Reaction{name, 0, make(map[*Reaction]int64), -1}
		reactions[name] = prod
	}

	return prod
}

var regex2 = regexp.MustCompile("(\\d*) (\\w*)")

func parsePiece(piece string) (int64, string) {
	submatch := regex2.FindStringSubmatch(piece)
	if submatch == nil {
		log.Fatalln("Not matching regex2", piece)
	}
	q, err := strconv.Atoi(submatch[1])
	if err != nil {
		log.Fatalln("Invalid number", submatch[1], err)
	}

	return int64(q), submatch[2]
}
