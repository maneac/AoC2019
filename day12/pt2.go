package main

import (
	"fmt"
	"reflect"
)

func (m moon) dir(i int) int {
	switch i {
	case 0:
		return m.x
	case 1:
		return m.y
	case 2:
		return m.z
	}
	return 0
}

func day2(ma moons) {
	periods := []int{}
	for i := 0; i < 3; i++ {
		dir := make(minmoons, 0, len(ma))
		for j := range ma {
			dir = append(dir, minmoon{ma[j].dir(i), 0})
		}
		periods = append(periods, findPeriod(dir))
	}
	fmt.Println(periods)
}

type minmoon struct {
	loc, vel int
}

type minmoons []minmoon

func findPeriod(dir minmoons) int {
	initialState := make(minmoons, len(dir))
	copy(initialState, dir)

	steps := 0
	for {
		steps++

		(&dir).gravity(0, 1)
		(&dir).gravity(0, 2)
		(&dir).gravity(0, 3)
		(&dir).gravity(1, 2)
		(&dir).gravity(1, 3)
		(&dir).gravity(2, 3)

		for i := range dir {
			dir[i].loc += dir[i].vel
		}

		if reflect.DeepEqual(dir, initialState) {
			break
		}
	}
	return steps
}

func (m *minmoons) gravity(i, j int) {
	if (*m)[i].loc != (*m)[j].loc {
		switch {
		case (*m)[i].loc > (*m)[j].loc:
			(*m)[i].vel--
			(*m)[j].vel++
		default:
			(*m)[i].vel++
			(*m)[j].vel--
		}
	}
}
