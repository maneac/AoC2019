package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

type input []int

var relIdx int
var debug bool = false

type dir int

const (
	up dir = iota
	right
	down
	left
)

type point struct {
	x int
	y int
}

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

	ioChans := make([]chan int, 2)
	for i := range ioChans {
		ioChans[i] = make(chan int, 10000)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	grid := make([][]int, 27)
	for i := range grid {
		grid[i] = make([]int, 40)
		if i == (len(grid) - 1) {
			for j := range grid[i] {
				grid[i][j] = 1
			}
		}
	}
	in.runProgram(0, &ioChans, &wg, &grid)

	// close(ioChans[1])

	// blocks := 0

	// for len(ioChans[1]) > 0 {
	// 	x := <-ioChans[1]
	// 	y := <-ioChans[1]
	// 	switch z := <-ioChans[1]; {
	// 	case z == 3:
	// 		paddleLoc = point{x, y}
	// 	case z == 4:
	// 		ballLoc[0] = point{x, y}
	// 	case z == 2:
	// 		blocks++
	// 		fallthrough
	// 	default:
	// 		grid[y][x] = z
	// 	}
	// }

	// printGrid(&grid)

	// fmt.Printf("Tiles: %d\n", blocks)
}

func printGrid(grid *[][]int) {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
<<<<<<< HEAD
=======
	time.Sleep(10 * time.Millisecond)
>>>>>>> 42c847077f28edd9f28bf86c719308cac42e2071
	fmt.Println("")
	for i := range *grid {
		for j := range (*grid)[i] {
			switch (*grid)[i][j] {
			case 0:
				fmt.Print(" ")
			case 1:
				switch i {
				case 0:
					fmt.Print("_")
				case 26:
					fmt.Print("^")
				default:
					fmt.Print("|")
				}
			case 2:
				fmt.Print("#")
			case 3:
				fmt.Print("_")
			case 4:
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
}

func (byteInput *input) runProgram(loc int, ioChans *[]chan int, wg *sync.WaitGroup, grid *[][]int) int {
	defer wg.Done()
	isHalt := false
	idx := 0
	ballLoc := point{-1, -1}
	paddleLoc := point{-1, -1}
	x, y, z := -2, -2, -2
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
			printGrid(grid)
<<<<<<< HEAD
			time.Sleep(time.Second / 30)
=======
			time.Sleep(40 * time.Millisecond)
>>>>>>> 42c847077f28edd9f28bf86c719308cac42e2071
			switch {
			case ballLoc.x > paddleLoc.x:
				(*ioChans)[loc] <- 1
			case ballLoc.x < paddleLoc.x:
				(*ioChans)[loc] <- -1
			default:
				(*ioChans)[loc] <- 0

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
			if x == -2 {
				x = <-(*ioChans)[1]
			} else if y == -2 {
				y = <-(*ioChans)[1]
			} else if z == -2 {
				t := <-(*ioChans)[1]
				if x == -1 && y == 0 {
					fmt.Printf("Score: %d\n", t)
				} else {
					switch {
					case t == 3:
						paddleLoc = point{x, y}
					case t == 4:
						ballLoc = point{x, y}
					}
					(*grid)[y][x] = t
				}
				x, y, z = -2, -2, -2
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
	printGrid(grid)
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
