package main

import (
	"../../util"
	"fmt"
	"log"
)

const (
	nbRows = 6  //2
	nbCols = 25 //3
)

func main() {
	layers := loadLayers()

	for row := 0; row < nbRows; row++ {
	columnLoop:
		for col := 0; col < nbCols; col++ {
			for layer := 0; layer < len(layers); layer++ {
				switch layers[layer][row][col] {
				case 0:
					fmt.Print(" ")
					continue columnLoop
				case 1:
					fmt.Print("#")
					continue columnLoop
				}
			}
			log.Fatal("Nothing found")
		}
		fmt.Println()
	}

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
