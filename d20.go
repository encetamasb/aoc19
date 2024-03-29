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

func (m Map) ShowPath(level int, path []Triplet) {
	fmt.Println("Level:", level)
	visited := make(map[Pair]bool)
	for i := 0; i < len(path); i++ {
		if path[i].level == level {
			visited[Pair{path[i].x, path[i].y}] = true
		}
	}
	fmt.Println()
	for j := 0; j < m.H; j++ {
		for i := 0; i < m.W; i++ {
			pos := Pair{i, j}
			_, exists := visited[pos]
			if exists {
				fmt.Print("@")
			} else {
				fmt.Print(m.At(pos, Nothing))
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func (m Map) ShowPathAround(level int, path []Triplet, p Pair, w, h int) {
	fmt.Println("Level:", level)
	visited := make(map[Pair]bool)
	for i := 0; i < len(path); i++ {
		if path[i].level == level {
			visited[Pair{path[i].x, path[i].y}] = true
		}
	}
	fmt.Println()
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			pos := Pair{p.x + i - w/2, p.y + j - h/2}
			_, exists := visited[pos]
			if exists {
				fmt.Print("@")
			} else {
				fmt.Print(m.At(pos, Nothing))
			}
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

func (m Map) LabelToPos(label string) (Pair, bool) {
	for pos, s := range m.Labels {
		if label == s {
			return pos, true
		}
	}
	return Pair{}, false
}

type Triplet struct{ level, x, y int }

func (m Map) MinPath(from, to string, recursive bool) ([]Triplet, int) {
	start, ok := m.LabelToPos(from)
	if !ok {
		panic("ops")
	}

	end, ok := m.LabelToPos(to)
	if !ok {
		panic("ops")
	}

	type Rec struct {
		p     Pair
		level int
		dist  int
		from  *Rec
	}

	visited := make(map[Triplet]*Rec)
	toPath := func(cur *Rec) []Triplet {
		path := make([]Triplet, 0)
		for cur != nil {
			path = append(path, Triplet{cur.level, cur.p.x, cur.p.y})
			cur = cur.from
		}
		return path
	}

	q := make([]*Rec, 0)
	q = append(q, &Rec{start, 0, 0, nil})

	for len(q) > 0 {
		rec := q[0]
		q = q[1:]

		if recursive {
			//m.ShowPathAround(rec.level, toPath(rec), rec.p, 15, 15)
			//var input string
			//fmt.Scanln(&input)
		}

		prev, exists := visited[Triplet{rec.level, rec.p.x, rec.p.y}]
		if exists {
			if prev.dist > rec.dist {
				visited[Triplet{rec.level, rec.p.x, rec.p.y}] = rec
			} else {
				continue
			}
		} else {
			visited[Triplet{rec.level, rec.p.x, rec.p.y}] = rec
		}

		for i := 0; i < 4; i++ {
			nextp := rec.p.NextTo(Dir(i))
			f := m.At(nextp, Wall)
			if f == Empty {
				q = append(q, &Rec{nextp, rec.level, rec.dist + 1, rec})
			} else if f.IsLabel() && !m.At(rec.p, Nothing).IsLabel() {
				inner := rec.p.x > 5 && rec.p.x < m.W-5 && rec.p.y > 5 && rec.p.y < m.H-5
				if !recursive {
					nextp, exists := m.Portals[rec.p]
					if exists {
						q = append(q, &Rec{nextp, rec.level, rec.dist + 1, rec})
					}
				} else {
					if rec.level == 0 && !inner {
						continue
					}

					if rec.level >= len(m.Portals)/2 {
						continue
					}

					nextp, exists := m.Portals[rec.p]
					if exists {
						if inner {
							q = append(q, &Rec{nextp, rec.level + 1, rec.dist + 1, rec})
						} else {
							q = append(q, &Rec{nextp, rec.level - 1, rec.dist + 1, rec})
						}
					}
				}
			}

		}
	}

	cur := visited[Triplet{0, end.x, end.y}]
	dist := cur.dist
	path := toPath(cur)
	return path, dist
}

func main() {
	m := LoadMap(os.Args[1])
	// m.Show()

	fmt.Println(m.Labels)
	fmt.Println(m.Portals)

	path, min := m.MinPath("AA", "ZZ", false)
	m.ShowPath(0, path)
	fmt.Println("Result1:", min)

	path, min = m.MinPath("AA", "ZZ", true)
	fmt.Println("Result2:", min)
}
