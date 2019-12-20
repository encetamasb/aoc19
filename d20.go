package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Field string

const (
	Wall    Field = "#"
	Empty   Field = "."
	Nothing Field = " "
)

func (f Field) IsWall() bool {
	return f == Wall
}

func (f Field) IsEmpty() bool {
	return f == Empty
}

func (f Field) IsNothing() bool {
	return f == Nothing
}

func (f Field) IsLabel() bool {
	return f >= Field("A") && f <= Field("Z")
}

type Fields []Field

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

type Map struct {
	W, H    int
	Fields  []Field
	Portals map[Pair]Pair
	Labels  map[Pair]string
}

func (m Map) At(p Pair, def Field) Field {
	if m.H > p.y && p.y >= 0 && m.W > p.x && p.x >= 0 {
		return m.Fields[p.y*m.W+p.x]
	}
	return def
}

func (m Map) Show() {
	fmt.Println()
	for j := 0; j < m.H; j++ {
		for i := 0; i < m.W; i++ {
			pos := Pair{i, j}
			fmt.Print(m.At(pos, Nothing))
		}
		fmt.Println()
	}
	fmt.Println()
}

func LoadMap(path string) Map {
	raw, _ := ioutil.ReadFile(path)
	lines := strings.Split(strings.Trim(string(raw), "\n"), "\n")

	w, h := len(lines[0]), len(lines)

	fields := make(Fields, 0, w*h)
	for _, line := range lines {
		for _, c := range line {
			f := Field(string(c))
			// fmt.Print(f)
			fields = append(fields, f)
		}
		// fmt.Println()
	}
	fmt.Println()

	labels := make(map[Pair]string)
	portals := make(map[Pair]Pair)

	m := Map{w, h, fields, portals, labels}

	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			if m.At(Pair{i, j}, Nothing).IsEmpty() {
				if m.At(Pair{i - 1, j}, Wall).IsLabel() {
					labels[Pair{i, j}] = string(lines[j][i-2]) + string(lines[j][i-1])
				}
				if m.At(Pair{i + 1, j}, Wall).IsLabel() {
					labels[Pair{i, j}] = string(lines[j][i+1]) + string(lines[j][i+2])
				}
				if m.At(Pair{i, j - 1}, Wall).IsLabel() {
					labels[Pair{i, j}] = string(lines[j-2][i]) + string(lines[j-1][i])
				}
				if m.At(Pair{i, j + 1}, Wall).IsLabel() {
					labels[Pair{i, j}] = string(lines[j+1][i]) + string(lines[j+2][i])
				}
			}
		}
	}

	for p1, l1 := range labels {
		for p2, l2 := range labels {
			if l1 == l2 && p1 != p2 {
				portals[p1] = p2
				portals[p2] = p1
			}
		}
	}

	return Map{w, h, fields, portals, labels}
}

func main() {
	m := LoadMap(os.Args[1])
	m.Show()
	fmt.Println(m.Labels)
	fmt.Println(m.Portals)
}
