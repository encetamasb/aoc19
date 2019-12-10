package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strings"
)

type Pair struct{ x, y int }
type FPair struct{ x, y float64 }
type Pairs []Pair

type VMap map[int][]int

func (m VMap) put(at int, v int) {
	arr, ok := m[at]
	if !ok {
		arr = make([]int, 0)
	}
	arr = append(arr, v)
	m[at] = arr
}

func unit(a, b Pair) FPair {
	d := dist(a, b)
	u := FPair{float64(a.y-b.y) / d, float64(a.x-b.x) / d}
	return u
}

func dist(a, b Pair) float64 {
	return math.Sqrt(math.Pow(float64(a.x-b.x), 2) + math.Pow(float64(a.y-b.y), 2))
}

func fdist(a, b FPair) float64 {
	return math.Sqrt(math.Pow(a.x-b.x, 2) + math.Pow(a.y-b.y, 2))
}

func rot(a, b FPair) float64 {
	v := math.Atan2(a.x*b.y-a.y*b.x, a.x*b.x+a.y*b.y)
	if v < 0 {
		v += 2 * math.Pi
	}
	return v / math.Pi * 180
}

func main() {
	raw, _ := ioutil.ReadFile(os.Args[1])
	lines := strings.Split(strings.Trim(string(raw), "\n "), "\n")

	mets := make([]Pair, 0)

	for j, line := range lines {
		for i, c := range line {
			if c == '#' {
				mets = append(mets, Pair{i, j})
			}
		}
	}

	vm := make(VMap)
	for i := 0; i < len(mets)-1; i++ {
		for j := i + 1; j < len(mets); j++ {
			inSight := true
			a, b := mets[i], mets[j]
			u := unit(a, b)
			d := dist(a, b)
			for k := 0; k < len(mets); k++ {
				if k == i || k == j {
					continue
				}
				c := mets[k]

				uc := unit(a, c)
				if fdist(u, uc) < 0.0001 && dist(a, c) < d {
					inSight = false
					break
				}

			}

			if inSight {
				vm.put(i, j)
				vm.put(j, i)
			}
		}
	}

	maxv := 0
	maxi := 0
	for i, v := range vm {
		if len(v) > maxv {
			maxv = len(v)
			maxi = i
		}
	}

	fmt.Println("Result1:", maxv)

	type Rec struct {
		u FPair
		r float64
		n float64
		d float64
		i int
	}

	type Records []Rec

	up := FPair{1, 0}

	recs := make(Records, 0, len(mets)-1)
	mstation := mets[maxi]

	for i := 0; i < len(mets); i++ {
		if i == maxi {
			continue
		}
		b := mets[i]
		ro := rot(unit(mstation, b), up)
		recs = append(recs, Rec{unit(mstation, b), ro, 0, dist(mstation, b), i})
	}

	Less := func(i, j int) bool {
		a, b := recs[i], recs[j]
		ar, br := a.r+360*a.n, b.r+360*b.n
		return ar < br || (ar == br && a.d < b.d)
	}

	sort.Slice(recs, Less)
	cnt := 1

	for i := 1; i < len(recs); i++ {
		if recs[i-1].r == recs[i].r {
			recs[i].n = float64(cnt)
			cnt += 1
		} else {
			cnt = 1
		}
	}

	sort.Slice(recs, Less)
	fmt.Println("Result2:", mets[recs[199].i].x*100+mets[recs[199].i].y)

}
