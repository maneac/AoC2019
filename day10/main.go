package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"sort"
)

type point struct {
	x    float64
	y    float64
	rawx int
	rawy int
	size float64
}

type points []point

func (p points) Len() int      { return len(p) }
func (p points) Swap(a, b int) { p[a], p[b] = p[b], p[a] }
func (p points) Less(a, b int) bool {
	if p[a].x == 0 && p[a].y == 0 {
		return false
	}
	if p[b].x == 0 && p[b].y == 0 {
		return true
	}
	angA := math.Acos(p[a].y)
	angB := math.Acos(p[b].y)
	if p[a].x < 0 {
		angA = 360 - angA
	}
	if p[b].x < 0 {
		angB = 360 - angB
	}
	return angA < angB
	// if p[a].x == 0 {
	// 	if p[a].y == 0 {
	// 		return false
	// 	} else if p[a].y > 0 {
	// 		return true
	// 	}
	// }
	// aSec := getSector(p[a])
	// bSec := getSector(p[b])
	// if aSec != bSec {
	// 	return aSec < bSec
	// }
	// switch aSec {
	// case 1, 4:
	// 	return p[a].y < p[b].y
	// case 2, 3:
	// 	return p[a].y > p[b].y
	// }
	// return false
}

func getSector(p point) int {
	if p.x == 0 && p.y == 0 {
		return -1
	}
	if p.x >= 0 {
		if p.y > 0 {
			return 1
		}
		return 2
	}
	if p.y < 0 {
		return 3
	}
	return 4
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	inputReader := bufio.NewReader(f)

	height := 0
	asteroids := points{}
	for {
		x := 0
		line, err := inputReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		strByte := []byte(line[:len(line)-1])
		for _, b := range strByte {
			if b == '#' {
				asteroids = append(asteroids, point{float64(x), float64(height), x, height, 0})
			}
			x++
		}
		height++
	}

	best := 0
	bestSeen := points{}

	for a := range asteroids {
		seen := points{}
		for b := range asteroids {
			dx := float64(asteroids[b].x - asteroids[a].x)
			dy := float64(asteroids[a].y - asteroids[b].y)
			var nb point
			if b != a {
				nb = normalise(point{dx, dy, asteroids[b].rawx, asteroids[b].rawy, 0})
			} else {
				nb = point{0, 0, asteroids[b].rawx, asteroids[b].rawy, 0}
			}
			hasSeen := false
			for s := range seen {
				if fmt.Sprintf("%0.05f,%0.05f", seen[s].x, seen[s].y) == fmt.Sprintf("%0.05f,%0.05f", nb.x, nb.y) {
					hasSeen = true
					if nb.size < seen[s].size {
						seen[s] = nb
					}
				}
			}
			if !hasSeen {
				seen = append(seen, nb)
			}
		}
		if len(seen) > len(bestSeen) {
			bestSeen = seen
			best = a
		}
	}

	fmt.Printf("Most seen from: (%v) - %d\n", asteroids[best], len(bestSeen)-1)

	sort.Sort(bestSeen)
	strOut := ""
	for i := range bestSeen {
		strOut += fmt.Sprintf("%v\n", bestSeen[i])
	}
	ioutil.WriteFile("./bestSeen.txt", []byte(strOut), 0777)
	fmt.Println(bestSeen[199])
	// if len(bestSeen) >= 200 {
	// 	fmt.Printf("200th best: %v\n", bestSeen[198:202])
	// }
}

func normalise(p point) point {
	p.size = math.Sqrt(math.Pow(p.x, 2) + math.Pow(p.y, 2))
	if p.size == 0 {
		fmt.Println(p)
	}
	p.x /= p.size
	p.y /= p.size

	return p
}
