package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	f, err := os.Open("./input.csv")
	if err != nil {
		panic(err)
	}

	inputReader := csv.NewReader(f)

	total := 0

	for {
		rec, err := inputReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		n, err := strconv.Atoi(rec[0])
		if err != nil {
			panic(err)
		}

		for {
			subtotal := (n / 3) - 2
			if subtotal <= 0 {
				break
			}
			total += subtotal
			n = subtotal
			//fmt.Printf("n: %d, subtotal: %d, total: %d\n", n, subtotal, total)
		}
		//break

	}
	fmt.Printf("Total: %d\n", total)
}
