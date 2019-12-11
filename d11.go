package main

import (
	. "./intcode"
	"fmt"
	"os"
)

type Dir int

const (
	Up Dir = iota
	Right
	Down
	Left
)

type Rot int

const (
	CCW Rot = 0
	CW  Rot = 1
)

func (d Dir) Rotate(r Rot) Dir {
	if r == CW {
		if d == Left {
			return Up
		}
		return Dir((d + 1) % 4)
	}
	if d == Up {
		return Left
	}
	return Dir((d - 1) % 4)
}

func (d Dir) Vec() Pair {
	return [4]Pair{
		Pair{0, -1},
		Pair{1, 0},
		Pair{0, 1},
		Pair{-1, 0}}[d]
}

type Pair struct{ x, y int }

func show(m map[Pair]int, tx, ty, w, h int) {
	for j := ty; j < ty+h; j++ {
		for i := tx; i < tx+w; i++ {
			v, ok := m[Pair{i, j}]
			if !ok {
				v = 0
			}
			if v == 0 {
				fmt.Print(" ")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println()
	}

}

func main() {
	prog := LoadIntProg(os.Args[1])
	in := make(chan int, 2)
	out := make(chan int, 2)

	vm := &VM{prog.Clone(), Position(0), Position(0), Running, ChannelIO{in, out}}
	go vm.Run()

	m := make(map[Pair]int)
	curDir := Up
	curPos := Pair{0, 0}

	for vm.State == Running {
		v, ok := m[curPos]
		if !ok {
			v = 0
		}

		in <- v
		newv := <-out

		m[curPos] = newv

		rot := Rot(<-out)
		curDir = curDir.Rotate(rot)
		off := curDir.Vec()
		curPos = Pair{curPos.x + off.x, curPos.y + off.y}
	}

	fmt.Println("Result1:", len(m))

	in = make(chan int, 2)
	out = make(chan int, 2)

	vm = &VM{prog.Clone(), Position(0), Position(0), Running, ChannelIO{in, out}}
	go vm.Run()

	m = make(map[Pair]int)
	curDir = Up
	curPos = Pair{0, 0}
	m[curPos] = 1

	for vm.State == Running {
		fmt.Println("POS", curPos)
		v, ok := m[curPos]
		if !ok {
			v = 0
		}

		in <- v
		newv := <-out

		m[curPos] = newv

		rot := Rot(<-out)
		curDir = curDir.Rotate(rot)
		off := curDir.Vec()
		curPos = Pair{curPos.x + off.x, curPos.y + off.y}
	}

	show(m, 0, 0, 40, 1)
}
