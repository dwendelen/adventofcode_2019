package main

import (
	"fmt"
	"strconv"
)

func main() {
	from := 382345
	to := 843167

	nb := 0
	for p := from; p <= to; p++ {
		if meetsRules(p) {
			nb++
		}
	}

	fmt.Println(nb)
}

func meetsRules(p int) bool {
	str := strconv.Itoa(p)

	upRule := str[0] <= str[1] &&
		str[1] <= str[2] &&
		str[2] <= str[3] &&
		str[3] <= str[4] &&
		str[4] <= str[5]

	doubleRule := str[0] == str[1] ||
		str[1] == str[2] ||
		str[2] == str[3] ||
		str[3] == str[4] ||
		str[4] == str[5]

	return upRule && doubleRule
}
