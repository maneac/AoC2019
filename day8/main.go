package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

type layer struct {
	byteMap map[int]int
	width   int
	height  int
	rows    [][]int
}

func main() {
	width := 25
	height := 6
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	read := bufio.NewReader(f)
	isEOF := false

	layers := []layer{layer{map[int]int{}, width, height, make([][]int, height)}}
	layers[0].rows[0] = make([]int, width)
	layerIdx := 0
	rowIdx := 0
	colIdx := 0
	for !isEOF {
		str, err := read.ReadByte()
		if err != nil {
			if err == io.EOF {
				isEOF = true
			}
			panic(err)
		}
		val, err := strconv.Atoi(fmt.Sprintf("%c", str))
		if err != nil {
			if str == '\n' {
				isEOF = true
			} else {
				panic(err)
			}
		}

		if colIdx >= width {
			if rowIdx >= (height - 1) {
				if isEOF {
					break
				}
				layers = append(layers, layer{map[int]int{}, width, height, make([][]int, height)})
				layerIdx++
				rowIdx = 0
			} else {
				rowIdx++
			}
			colIdx = 0
			layers[layerIdx].rows[rowIdx] = make([]int, width)
		}
		layers[layerIdx].byteMap[val]++
		layers[layerIdx].rows[rowIdx][colIdx] = val
		colIdx++
	}
	bestIdx := 0
	maxZero := -1
	for idx, layer := range layers {
		if maxZero < 0 || layer.byteMap[0] < maxZero {
			fmt.Println(layer.byteMap)
			bestIdx = idx
			maxZero = layer.byteMap[0]
		}
	}
	fmt.Println(layers[bestIdx].byteMap[1] * layers[bestIdx].byteMap[2])

	finalImage := make([][]int, height)
	for i := 0; i < len(finalImage); i++ {
		finalImage[i] = make([]int, width)
	}

	for i := len(layers) - 1; i >= 0; i-- {
		for r, row := range layers[i].rows {
			for c, val := range row {
				switch val {
				case 0, 1:
					finalImage[r][c] = val
				}
			}
		}
	}
	for _, row := range finalImage {
		for _, v := range row {
			switch v {
			case 0:
				fmt.Printf(" ")
			case 1:
				fmt.Printf("0")
			case 2:
				fmt.Printf("2")
			}
		}
		fmt.Printf("\n")
	}
}
