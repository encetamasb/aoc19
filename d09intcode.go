package main

import (
	. "./intcode"
	"fmt"
	"os"
)

func main() {
	prog := LoadIntProg(os.Args[1])
	in := make(chan int)
	out := make(chan int)

	vm := &VM{prog.Clone(), Position(0), Position(0), Running, ChannelInput{in}, ChannelOutput{out}}
	go vm.Run()
	in <- 1
	fmt.Println("Result1:", <-out)

	vm = &VM{prog.Clone(), Position(0), Position(0), Running, ChannelInput{in}, ChannelOutput{out}}
	go vm.Run()
	in <- 2
	fmt.Println("Result2:", <-out)
}
