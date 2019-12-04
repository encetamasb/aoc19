package main

import (
	"fmt"
)

func lessEq(a, b [6]int) bool {
	for i := 0; i < 6; i++ {
		if a[i] == b[i] {
			continue
		}
		if a[i] < b[i] {
			return true
		}
		return false
	}
	return true
}

func inRange(a, min, max [6]int) bool {
	return lessEq(min, a) && lessEq(a, max)
}

func hasPair(a [6]int) bool {
	for i := 0; i < 5; i++ {
		if a[i] == a[i+1] {
			return true
		}
	}
	return false
}

func hasPair2(a [6]int) bool {
	for i := 1; i < 4; i++ {
		if a[i-1] != a[i] && a[i] == a[i+1] && a[i] != a[i+2] {
			return true
		}
	}
	if a[2] != a[1] && a[1] == a[0] {
		return true
	}
	if a[3] != a[4] && a[4] == a[5] {
		return true
	}
	return false
}

func do(prev [6]int, pos int, min, max [6]int, hasPair func([6]int) bool) int {
	acc := 0
	from := prev[0]
	if pos > 0 {
		from = prev[pos-1]
	}

	for n := from; n <= 9; n++ {
		next := [6]int{prev[0], prev[1], prev[2], prev[3], prev[4], prev[5]}
		next[pos] = n

		if pos >= 5 {
			if inRange(next, min, max) && hasPair(next) {
				acc += 1
			}
		} else {
			acc += do(next, pos+1, min, max, hasPair)
		}
	}
	return acc
}

func main() {
	min := [6]int{1, 2, 3, 2, 5, 7}
	max := [6]int{6, 4, 7, 0, 1, 5}
	fmt.Println("Result1:", do(min, 0, min, max, hasPair))
	fmt.Println("Result2:", do(min, 0, min, max, hasPair2))
}
