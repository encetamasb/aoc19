package main

import (
	. "./intcode"
	"fmt"
	"os"
	"time"
)

type Tile int

const (
	Empty Tile = iota
	Wall
	Block
	Paddle
	Ball
)

func (t Tile) String() string {
	return [...]string{" ", "#", "x", "-", "o"}[t]
}

const (
	Left    = -1
	Neutral = 0
	Right   = 1
)

type Pair [2]int

func draw(cnt int, scores []int, max Pair, m map[Pair]Tile) {
	ball := Pair{}
	fmt.Println("Cnt:", cnt, "Score:", scores[len(scores)-1])
	for j := 0; j < max[1]+1; j++ {
		for i := 0; i < max[0]+1; i++ {
			v, ok := m[Pair{i, j}]
			if !ok {
				v = Empty
			}
			fmt.Print(v)
			if v == Ball {
				ball = Pair{i, j}
			}
		}
		fmt.Println()
	}
	fmt.Println("Ball", ball)
	if cnt > max[0]*max[1] {
		time.Sleep(10 * time.Millisecond)
	}
}

func main() {
	prog := LoadIntProg(os.Args[1])
	in := make(chan int, 2)
	out := make(chan int, 2)

	vm := &VM{prog.Clone(), Position(0), Position(0), Running, ChannelIO{in, out}}
	go vm.Run()

	outp := make([]int, 0, 3)
	cnt := 0
	max := Pair{0, 0}
	for v := range out {
		outp = append(outp, v)
		if len(outp) == 3 {
			if Tile(outp[2]) == Block {
				cnt += 1
			}
			if max[0] < outp[0] {
				max[0] = outp[0]
			}
			if max[1] < outp[1] {
				max[1] = outp[1]
			}
			outp = make([]int, 0, 3)
		}
	}

	fmt.Println("Result1:", cnt)

	prog = prog.Clone()
	prog[0] = 2
	in = make(chan int, 3)
	out = make(chan int, 2)
	vm = &VM{prog.Clone(), Position(0), Position(0), Running, ChannelIO{in, out}}
	go vm.Run()

	m := make(map[Pair]Tile)
	scores := make([]int, 0)
	scores = append(scores, 0)
	balls := make([]Pair, 0)

	in <- Neutral
	in <- Neutral
	in <- Neutral

	cnt = 0
	for v := range out {
		outp = append(outp, v)
		if len(outp) == 3 {
			cnt += 1
			if Tile(outp[2]) == Ball {
				balls = append(balls, Pair{outp[0], outp[1]})
			}
			if outp[0] == -1 && outp[1] == 0 {
				scores = append(scores, outp[2])
			} else {
				m[Pair{outp[0], outp[1]}] = Tile(outp[2])
			}
			if cnt%987 == 0 || cnt > 30000 {
				draw(cnt, scores, max, m)
			}
			if Tile(outp[2]) == Ball && len(balls) > 3 {
				in <- (balls[len(balls)-1][0] - balls[len(balls)-2][0])
			}
			outp = make([]int, 0, 3)
		}
	}
	draw(cnt, scores, max, m)

	fmt.Println("Result2:", scores[len(scores)-1])
}
