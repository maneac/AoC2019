package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

type material struct {
	name string
	amt  int
}

type reaction struct {
	materials []material
	amt       int
}

type reactions map[string]reaction

func main() {
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}

	inputReader := bufio.NewReader(f)

	reactions := reactions{}

	for {
		line, err := inputReader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		react := reaction{}
		fl := strings.Split(line, "=>")

		for _, v := range strings.Split(fl[0], ",") {
			amt := strings.Split(strings.TrimSpace(v), " ")
			i, err := strconv.Atoi(amt[0])
			if err != nil {
				panic(err)
			}
			react.materials = append(react.materials, material{strings.TrimSpace(amt[1]), i})
		}
		names := strings.Split(strings.TrimSpace(fl[1]), " ")
		amt, err := strconv.Atoi(names[0])
		if err != nil {
			panic(err)
		}
		react.amt = amt
		reactions[names[1]] = react
	}
	fmt.Println(reactions)
	rawMaterials := (&reactions).getOreReq("FUEL")
	total := 0
	for mat, amt := range rawMaterials {

		mult := int(math.Ceil(float64(amt) / float64(reactions[mat].amt)))
		fmt.Println(mat, amt, reactions[mat].amt, mult)
		total += reactions[mat].materials[0].amt * mult
	}

	fmt.Printf("Ore required for fuel: %d\n", total)
}

// Returns amount of product made and amount of ore required
func (reacts *reactions) getOreReq(mat string) map[string]int {
	if (*reacts)[mat].materials[0].name == "ORE" {
		return nil
	}

	out := map[string]int{}
	for _, item := range (*reacts)[mat].materials {
		ret := reacts.getOreReq(item.name)
		if ret == nil {
			out[item.name] += item.amt
		} else {
			for k, v := range ret {
				out[k] += (v * (item.amt / (*reacts)[item.name].amt))
				if item.amt%(*reacts)[item.name].amt != 0 {
					out[k] += v
				}
			}
		}
	}
	fmt.Println(mat, out)
	return out
}
