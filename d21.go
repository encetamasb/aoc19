package main

import (
	. "./intcode"
	"fmt"
	"os"
)

func send(ch chan int, s string) {
	for i := 0; i < len(s); i++ {
		ch <- int(s[i])
	}
	ch <- 10
}

func main() {
	prog := LoadIntProg(os.Args[1])
	in := make(chan int, 1000)
	out := make(chan int)
	vm := &VM{prog.Clone(), Position(0), Position(0), Running, ChannelIO{in, out}}

	go vm.Run()

	send(in, "OR D J")
	send(in, "NOT C T")
	send(in, "AND T J")
	send(in, "NOT A T")
	send(in, "OR T J")
	for i := 0; i < 10; i++ {
		send(in, "OR J J")
	}
	send(in, "WALK")

	for c := range out {
		if c > 255 {
			fmt.Println(c)
		} else {
			fmt.Print(string(c))
		}
	}

	in = make(chan int, 1000)
	out = make(chan int)
	vm = &VM{prog.Clone(), Position(0), Position(0), Running, ChannelIO{in, out}}

	go vm.Run()

	send(in, "OR E J")
	send(in, "OR H J")
	send(in, "AND D J")

	send(in, "NOT T T")
	send(in, "AND A T")
	send(in, "AND B T")
	send(in, "AND C T")
	send(in, "NOT T T")

	send(in, "AND T J")
	for i := 0; i < 6; i++ {
		send(in, "OR J J")
	}
	send(in, "RUN")

	for c := range out {
		if c > 255 {
			fmt.Println(c)
		} else {
			fmt.Print(string(c))
		}
	}
}
