package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type point struct {
	x    int
	y    int
	step int
}

type byX []point
type bySteps []point

func (in byX) Len() int      { return len(in) }
func (in byX) Swap(a, b int) { in[a], in[b] = in[b], in[a] }
func (in byX) Less(a, b int) bool {
	if in[a].x == in[b].x {
		return in[a].y < in[b].y
	}
	return in[a].x < in[b].x
}

func (in bySteps) Len() int      { return len(in) }
func (in bySteps) Swap(a, b int) { in[a], in[b] = in[b], in[a] }
func (in bySteps) Less(a, b int) bool {
	return in[a].step < in[b].step
}

func main() {
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	inReader := bufio.NewReader(f)

	wires := make([][]string, 0)

	for {
		line, err := inReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		values := strings.Split(line[:len(line)-1], ",")
		wires = append(wires, values)
	}

	wirePaths := make([][]point, 0)
	for _, wire := range wires {
		wirePaths = append(wirePaths, getPath(wire))
	}

	interPt, dist := findClosestIntersect(wirePaths)
	fmt.Printf("Closest: (%d,%d) - %d\n", interPt.x, interPt.y, dist)
}

func getPath(wire []string) []point {
	visited := make([]point, 0)
	curX, curY := 0, 0
	step := 0
	for i := 0; i < len(wire); i++ {

		switch wire[i][0] {
		case 'R', 'L', 'U', 'D':
		default:
			fmt.Printf("Invalid direction: %v\n", wire[i])
			os.Exit(0)
		}

		dist, err := strconv.Atoi(wire[i][1:])
		if err != nil {
			panic(err)
		}

		switch wire[i][0] {
		case 'R':
			for j := 0; j < dist; j++ {
				visited = append(visited, point{curX, curY, step})
				curX++
				step++
			}
		case 'L':
			for j := 0; j < dist; j++ {
				visited = append(visited, point{curX, curY, step})
				curX--
				step++
			}
		case 'U':
			for j := 0; j < dist; j++ {
				visited = append(visited, point{curX, curY, step})
				curY++
				step++
			}
		case 'D':
			for j := 0; j < dist; j++ {
				visited = append(visited, point{curX, curY, step})
				curY--
				step++
			}
		}
	}
	visited = append(visited[1:], point{curX, curY, len(wire)})
	fmt.Println(visited)
	return visited
}

func findClosestIntersect(paths [][]point) (point, int) {
	if len(paths) != 2 {
		return point{}, 0
	}
	sort.Sort(byX(paths[0]))
	sort.Sort(byX(paths[1]))

	intersection := make([]point, 0)

	for i := 0; i < len(paths[0]); i++ {
		for len(paths[1]) > 0 && paths[1][0].x < paths[0][i].x {
			paths[1] = paths[1][1:]
		}
		if len(paths[1]) == 0 {

			break
		}

		if paths[0][i].x == paths[1][0].x {
			for len(paths[1]) > 0 && (paths[0][i].x == paths[1][0].x) && paths[1][0].y < paths[0][i].y {
				paths[1] = paths[1][1:]
			}
			if len(paths[1]) == 0 {
				break
			}
			if paths[1][0].x == paths[0][i].x && paths[1][0].y == paths[0][i].y {
				pt := paths[1][0]
				pt.step += paths[0][i].step
				intersection = append(intersection, pt)
			}
		}
	}

	//pt1
	// bestIdx, bestDist := 0, -1
	// for idx, point := range intersection {
	// 	dist := abs(point.x) + abs(point.y)
	// 	if bestDist < 0 || dist < bestDist {
	// 		bestDist = dist
	// 		bestIdx = idx
	// 	}
	// }
	// return intersection[bestIdx], bestDist

	//pt2
	sort.Sort(bySteps(intersection))

	return intersection[0], intersection[0].step
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
