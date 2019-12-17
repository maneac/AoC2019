package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
)

type input []int

var relIdx int
var debug bool = false

func main() {

	flag.BoolVar(&debug, "debug", false, "Provides debugging output")
	flag.Parse()

	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}

	inData := bufio.NewReader(f)

	isEOF := false
	in := input{}

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
		in = append(in, num)
	}

	byteInputs := make([]input, 5)
	for i := 0; i < 5; i++ {
		byteInputs[i] = make(input, len(in))
		copy(byteInputs[i], in)
	}

	ioChans := make([]chan int, 2)
	for i := range ioChans {
		ioChans[i] = make(chan int, 100)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go in.runProgram(0, &ioChans, &wg)
	ioChans[0] <- 2
	wg.Wait()
	for {
		if len(ioChans[1]) == 0 {
			break
		}
		fmt.Println(<-ioChans[1])
	}
}

func (byteInput *input) runProgram(loc int, ioChans *[]chan int, wg *sync.WaitGroup) int {
	defer wg.Done()
	isHalt := false
	idx := 0
	for !isHalt {
		opcode := fmt.Sprintf("%05d", (*byteInput)[idx])
		switch opcode[len(opcode)-2:] {
		case "01":
			in1, in2, outIdx := byteInput.getIn(opcode, idx)
			if debug {
				fmt.Printf("Adding index %d (%d) to %d (%d), outputting to %d\n", idx+1, in1, idx+2, in2, (*byteInput)[idx+3])
			}
			(*byteInput)[outIdx] = in1 + in2
			idx += 4
		case "02":
			in1, in2, outIdx := byteInput.getIn(opcode, idx)
			if debug {
				fmt.Printf("Multiplying index %d (%d) to %d (%d), outputting to %d\n", idx+1, in1, idx+2, in2, (*byteInput)[idx+3])
			}
			(*byteInput)[outIdx] = in1 * in2
			idx += 4
		case "03":
			if debug {
				fmt.Printf("Waiting for input at %d\n", loc)
			}
			switch opcode[2] {
			case '0':
				(*byteInput)[(*byteInput)[idx+1]] = <-(*ioChans)[loc]
			case '2':
				(*byteInput)[relIdx+(*byteInput)[idx+1]] = <-(*ioChans)[loc]
			}
			if debug {
				fmt.Printf("Received input at %d\n", loc)
			}
			idx += 2
		case "04":
			if debug {
				fmt.Printf("Outputting value at index %d, mode %d\n", (*byteInput)[idx+1], opcode[2])
			}
			switch opcode[2] {
			case '0':
				if loc == 4 {
					(*ioChans)[0] <- (*byteInput)[(*byteInput)[idx+1]]
				} else {
					(*ioChans)[loc+1] <- (*byteInput)[(*byteInput)[idx+1]]
				}
			case '1':
				if loc == 4 {
					(*ioChans)[0] <- (*byteInput)[idx+1]
				} else {
					(*ioChans)[loc+1] <- (*byteInput)[idx+1]
				}
			case '2':
				if loc == 4 {
					(*ioChans)[0] <- (*byteInput)[relIdx+(*byteInput)[idx+1]]
				} else {
					(*ioChans)[loc+1] <- (*byteInput)[relIdx+(*byteInput)[idx+1]]
				}
			}
			if debug {
				fmt.Printf("Outputting to %d\n", loc+1)
			}
			idx += 2
		case "05":
			in1, in2, _ := byteInput.getIn(opcode, idx)
			if in1 != 0 {
				idx = in2
			} else {
				idx += 3
			}
		case "06":
			in1, in2, _ := byteInput.getIn(opcode, idx)
			if in1 == 0 {
				idx = in2
			} else {
				idx += 3
			}
		case "07":
			in1, in2, outIdx := byteInput.getIn(opcode, idx)
			if in1 < in2 {
				(*byteInput)[outIdx] = 1
			} else {
				(*byteInput)[outIdx] = 0
			}
			idx += 4
		case "08":
			in1, in2, outIdx := byteInput.getIn(opcode, idx)
			if in1 == in2 {
				(*byteInput)[outIdx] = 1
			} else {
				(*byteInput)[outIdx] = 0
			}
			idx += 4
		case "09":
			in1, _, _ := byteInput.getIn(opcode, idx)
			relIdx += in1
			idx += 2
		case "99":
			if debug {
				fmt.Printf("Halting %d\n", loc)
			}
			isHalt = true
		default:
			if idx == 0 {
				fmt.Printf("Not an opcode (%d): %q (%d)\n", loc, opcode, (*byteInput)[:idx+3])
			} else {
				fmt.Printf("Not an opcode (%d): %q (%d)\n", loc, opcode, (*byteInput)[idx-1:idx+2])
			}
			fmt.Println((*byteInput))
			os.Exit(2)
		}
	}
	return (*byteInput)[0]
}

