package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

var debug bool = false

func main() {

	flag.BoolVar(&debug, "debug", false, "Provides debugging output")
	flag.Parse()

	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}

	inData := bufio.NewReader(f)
	byteInput := make([]int, 0)
	isEOF := false
	for !isEOF {
		rawNum, err := inData.ReadBytes(',')
		if err != nil {
			if err == io.EOF {
				isEOF = true
			} else {
				panic(err)
			}
		}
		num, err := strconv.Atoi(string(rawNum[:len(rawNum)-1]))
		if err != nil {
			panic(err)
		}
		byteInput = append(byteInput, num)
	}

	//pt1
	//fmt.Printf("Output: %d\n", runProgram(byteInput))

	//pt2
	rawProg := make([]int, len(byteInput))
	copy(rawProg, byteInput)
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			byteInput[1] = i
			byteInput[2] = j
			if output := runProgram(byteInput); output == 19690720 {
				fmt.Printf("Valid: %v\n", byteInput[1:3])
				return
			}
			copy(byteInput, rawProg)
		}

	}
}

func runProgram(byteInput []int) int {
	idx := 0
	isHalt := false
	for !isHalt {
		switch byteInput[idx] {
		case 1:
			if debug {
				fmt.Printf("Adding index %d (%d) to %d (%d), outputting to %d\n", idx+1, byteInput[idx+1], idx+2, byteInput[idx+2], byteInput[idx+3])
			}
			byteInput[byteInput[idx+3]] = byteInput[byteInput[idx+1]] + byteInput[byteInput[idx+2]]
		case 2:
			if debug {
				fmt.Printf("Multiplying index %d (%d) to %d (%d), outputting to %d\n", idx+1, byteInput[idx+1], idx+2, byteInput[idx+2], byteInput[idx+3])
			}
			byteInput[byteInput[idx+3]] = byteInput[byteInput[idx+1]] * byteInput[byteInput[idx+2]]
		case 99:
			isHalt = true
		default:
			fmt.Printf("Not an opcode: %d (%d)\n", byteInput[idx], idx)
			os.Exit(2)
		}
		idx += 4
	}
	return byteInput[0]
}
