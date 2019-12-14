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
	ore.amountLeft = 1000000000000
	ore.commit()

	fuel := reactions["FUEL"]

	res := int64(0)
	for batch := int64(1) << 20; batch >= 1; batch = batch >> 4 {
		fmt.Println("Batch", batch)
		for {
			success := fuel.takeAmount(batch)
			if success {
				res += batch
				for _, reaction := range reactions {
					reaction.commit()
				}
			} else {
				for _, reaction := range reactions {
					reaction.rollback()
				}
				break
			}
		}
	}
	fmt.Println(res)
}

type Reaction struct {
	name                string
	quantity            int64
	sources             map[*Reaction]int64
	amountLeft          int64
	amountLeftCommitted int64
}

func (p *Reaction) addSource(quantity int64, source *Reaction) {
	p.sources[source] = quantity
}

func (p *Reaction) commit() {
	p.amountLeftCommitted = p.amountLeft
}

func (p *Reaction) rollback() {
	p.amountLeft = p.amountLeftCommitted
}

func (p *Reaction) takeAmount(amount int64) bool {
	amountToProduce := util.Max(0, amount-p.amountLeft)

	if p.name == "ORE" && amountToProduce > 0 {
		return false
	}

	nbOfReactions := util.Ceil64(amountToProduce, p.quantity)

	for source, amount := range p.sources {
		if !source.takeAmount(amount * nbOfReactions) {
			return false
		}
	}

	amountProduced := p.quantity * nbOfReactions

	p.amountLeft = p.amountLeft - amount + amountProduced
	if p.amountLeft < 0 {
		log.Fatalln("Should not happen")
	}

	return true
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
		prod = &Reaction{name, 0, make(map[*Reaction]int64), 0, 0}
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
