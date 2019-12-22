package main

import (
	"../../util"
	"fmt"
	"log"
	"math/big"
	"regexp"
	"strconv"
)

var regex = regexp.MustCompile("(deal into new stack)|cut (-?\\d*)|deal with increment (\\d*)")

const file = "day22/input.txt"
const size = int64(119315717514047)
const timesShuffled = int64(101741582076661)
const index = int64(2020)

var bigSize = big.NewInt(size)

func main() {
	allLines := util.ReadAllLines(file)

	//oldIndex = (factor * index + sum) mod bigSize
	factor := big.NewInt(1)
	sum := big.NewInt(0)

	for i := len(allLines) - 1; i >= 0; i-- {
		pieces := regex.FindStringSubmatch(allLines[i])
		switch {
		case pieces[1] != "":
			//j = -i - 1
			//i = -j - 1
			factor = mod(neg(factor), bigSize)
			sum = mod(add(neg(sum), big.NewInt(-1)), bigSize)
		case pieces[2] != "":
			//j = i - cut
			//i = j + cut
			asInt, err := strconv.Atoi(pieces[2])
			if err != nil {
				log.Fatalln("Should not happen")
			}

			sum = mod(add(sum, big.NewInt(int64(asInt))), bigSize)
		case pieces[3] != "":
			//j = incr * i
			//i = j / incr
			asInt, err := strconv.Atoi(pieces[3])
			if err != nil {
				log.Fatalln("Should not happen")
			}
			inverseMod := inverseMultipleMod(int64(asInt), size)
			bigInverse := big.NewInt(inverseMod)

			factor = mod(mul(bigInverse, factor), bigSize)
			sum = mod(mul(bigInverse, sum), bigSize)
		default:
			log.Fatalln("Should not happen")
		}
	}

	ff, ss := applyNTimes(factor, sum, big.NewInt(1), big.NewInt(0), timesShuffled)
	res := mod(add(mul(ff, big.NewInt(index)), ss), bigSize)

	fmt.Println(res)
}

/*
	How to compose h(index) = g(f(index))
	-> 	factor_h = factor_g * factor_f
		sum_h 	 = factor_g * sum_f + sum_g

	So g(index) = f^2(index) = f(f(index))
	-> 	factor_g = factor_f * factor_f
		sum_g    = factor_f * sum_f + sum_f

	Then we can use the same principle to calculate f^n where n is a power of 2

	With our powers of two, we can easily use this to calculate f^n for any arbitrary n
*/
func applyNTimes(factor, sum, accFactor, accSum *big.Int, n int64) (*big.Int, *big.Int) {
	if n == 0 {
		return accFactor, accSum
	}

	if n%2 == 1 {
		accFactor = mod(mul(accFactor, factor), bigSize)
		accSum = mod(add(mul(accSum, factor), sum), bigSize)
	}

	factor2 := mod(mul(factor, factor), bigSize)
	sum2 := mod(add(mul(factor, sum), sum), bigSize)

	return applyNTimes(factor2, sum2, accFactor, accSum, n/2)
}

func inverseMultipleMod(a, b int64) int64 {
	s := int64(0)
	old_s := int64(1)
	r := b
	old_r := a

	for r != 0 {
		quotient := old_r / r
		old_r, r = r, old_r-quotient*r
		old_s, s = s, old_s-quotient*s
	}

	return old_s
}

func neg(a *big.Int) *big.Int {
	return big.NewInt(0).Neg(a)
}

func add(a, b *big.Int) *big.Int {
	return big.NewInt(0).Add(a, b)
}

func mul(a, b *big.Int) *big.Int {
	return big.NewInt(0).Mul(a, b)
}

func mod(a, m *big.Int) *big.Int {
	tmp := big.NewInt(0)
	tmp.Mod(a, m)

	if tmp.Cmp(big.NewInt(0)) < 0 {
		tmp.Add(tmp, m)
	}
	return tmp
}
