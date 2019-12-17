package main

import (
	"fmt"
	"strconv"
)

/*
	inputMin := 264793
	inputMax := 803935
*/
func main() {
	inputMin := 264793
	inputMax := 803935
	validOutputs := 0
	outputs := make([]int, 0)
	for i := inputMin; i <= inputMax; i++ {
		num := strconv.Itoa(i)
		isValid := true
		hasDouble := false
		for j := 0; j < 5; j++ {
			if num[j+1] < num[j] {
				isValid = false
				break
			}
			if num[j+1] == num[j] {
				switch j {
				case 0:
					if num[j+2] != num[j+1] {
						hasDouble = true
					}
				case 1, 2, 3:
					if num[j+2] != num[j+1] && num[j-1] != num[j] {
						hasDouble = true
					}
				case 4:
					if num[j-1] != num[j] {
						hasDouble = true
					}
				}
			}
		}
		if !isValid || !hasDouble {
			continue
		}
		outputs = append(outputs, i)
		validOutputs++
	}
	//fmt.Println(outputs)
	fmt.Printf("Valid outputs: %d\n", validOutputs)
}
