package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type moon struct {
	x, y, z, vx, vy, vz int
}

type moons []moon

func main() {
	var err error
	f, err := os.Open("./t1.txt")
	if err != nil {
		panic(err)
	}

	inData := bufio.NewReader(f)

	isEOF := false

	moonArr := moons{}

	for !isEOF {
		rawNum, err := inData.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				isEOF = true
				break
			} else {
				panic(err)
			}
		}
		nums := strings.Split(string(rawNum[:len(rawNum)-1]), ",")
		tmpMoon := moon{}
		x, _ := strconv.Atoi(nums[0][3:])
		tmpMoon.x = x
		y, _ := strconv.Atoi(nums[1][3:])
		tmpMoon.y = y
		z, _ := strconv.Atoi(nums[2][3 : len(nums[2])-1])
		tmpMoon.z = z
		moonArr = append(moonArr, tmpMoon)
	}

	visited := map[string]struct{}{}
	for {
		stateStr := fmt.Sprintf("%v%v%v%v", moonArr[0], moonArr[1], moonArr[2], moonArr[3])
		if _, ok := visited[stateStr]; ok {
			break
		} else {
			visited[stateStr] = struct{}{}
		}
		//fmt.Println(moonArr)
		(&moonArr).calcGrav()
		(&moonArr).move()
	}
	fmt.Println(len(visited))
	fmt.Printf("Total energy: %d\n", (&moonArr).getEnergy())
}

func (m *moons) calcGrav() {
	m.gravity(0, 1)
	m.gravity(0, 2)
	m.gravity(0, 3)
	m.gravity(1, 2)
	m.gravity(1, 3)
	m.gravity(2, 3)
}

func (m *moons) gravity(i, j int) {
	if (*m)[i].x != (*m)[j].x {
		switch {
		case (*m)[i].x > (*m)[j].x:
			(*m)[i].vx--
			(*m)[j].vx++
		default:
			(*m)[i].vx++
			(*m)[j].vx--
		}
	}
	if (*m)[i].y != (*m)[j].y {
		switch {
		case (*m)[i].y > (*m)[j].y:
			(*m)[i].vy--
			(*m)[j].vy++
		default:
			(*m)[i].vy++
			(*m)[j].vy--
		}
	}
	if (*m)[i].z != (*m)[j].z {
		switch {
		case (*m)[i].z > (*m)[j].z:
			(*m)[i].vz--
			(*m)[j].vz++
		default:
			(*m)[i].vz++
			(*m)[j].vz--
		}
	}
}

func (m *moons) move() {
	for i := 0; i < len(*m); i++ {
		(*m)[i].x += (*m)[i].vx
		(*m)[i].y += (*m)[i].vy
		(*m)[i].z += (*m)[i].vz
	}
}

func (m *moons) getEnergy() int {
	total := 0
	for i := 0; i < len(*m); i++ {
		pot := 0
		if (*m)[i].x < 0 {
			pot -= (*m)[i].x
		} else {
			pot += (*m)[i].x
		}
		if (*m)[i].y < 0 {
			pot -= (*m)[i].y
		} else {
			pot += (*m)[i].y
		}
		if (*m)[i].z < 0 {
			pot -= (*m)[i].z
		} else {
			pot += (*m)[i].z
		}

		kin := 0
		if (*m)[i].vx < 0 {
			kin -= (*m)[i].vx
		} else {
			kin += (*m)[i].vx
		}
		if (*m)[i].vy < 0 {
			kin -= (*m)[i].vy
		} else {
			kin += (*m)[i].vy
		}
		if (*m)[i].vz < 0 {
			kin -= (*m)[i].vz
		} else {
			kin += (*m)[i].vz
		}
		total += (kin * pot)
	}
	return total
}
