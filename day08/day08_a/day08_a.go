package main

import (
	"../../util"
	"fmt"
	"math"
)

const (
	nbRows = 6  //2
	nbCols = 25 //3
)

func main() {
	layers := loadLayers()

	lowestNbOf0s := math.MaxInt32
	score := 0

	for _, layer := range layers {
		nb0s := 0
		nb1s := 0
		nb2s := 0
		for _, row := range layer {
			for _, cell := range row {
				switch cell {
				case 0:
					nb0s++
					break
				case 1:
					nb1s++
					break
				case 2:
					nb2s++
					break
				}
			}
		}
		if nb0s < lowestNbOf0s {
			lowestNbOf0s = nb0s
			score = nb1s * nb2s
		}
	}

	fmt.Println(score)
}

func loadLayers() [][][]byte {
	layers := make([][][]byte, 0)
	layer := -1
	row := nbRows - 1
	col := nbCols - 1

	util.ReadBytes("day08/input.txt", func(data byte) {
		col++
		if col == nbCols {
			row++
			if row == nbRows {
				layer++
				layers = append(layers, make([][]byte, 0, nbRows))
				row = 0
			}
			layers[layer] = append(layers[layer], make([]byte, 0, nbCols))
			col = 0
		}
		layers[layer][row] = append(layers[layer][row], data-'0')
	})
	return layers
}