//Returns input1, input2, and the index to write output to
func (byteInput *input) getIn(opcode string, idx int) (int, int, int) {
	var in1, in2, outIdx int
	switch opcode[2] {
	case '0':
		if idx+1 >= len((*byteInput)) {
			(*byteInput) = append((*byteInput), make([]int, (idx+1)-len((*byteInput))+1)...)
		}
		if (*byteInput)[idx+1] >= len((*byteInput)) {
			(*byteInput) = append((*byteInput), make([]int, (*byteInput)[idx+1]-len((*byteInput))+1)...)
		}
		in1 = (*byteInput)[(*byteInput)[idx+1]]
	case '1':
		if idx+1 >= len((*byteInput)) {
			(*byteInput) = append((*byteInput), make([]int, (idx+1)-len((*byteInput))+1)...)
		}
		in1 = (*byteInput)[idx+1]
	case '2':
		if idx+1 >= len((*byteInput)) {
			(*byteInput) = append((*byteInput), make([]int, (idx+1)-len((*byteInput))+1)...)
		}
		if relIdx+(*byteInput)[idx+1] >= len((*byteInput)) {
			(*byteInput) = append((*byteInput), make([]int, relIdx+(*byteInput)[idx+1]-len((*byteInput))+1)...)
		}
		in1 = (*byteInput)[relIdx+(*byteInput)[idx+1]]
	}

	switch opcode[1] {
	case '0':
		if idx+2 >= len((*byteInput)) {
			(*byteInput) = append((*byteInput), make([]int, (idx+2)-len((*byteInput))+1)...)
		}
		if (*byteInput)[idx+2] >= len((*byteInput)) {
			(*byteInput) = append((*byteInput), make([]int, (*byteInput)[idx+2]-len((*byteInput))+1)...)
		}
		in2 = (*byteInput)[(*byteInput)[idx+2]]
	case '1':
		if idx+2 >= len((*byteInput)) {
			(*byteInput) = append((*byteInput), make([]int, (idx+2)-len((*byteInput))+1)...)
		}
		in2 = (*byteInput)[idx+2]
	case '2':
		if idx+2 >= len((*byteInput)) {
			(*byteInput) = append((*byteInput), make([]int, (idx+2)-len((*byteInput))+1)...)
		}
		if relIdx+(*byteInput)[idx+2] >= len((*byteInput)) {
			(*byteInput) = append((*byteInput), make([]int, relIdx+(*byteInput)[idx+2]-len((*byteInput))+1)...)
		}
		in2 = (*byteInput)[relIdx+(*byteInput)[idx+2]]
	}
	switch opcode[0] {
	case '0':
		if idx+3 >= len((*byteInput)) {
			(*byteInput) = append((*byteInput), make([]int, (idx+3)-len((*byteInput))+1)...)
		}
		if (*byteInput)[idx+3] >= len((*byteInput)) {
			(*byteInput) = append((*byteInput), make([]int, (*byteInput)[idx+3]-len((*byteInput))+1)...)
		}
		outIdx = (*byteInput)[idx+3]
	case '1':
		if idx+3 >= len((*byteInput)) {
			(*byteInput) = append((*byteInput), make([]int, (idx+3)-len((*byteInput))+1)...)
		}
		outIdx = idx + 3
	case '2':
		if idx+3 >= len((*byteInput)) {
			(*byteInput) = append((*byteInput), make([]int, (idx+3)-len((*byteInput))+1)...)
		}
		if relIdx+(*byteInput)[idx+3] >= len((*byteInput)) {
			(*byteInput) = append((*byteInput), make([]int, relIdx+(*byteInput)[idx+3]-len((*byteInput))+1)...)
		}
		outIdx = relIdx + (*byteInput)[idx+3]
	}
	return in1, in2, outIdx
}
