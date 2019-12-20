package main

import (
	. "./intcode"
	"fmt"
	"os"
)

func get(prog IntProg, i, j int) int {
	in := make(chan int, 2)
	out := make(chan int)
	vm := &VM{prog.Clone(), Position(0), Position(0), Running, ChannelIO{in, out}}

	go vm.Run()

	in <- i
	in <- j
	return <-out
}

func main() {
	prog := LoadIntProg(os.Args[1])
	sum := 0
	for j := 0; j < 50; j++ {
		for i := 0; i < 50; i++ {
			v := get(prog, i, j)
			sum += v
			// fmt.Print(map[int]string{0: ".", 1: "#"}[v])
		}
		// fmt.Println()
	}
	fmt.Println("Result1:", sum)

	type Row struct{ a, b int }

	j := 3
	rows := make([]Row, 0)
	rows = append(rows, Row{0, 0})
	rows = append(rows, Row{0, 0})
	rows = append(rows, Row{0, 0})
	n := 100
	for {
		last := rows[len(rows)-1]
		a := last.a
		for get(prog, a, j) != 1 {
			a += 1
		}

		b := a
		for get(prog, b, j) != 0 {
			b += 1
		}

		row := Row{a, b}
		rows = append(rows, row)

		if len(rows) >= n && row.b-row.a >= n && rows[j-n+1].b >= a+n {
			fmt.Println("Result2:", 10000*row.a+j-n+1)
			break
		}
		j += 1
	}

}
