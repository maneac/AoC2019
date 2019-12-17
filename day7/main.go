package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strconv"
	"sync"
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

	isEOF := false
	in := []int{}

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

	byteInputs := make([][]int, 5)
	for i := 0; i < 5; i++ {
		byteInputs[i] = make([]int, len(in))
		copy(byteInputs[i], in)
	}

	phases := [][]int{}
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			for k := 0; k < 5; k++ {
				for l := 0; l < 5; l++ {
					for m := 0; m < 5; m++ {
						tmpPhase := []int{i + 5, j + 5, k + 5, l + 5, m + 5}
						sort.Ints(tmpPhase)
						if reflect.DeepEqual(tmpPhase, []int{5, 6, 7, 8, 9}) {
							phases = append(phases, []int{i + 5, j + 5, k + 5, l + 5, m + 5})
						}
					}
				}
			}
		}
	}
	best := 0
	bestIdx := 0
	for idx, phase := range phases {
		ioChans := make([]chan int, 5)
		wg := sync.WaitGroup{}
		for i := 0; i < 5; i++ {
			ioChans[i] = make(chan int, 2)
			ioChans[i] <- phase[i]
		}
		for i := 0; i < 5; i++ {
			go runProgram(i, byteInputs[i], &ioChans, &wg)
		}
		wg.Add(5)
		ioChans[0] <- 0
		wg.Wait()
		res := <-ioChans[0]
		if debug {
			fmt.Printf("Outputs\n%d\n\n", res)
		}
		if res > best {
			bestIdx = idx
			best = res
		}
		for i := 0; i < 5; i++ {
			copy(byteInputs[i], in)
		}
	}
	fmt.Printf("Best: %d %v %d\n", bestIdx, phases[bestIdx], best)
}

func runProgram(loc int, byteInput []int, ioChans *[]chan int, wg *sync.WaitGroup) int {
	defer wg.Done()
	isHalt := false
	idx := 0
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
				fmt.Printf("Waiting for input at %d\n", loc)
			}
			byteInput[byteInput[idx+1]] = <-(*ioChans)[loc]
			if debug {
				fmt.Printf("Received input at %d\n", loc)
			}
			idx += 2
		case "04":
			if debug {
				fmt.Printf("Outputting value at index %d, mode %d\n", byteInput[idx+1], opcode[2])
			}
			switch opcode[2] {
			case '0':
				if loc == 4 {
					(*ioChans)[0] <- byteInput[byteInput[idx+1]]
					if debug {
						fmt.Printf("Outputting to 0\n")
					}
				} else {
					(*ioChans)[loc+1] <- byteInput[byteInput[idx+1]]
					if debug {
						fmt.Printf("Outputting to %d\n", loc+1)
					}
				}
			case '1':
				if loc == 4 {
					(*ioChans)[0] <- byteInput[idx+1]
				} else {
					(*ioChans)[loc+1] <- byteInput[idx+1]
				}
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
			if debug {
				fmt.Printf("Halting %d\n", loc)
			}
			isHalt = true
		default:
			if idx == 0 {
				fmt.Printf("Not an opcode (%d): %q (%d)\n", loc, opcode, byteInput[:idx+3])
			} else {
				fmt.Printf("Not an opcode (%d): %q (%d)\n", loc, opcode, byteInput[idx-1:idx+2])
			}
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
