package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var basePattern = [4]int{0, 1, 0, -1}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func calcElem(v []int, index int, p [4]int, off int) int {
	take := index + 1
	cur := index

	m := +1
	max := len(v) + off
	sum := 0
	for cur < max {
		sum += m * v[cur-off]
		take = take - 1
		if take < 1 {
			cur += index + 2
			take = index + 1
			m = -m
		} else {
			cur += 1
		}
	}
	return abs(sum) % 10
}

func calcNext(v []int, p [4]int, off int) []int {
	w := make([]int, 0)
	for i := off; i < off+len(v); i++ {
		w = append(w, calcElem(v, i, p, off))
	}
	return w
}

func calcNext2(v []int, p [4]int, off int) []int {
	w := make([]int, 0)
	sum := 0
	for i := 0; i < len(v); i++ {
		sum += v[i]
	}

	for i := 0; i < len(v); i++ {
		if i > 0 {
			sum = sum - v[i-1]
			w = append(w, sum%10)
		} else {
			w = append(w, sum%10)
		}
	}
	return w
}

func toInt(v []int) int {
	sum := 0
	m := 1
	for i := len(v) - 1; i > -1; i = i - 1 {
		sum += m * v[i]
		m *= 10
	}
	return sum
}

func main() {
	raw, _ := ioutil.ReadFile(os.Args[1])
	v := make([]int, 0)
	s := strings.Trim(string(raw), "\n ")

	for i := 0; i < len(s); i++ {
		v = append(v, int(s[i])-48)
	}

	for i := 0; i < 100; i++ {
		v = calcNext(v, basePattern, 0)
	}

	fmt.Println("Result1:", toInt(v[:8]))

	v = make([]int, 0)
	for j := 0; j < 10000; j++ {
		for i := 0; i < len(s); i++ {
			v = append(v, int(s[i])-48)
		}
	}

	off := 0
	m := 1
	for i := 6; i > -1; i = i - 1 {
		off += m * v[i]
		m *= 10
	}

	v = v[off:]

	for i := 0; i < 100; i++ {
		v = calcNext2(v, basePattern, off)
	}

	fmt.Println("Result2:", toInt(v[:8]))
}
