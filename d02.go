package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func step(state *[]int, pos int) (int, bool) {
	st := (*state)
	switch st[pos] {
	case 1:
		st[st[pos+3]] = st[st[pos+1]] + st[st[pos+2]]
		return pos + 4, false
	case 2:
		st[st[pos+3]] = st[st[pos+1]] * st[st[pos+2]]
		return pos + 4, false
	case 99:
		return pos, true
	}
	return -1, true
}

func main() {
	raw, _ := ioutil.ReadFile(os.Args[1])
	input := strings.Split(strings.Trim(string(raw), "\n "), ",")

	prog := make([]int, 0)
	for i := 0; i < len(input); i++ {
		n, err := strconv.Atoi(input[i])
		if err != nil {
			panic(err)
		}
		prog = append(prog, n)
	}

	pos := 0
	halted := false
	state := make([]int, len(prog))
	copy(state, prog)
	state[1] = 12
	state[2] = 2
	for halted == false {
		pos, halted = step(&state, pos)
	}

	fmt.Println("Result1:", state[0])

	var noun, verb int
	found := false
	for noun = 0; noun < 100; noun++ {
		for verb = 0; verb < 100; verb++ {
			state = make([]int, len(prog))
			copy(state, prog)
			state[1] = noun
			state[2] = verb
			pos = 0
			for {
				pos, halted = step(&state, pos)
				if halted {
					break
				}

			}

			if pos > -1 && state[0] == 19690720 {
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	fmt.Println("Result2:", 100*noun+verb)
}
