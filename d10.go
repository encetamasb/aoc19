package main

import (
	"fmt"
	"strings"
	"os"
	"io/ioutil"
	"math"
)


type Pair struct {x, y int}
type FPair struct {x, y float64}
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
	u := FPair{float64(a.y - b.y)/d, float64(a.x - b.x) / d}
	return u
}

func dist(a, b Pair) float64 {
	return math.Sqrt(math.Pow(float64(a.x - b.x), 2) + math.Pow(float64(a.y - b.y), 2))
}

func fdist(a, b FPair) float64 {
	return math.Sqrt(math.Pow(a.x - b.x, 2) + math.Pow(a.y - b.y, 2))
}

func main() {
	raw, _ := ioutil.ReadFile(os.Args[1])
	lines := strings.Split(strings.Trim(string(raw), "\n "), "\n")
	
	w, h := len(lines[0]), len(lines)
	mets := make([]Pair, 0)

	for j, line := range lines {
		for i, c := range line {
			if c == '#' {
				mets = append(mets, Pair{i, j})
			}
		}
	}
	
	vm := make(VMap)
	for i := 0; i < len(mets)-1; i ++ {
		for j := i + 1; j < len(mets); j ++ {
			inSight := true
			a, b := mets[i], mets[j]
			u := unit(a, b)
			d := dist(a, b)	
			for k := 0; k < len(mets); k ++ {
				if k == i || k == j {
					continue
				}
				c := mets[k]
				
				//fmt.Println(a, b, c, u, d,  unit(a, c), dist(a, c))
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

	//fmt.Println(mets, vm)
	maxv := 0
	maxi := 0
	for i, v := range vm {
		if len(v) > maxv {
			maxv = len(v)
			maxi = i
		}
		if mets[i].x == 11 && mets[i].y == 13 {
		fmt.Println(i, mets[i], len(v))
		}

	}
	
	
	fmt.Println("Result1:", w, h, maxi, mets[maxi], maxv)
}



