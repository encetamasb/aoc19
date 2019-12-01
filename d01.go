package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func calc(n int) int {
	return (n / 3) - 2
}

func main() {
	raw, _ := ioutil.ReadFile(os.Args[1])
	lines := strings.Split(strings.Trim(string(raw), "\n "), "\n")
	mods := make([]int, len(lines))
	sum := 0
	for i, line := range lines {
		n, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}
		mods[i] = n
		sum += calc(n)
	}

	fmt.Println("Result1:", sum)

	sum2 := 0
	for _, n := range mods {
		for n > 0 {
			n2 := calc(n)
			if n2 > 0 {
				sum2 += n2
			}
			n = n2
		}
	}
	fmt.Println("Result2:", sum2)
}
