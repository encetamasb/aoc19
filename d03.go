package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Direction string

const (
	U Direction = "U"
	R Direction = "R"
	D Direction = "D"
	L Direction = "L"
)

type Turn struct {
	dir    Direction
	length int
}

type Path []Turn
type Pair struct{ x, y int }

func dirToVec(d Direction) Pair {
	return map[Direction]Pair{
		U: Pair{0, -1},
		R: Pair{1, 0},
		D: Pair{0, 1},
		L: Pair{-1, 0},
	}[d]
}

func parsePath(line string) Path {
	path := make(Path, 0)
	rawturns := strings.Split(line, ",")
	for _, raw := range rawturns {
		length, err := strconv.Atoi(raw[1:])
		if err != nil {
			panic(err)
		}

		turn := Turn{Direction(string(raw[0])), length}
		path = append(path, turn)
	}
	return path
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func dist(a, b Pair) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func main() {
	raw, _ := ioutil.ReadFile(os.Args[1])
	lines := strings.Split(strings.Trim(string(raw), "\n "), "\n")

	m := make(map[Pair]int)
	bestpos := Pair{int(^uint(0) >> 2), int(^uint(0) >> 2)}
	for cur, line := range lines {
		path := parsePath(line)
		curpos := Pair{0, 0}
		for _, t := range path {
			vec := dirToVec(t.dir)
			for i := 0; i < t.length; i++ {
				nextpos := Pair{curpos.x + vec.x, curpos.y + vec.y}
				last, ok := m[nextpos]
				if ok && last != cur {
					if dist(bestpos, Pair{0, 0}) > dist(nextpos, Pair{0, 0}) {
						bestpos = nextpos
					}
				}
				m[nextpos] = cur
				curpos = nextpos
			}
		}
	}

	fmt.Println("Result1:", dist(bestpos, Pair{0, 0}))

	m = make(map[Pair]int)
	best := int(^uint(0) >> 1)
	for cur, line := range lines {
		step := 0
		path := parsePath(line)
		curpos := Pair{0, 0}
		for _, t := range path {
			vec := dirToVec(t.dir)
			for i := 0; i < t.length; i++ {
				step += 1
				nextpos := Pair{curpos.x + vec.x, curpos.y + vec.y}
				last, ok := m[nextpos]
				if cur == 0 {
					if !ok {
						m[nextpos] = step
					}
				} else {
					if ok {
						if best > last+step {
							best = last + step
						}
					}
				}
				curpos = nextpos
			}
		}
	}

	fmt.Println("Result2:", best)

}
