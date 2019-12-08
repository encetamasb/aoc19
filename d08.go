package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type M map[byte]int

func (m M) get(b byte, def int) int {
	n, ok := m[b]
	if !ok {
		n = def
	}
	return n
}

func (m M) inc(b byte) {
	_, ok := m[b]
	if !ok {
		m[b] = 0
	}
	m[b] += 1
}

func main() {
	raw, _ := ioutil.ReadFile(os.Args[1])
	content := strings.Trim(string(raw), "\n ")

	W, H := 25, 6
	layerLen := W * H
	//layerLen := 3 * 2
	layerCount := len(content) / layerLen
	minzero := layerLen
	var result1 int
	for j := 0; j < layerCount; j++ {
		m := make(M)
		for i := 0; i < layerLen; i++ {
			r := content[j*layerLen+i]
			m.inc(r)
		}

		zerocnt := m.get('0', 0)
		if zerocnt < minzero {
			minzero = zerocnt
			result1 = m.get('1', 0) * m.get('2', 0)
		}
	}

	fmt.Println("Result1:", result1)

	im := make([][]byte, 0, H)
	for j := 0; j < H; j++ {
		row := make([]byte, 0, W)
		for i := 0; i < W; i++ {
			row = append(row, '2')
		}
		im = append(im, row)
	}

	for z := 0; z < layerCount; z++ {
		for j := 0; j < H; j++ {
			for i := 0; i < W; i++ {
				if im[j][i] != '2' {
					continue
				}

				r := content[z*layerLen+j*W+i]
				im[j][i] = r
			}
		}
	}

	fmt.Println("Result2:")
	for _, row := range im {
		for _, r := range row {
			if r == '1' {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
