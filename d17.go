package main

import (
	. "./intcode"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Rotation int

const (
	CW  Rotation = 0
	CCW Rotation = 1
)

func (r Rotation) String() string {
	return [...]string{"CW", "CCW"}[r]
}

type Dir int

const (
	Up    Dir = 0
	Right Dir = 1
	Down  Dir = 2
	Left  Dir = 3
)

func (d Dir) String() string {
	return [...]string{"U", "R", "D", "L"}[d]
}

func (d Dir) RotateCW() Dir {
	if d == Left {
		return Up
	}
	return Dir((d + 1) % 4)
}

func (d Dir) RotateCCW() Dir {
	if d == Up {
		return Left
	}
	return Dir((d - 1) % 4)
}

func (d Dir) Rotate(r Rotation) Dir {
	if r == CW {
		return d.RotateCW()
	}
	return d.RotateCCW()
}

func (d Dir) AsPair() Pair {
	return [4]Pair{
		Pair{0, -1},
		Pair{1, 0},
		Pair{0, 1},
		Pair{-1, 0}}[d]
}

type Pair struct{ x, y int }

func (p Pair) Add(v Pair) Pair {
	return Pair{p.x + v.x, p.y + v.y}
}

func (p Pair) NextTo(d Dir) Pair {
	return p.Add(d.AsPair())
}

type Field int

const (
	Empty    = 46
	Scaffold = 35
)

func (f Field) String() string {
	if f == Empty {
		return "."
	}
	if f == Scaffold {
		return "#"
	}
	panic("ops")
}

type Map [][]Field

func (m Map) At(pos Pair, def Field) Field {
	if pos.y >= 0 && pos.y < len(m) && pos.x >= 0 && pos.x < len(m[pos.y]) {
		return m[pos.y][pos.x]
	}
	return def
}

func (m Map) Show() {
	fmt.Println()
	for j := 0; j < len(m); j++ {
		for i := 0; i < len(m[j]); i++ {
			fmt.Print(m[j][i])
		}
		fmt.Println()
	}
	fmt.Println()
}

func (m Map) ShowAround(pos Pair, w, h int) {
	fmt.Println()
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			cur := Pair{pos.x + i - w/2, pos.y + j - h/2}
			if pos != cur {
				fmt.Print(m.At(cur, Empty))
			} else {
				fmt.Print("o")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

type Step struct {
	r Rotation
	n int
}

func (s Step) String() string {
	return [2]string{"R", "L"}[s.r] + "," + strconv.Itoa(s.n)
}

type Steps []Step

func (steps Steps) String() string {
	m := make([]string, 0, len(steps))
	for i := 0; i < len(steps); i++ {
		m = append(m, steps[i].String())
	}
	return strings.Join(m, ",")
}

func Follow(m Map, curp Pair) Steps {
	curd := Up
	curr := CW

	if m.At(curp.NextTo(curd.RotateCCW()), Empty) == Scaffold {
		curr = CCW
	}

	curd = curd.Rotate(curr)

	cnt := 0
	steps := make(Steps, 0)
	for {
		nextp := curp.NextTo(curd)

		if m.At(nextp, Empty) == Scaffold {
			curp = nextp
			cnt += 1
		} else {
			steps = append(steps, Step{curr, cnt})
			cnt = 0

			if m.At(curp.NextTo(curd.RotateCW()), Empty) == Scaffold {
				curr = CW
				curd = curd.RotateCW()
			} else if m.At(curp.NextTo(curd.RotateCCW()), Empty) == Scaffold {
				curr = CCW
				curd = curd.RotateCCW()
			} else {
				break
			}
		}
	}
	return steps
}

func buildWin(steps Steps, offset int, maxLen int, maxSLen int) Steps {
	win := make(Steps, 0)
	size := 0
	for i := offset; i < len(steps); i++ {
		s := steps[i]
		if len(win) < maxLen && size+1+len(s.String()) < maxSLen {
			win = append(win, s)
			size = len(win.String())
		} else {
			break
		}
	}
	return win

}

func (steps Steps) FindABC(maxlen int) (string, string, string, string) {
	m := make(map[string]int)

	// Let's create a pattern/occurence map
	// (decreasing elem count, constant max pattern string length)
	for max := len(steps); max > 1; max-- {
		i := 0
		for i < len(steps) {
			w := buildWin(steps, i, max, maxlen+1)
			if len(w) == max {
				m[w.String()] += 1
			}
			i += 1
		}
	}

	type Rec struct {
		s   string
		cnt int
	}

	c := make([]Rec, 0)
	for k, v := range m {
		c = append(c, Rec{k, v})
	}

	sort.Slice(c, func(i, j int) bool {
		// Longer or more frequent to the front!
		return !(len(c[i].s) < len(c[j].s) || (c[i].cnt < c[j].cnt && len(c[i].s) == len(c[j].s)))
	})

	whole := steps.String()
	minln := 999999
	var mina, minb, minc Rec
	var mincur string
	// Check every combinations of A B C
	for i := 0; i < len(c)-2; i++ {
		for j := i + 1; j < len(c)-1; j++ {
			for k := j + 1; k < len(c); k++ {
				cur := whole
				a, b, c := c[i], c[j], c[k]
				cur = strings.ReplaceAll(cur, a.s, "A")
				cur = strings.ReplaceAll(cur, b.s, "B")
				cur = strings.ReplaceAll(cur, c.s, "C")
				// Replacing with winning a b c triplet should result in smallest main routine
				if minln > len(cur) {
					minln = len(cur)
					mina, minb, minc = a, b, c
					mincur = cur
				}
			}
		}
	}
	// Let's hope for the best
	return mincur, mina.s, minb.s, minc.s
}

func main() {
	prog := LoadIntProg(os.Args[1])
	in := make(chan int)
	out := make(chan int)

	vm := &VM{prog.Clone(), Position(0), Position(0), Running, ChannelIO{in, out}}
	go vm.Run()

	m := make(Map, 0)
	row := make([]Field, 0)
	var curp Pair
	i, j := 0, 0
	for {
		v := <-out

		if v == 94 {
			curp = Pair{i, j}
			v = 35
		}
		if v == 0 {
			m = append(m, row)
			break
		}
		if v == 10 {
			m = append(m, row)
			row = make([]Field, 0)
			j += 1
			i = 0
		} else {
			row = append(row, Field(v))
			i += 1
		}

	}

	sum := 0
	for j := 2; j < len(m); j++ {
		for i := 1; i < len(m[j])-1; i++ {
			a := [5]Field{m[j][i], m[j-1][i-1], m[j-1][i], m[j-2][i], m[j-1][i+1]}
			b := [5]Field{Scaffold, Scaffold, Scaffold, Scaffold, Scaffold}
			if a == b {
				sum += i * (j - 1)
			}
		}
	}
	fmt.Println("Result1:", sum)

	steps := Follow(m, curp)
	fmt.Println("Steps:", steps)

	main, A, B, C := steps.FindABC(20)

	fmt.Println("Main:", main, "\nA:", A, "\nB:", B, "\nC:", C)

	//By Hand version:
	//main := "B,A,B,C,B,A,C,A,C,A\n"
	//A := "R,8,L,12,R,4,R,4\n"
	//B := "R,8,L,10,L,12,R,4\n"
	//C := "R,8,L,10,R,8\n"

	payload := main + "\n" + A + "\n" + B + "\n" + C + "\nn\n"

	prog2 := prog.Clone()
	prog2[0] = 2
	in = make(chan int, len(payload))
	out = make(chan int)

	vm = &VM{prog2, Position(0), Position(0), Running, ChannelIO{in, out}}
	go vm.Run()

	for i := 0; i < len(payload); i++ {
		in <- int(payload[i])
	}

	last := 0
	for {
		v := <-out
		if v == 0 {
			break
		}
		last = v
		// fmt.Print(string(last))
	}
	fmt.Println("Result2:", last)
}
