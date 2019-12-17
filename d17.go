package main

import (
	. "./intcode"
	"fmt"
	"os"
)

func show(m [][]int) {
	fmt.Println()
	for j := 0; j < len(m); j++ {
		for i := 0; i < len(m[j]); i++ {
			fmt.Print(m[j][i])
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	prog := LoadIntProg(os.Args[1])
	in := make(chan int)
	out := make(chan int)

	vm := &VM{prog.Clone(), Position(0), Position(0), Running, ChannelIO{in, out}}
	go vm.Run()

	m := make([][]int, 0)
	row := make([]int, 0)
	for {
		v := <-out
		fmt.Print(string(v))

		if v == 0 {
			m = append(m, row)
			break
		}
		if v == 10 {
			m = append(m, row)
			row = make([]int, 0)
		} else {
			row = append(row, v)
		}

	}
	sum := 0

	for j := 2; j < len(m); j++ {
		for i := 1; i < len(m[j])-1; i++ {
			if [5]int{m[j][i], m[j-1][i-1], m[j-1][i], m[j-2][i], m[j-1][i+1]} == [5]int{35, 35, 35, 35, 35} {
				sum += i * (j - 1)
			}
		}
	}

	fmt.Println("Result1:", sum)

	prog2 = prog.Clone()
	prog2[0] = 2
	in = make(chan int)
	out = make(chan int)

	vm = &VM{prog2, Position(0), Position(0), Running, ChannelIO{in, out}}
	go vm.Run()
}
