package main

import (
	"../../util"
	"fmt"
	"log"
	"strconv"
)

const (
	EXAMPLE1 = "03036732577212944063491565474664"
	EXAMPLE2 = "02935109699940807407585447034323"
	EXAMPLE3 = "03081770884921959731165446850517"
	INPUT    = "59708372326282850478374632294363143285591907230244898069506559289353324363446827480040836943068215774680673708005813752468017892971245448103168634442773462686566173338029941559688604621181240586891859988614902179556407022792661948523370366667688937217081165148397649462617248164167011250975576380324668693910824497627133242485090976104918375531998433324622853428842410855024093891994449937031688743195134239353469076295752542683739823044981442437538627404276327027998857400463920633633578266795454389967583600019852126383407785643022367809199144154166725123539386550399024919155708875622641704428963905767166129198009532884347151391845112189952083025"
)

type Key struct {
	iteration int64
	index     int64
}

const BASE_VECTOR = INPUT
const FACTOR = 10000
const NB_ITERATIONS int64 = 100

/*
		if v_(n+1) = M * v_n with v a vector and M a matrix

		Then the bottom half of M looks like this

		0 0 0 0 1 1 1 1
	 	0 0 0 0 0 1 1 1
		0 0 0 0 0 0 1 1
	    0 0 0 0 0 0 0 1

		If we are only interested in the very tail of 10000 times the input,
		we can easily abuse this form.

		Note that vector*row_n = vector*row_n+1 + vector[n] (* is vector product)
*/
func main() {
	max := int64(len(BASE_VECTOR)) * FACTOR

	offsetInt, err := strconv.Atoi(BASE_VECTOR[:7])
	if err != nil {
		log.Fatalln("Should not happen", err)
	}
	offset := int64(offsetInt)

	min := offset
	if min < max/2 {
		log.Fatalln("We assume that we are only working at the end")
	}
	size := max - min
	vector := make([]int64, size)
	newVector := make([]int64, size)

	for i := max - 1; i >= min; i-- {
		vector[i-min] = int64(BASE_VECTOR[i%int64(len(BASE_VECTOR))] - '0')
	}

	for i := 0; i < 100; i++ {
		sum := int64(0)
		for j := size - 1; j >= 0; j-- {
			sum += vector[j]
			newVector[j] = util.Abs64(sum % 10)
		}
		vector = newVector
	}

	fmt.Println(vector[:8])
}
