package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

var comparePaths [][]string = [][]string{}

func main() {

	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	fileReader := bufio.NewReader(f)

	objects := map[string]int{}

	paths := map[string][]string{}

	for {
		path, err := fileReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		orbit := strings.Split(path[:len(path)-1], ")")
		if _, ok := paths[orbit[0]]; !ok {
			paths[orbit[0]] = []string{}
		}
		paths[orbit[0]] = append(paths[orbit[0]], orbit[1])
	}
	f.Close()

	fmt.Println(processOrbits("COM", &paths, &objects, []string{}))
	fmt.Println(minIntersect())
}

func processOrbits(target string, paths *map[string][]string, objects *map[string]int, path []string) int {
	total := (*objects)[target]
	path = append(path, target)
	if _, ok := (*paths)[target]; !ok {
		switch target {
		case "YOU", "SAN":
			fmt.Printf("----- %s -----\n%v\n", target, path)
			comparePaths = append(comparePaths, path)
		}
		return (*objects)[target]
	}
	for _, child := range (*paths)[target] {
		(*objects)[child] = (*objects)[target] + 1

		total += processOrbits(child, paths, objects, path)

	}
	return total
}

func minIntersect() int {
	for i := 0; i < len(comparePaths[0]); i++ {
		if !reflect.DeepEqual(comparePaths[0][:i], comparePaths[1][:i]) {
			return len(comparePaths[0][i-1:]) + len(comparePaths[1][i-1:])
		}
	}
	return len(comparePaths[0])
}
