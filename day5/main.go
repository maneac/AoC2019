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
	runProgram(byteInput)
}

func runProgram(byteInput []int) int {
	idx := 0
	isHalt := false
	scanner := bufio.NewScanner(os.Stdin)
	for !isHalt {
		opcode := fmt.Sprintf("%05d", byteInput[idx])
		switch opcode[len(opcode)-2:] {
		case "01":
			in1, in2 := getIn(byteInput, opcode, idx)
			if debug {
				fmt.Printf("Adding index %d (%d) to %d (%d), outputting to %d\n", idx+1, in1, idx+2, in2, byteInput[idx+3])
			}
			byteInput[byteInput[idx+3]] = in1 + in2
			idx += 4
		case "02":
			in1, in2 := getIn(byteInput, opcode, idx)
			if debug {
				fmt.Printf("Multiplying index %d (%d) to %d (%d), outputting to %d\n", idx+1, in1, idx+2, in2, byteInput[idx+3])
			}
			byteInput[byteInput[idx+3]] = in1 * in2
			idx += 4
		case "03":
			if debug {
				fmt.Printf("Reading input to index %d\n", byteInput[idx+1])
			}
			fmt.Printf("Waiting for input...\n")
			scanner.Scan()
			n, err := strconv.Atoi(scanner.Text())
			if err != nil {
				panic(err)
			}
			byteInput[byteInput[idx+1]] = n
			idx += 2
		case "04":
			if debug {
				fmt.Printf("Outputting value at index %d, mode %d\n", byteInput[idx+1], opcode[2])
			}
			switch opcode[2] {
			case '0':
				fmt.Printf("Ouptut: %d\n", byteInput[byteInput[idx+1]])
			case '1':
				fmt.Printf("Ouptut: %d\n", byteInput[idx+1])
			}
			idx += 2
		case "05":
			in1, in2 := getIn(byteInput, opcode, idx)
			if in1 != 0 {
				idx = in2
			} else {
				idx += 3
			}
		case "06":
			in1, in2 := getIn(byteInput, opcode, idx)
			if in1 == 0 {
				idx = in2
			} else {
				idx += 3
			}
		case "07":
			in1, in2 := getIn(byteInput, opcode, idx)
			if in1 < in2 {
				byteInput[byteInput[idx+3]] = 1
			} else {
				byteInput[byteInput[idx+3]] = 0
			}
			idx += 4
		case "08":
			in1, in2 := getIn(byteInput, opcode, idx)
			if in1 == in2 {
				byteInput[byteInput[idx+3]] = 1
			} else {
				byteInput[byteInput[idx+3]] = 0
			}
			idx += 4
		case "99":
			isHalt = true
		default:
			fmt.Printf("Not an opcode: %d (%d)\n", byteInput[idx], byteInput[idx-1:idx+2])
			fmt.Println(byteInput)
			os.Exit(2)
		}
	}
	return byteInput[0]
}

func getIn(byteInput []int, opcode string, idx int) (int, int) {
	var in1, in2 int
	switch opcode[2] {
	case '0':
		in1 = byteInput[byteInput[idx+1]]
	case '1':
		in1 = byteInput[idx+1]
	}
	switch opcode[1] {
	case '0':
		in2 = byteInput[byteInput[idx+2]]
	case '1':
		in2 = byteInput[idx+2]
	}
	return in1, in2
}
