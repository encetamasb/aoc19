package main

import (
	. "./intcode"
	"fmt"
	"os"
)

type Dir int

const (
	North Dir = 1
	South Dir = 2
	West  Dir = 3
	East  Dir = 4
	None  Dir = 5
)

func (d Dir) String() string {
	return [...]string{"?", "N", "S", "W", "E"}[d]
}

func (d Dir) Rotate180() Dir {
	return map[Dir]Dir{
		North: South,
		South: North,
		East:  West,
		West:  East,
	}[d]
}

type Field int

const (
	Wall   Field = 0
	Empty  Field = 1
	Tank   Field = 2
	Oxigen Field = 3
)

func (f Field) String() string {
	return [...]string{"#", ".", "T", "O"}[f]
}

type Pair struct{ x, y int }

func (p Pair) Step(d Dir) Pair {
	v := map[Dir]Pair{
		North: Pair{0, -1},
		South: Pair{0, 1},
		West:  Pair{-1, 0},
		East:  Pair{1, 0}}[d]
	return Pair{p.x + v.x, p.y + v.y}
}

type Head struct {
	pos  Pair
	f    Field
	from Dir
}

type Map map[Pair]Head

func show(m Map, p Pair, w, h int) {
	fmt.Println()
	for j := p.y - h/2; j < p.y+h/2; j++ {
		for i := p.x - w/2; i < p.x+w/2; i++ {
			if p.x == i && p.y == j {
				fmt.Print("X")
				continue
			}
			h, ok := m[Pair{i, j}]
			if ok {
				fmt.Print(h.f)
			} else {
				fmt.Print("?")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	prog := LoadIntProg(os.Args[1])
	in := make(chan int)
	out := make(chan int)

	vm := &VM{prog.Clone(), Position(0), Position(0), Running, ChannelIO{in, out}}
	go vm.Run()

	root := Head{Pair{0, 0}, Empty, None}
	m := make(Map)
	m[root.pos] = root
	cur := root
	dir := None
	steps := 0
	r1 := 0
	var tank Pair
	for {
		m, cur, dir, steps = step1(in, out, m, cur, steps)
		if cur.f == Tank {
			//fmt.Println(dir, cur, steps)
			tank = cur.pos
			r1 = steps
		}
		if dir == None {
			break
		}
	}

	//show(m, cur.pos, 60, 45)
	fmt.Println("Result1:", r1)

	t := 0
	q := make([]Pair, 0)
	q = append(q, tank)
	m[cur.pos] = Head{tank, Oxigen, None}
	grow := 0
	for len(q) > 0 {
		m, q, grow = step2(m, q)
		if grow == 0 {
			break

		}
		t += 1
	}

	fmt.Println("Result2:", t)
}

func step1(in chan int, out chan int, m Map, cur Head, steps int) (Map, Head, Dir, int) {
	curpos := cur.pos
	var d Dir
	dirs := [4]Dir{North, East, South, West}
	for _, d := range dirs {
		nextpos := cur.pos.Step(d)
		_, ok := m[nextpos]
		if !ok {
			in <- int(d)
			f := Field(<-out)
			m[nextpos] = Head{nextpos, f, d}

			if f == Wall {
				continue
			}

			cur = m[nextpos]
			steps += 1
			break
		}
	}

	if curpos == cur.pos {
		if cur.from == None {
			d = None
		} else {
			d = cur.from.Rotate180()
			in <- int(d)
			f := Field(<-out)
			steps -= 1
			if f == Wall {
				panic("ops")
			}
			cur = m[cur.pos.Step(d)]
		}
	}
	return m, cur, d, steps
}

func step2(m Map, q []Pair) (Map, []Pair, int) {
	newq := make([]Pair, 0)
	cnt := 0
	for i := 0; i < len(q); i++ {
		dirs := [4]Dir{North, East, South, West}
		for _, d := range dirs {
			nextpos := q[i].Step(d)
			if m[nextpos].f == Empty {
				m[nextpos] = Head{nextpos, Oxigen, None}
				newq = append(newq, nextpos)
				cnt += 1
			}
		}
	}

	return m, newq, cnt
}
