package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	Wall  string = "#"
	Empty string = "."
)

type Pair struct{ x, y int }

type Map struct {
	W, H  int
	Tiles [][]string
	Keys  map[string]Pair
	Doors map[string]Pair
	Entry Pair
}

func (m Map) At(p Pair, def string) string {
	if m.H > p.y && p.y >= 0 && m.W > p.x && p.x >= 0 {
		return m.Tiles[p.y][p.x]
	}
	return def
}

func (m Map) ShowAround(p Pair, w, h int) {
	fmt.Println()
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			x, y := p.x+i-w/2, p.y+j-h/2
			fmt.Print(m.Tiles[y][x])
		}
		fmt.Println()
	}
	fmt.Println()
}

func (m Map) Show() {
	fmt.Println()
	for j := 0; j < m.H; j++ {
		for i := 0; i < m.W; i++ {
			fmt.Print(m.Tiles[j][i])
		}
		fmt.Println()
	}
	fmt.Println()
}

func LoadMap(path string) Map {
	raw, _ := ioutil.ReadFile(path)
	lines := strings.Split(strings.Trim(string(raw), "\n "), "\n")

	doors := make(map[string]Pair)
	keys := make(map[string]Pair)
	rows := make([][]string, 0)
	var entry Pair
	for j, line := range lines {
		row := make([]string, 0, len(line))
		for i, c := range line {
			cur := string(c)
			row = append(row, cur)

			if c >= 'a' && c <= 'z' {
				keys[cur] = Pair{i, j}
			}

			if c >= 'A' && c <= 'Z' {
				doors[cur] = Pair{i, j}
			}

			if c == '@' {
				entry = Pair{i, j}
			}
		}
		rows = append(rows, row)
	}

	return Map{len(rows[0]), len(rows), rows, keys, doors, entry}
}

func main() {
	m := LoadMap(os.Args[1])
	m.Show()
	fmt.Println(m)
}
