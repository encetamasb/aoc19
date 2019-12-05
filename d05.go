package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Position int

type Success bool

const (
	Running Success = true
	Halted  Success = false
)

type Mode int

const (
	PositionMode  Mode = 0
	ImmediateMode Mode = 1
)

type Instr int

const (
	Add      Instr = 1
	Mul      Instr = 2
	In       Instr = 3
	Out      Instr = 4
	JmpTrue  Instr = 5
	JmpFalse Instr = 6
	LessThan Instr = 7
	Equals   Instr = 8
	Halt     Instr = 99
)

type State []int

type Op struct {
	C, B, A Mode // order not by spec!
	DE      Instr
}

func extractOp(o int) Op {
	return Op{
		Mode(o / 10000),
		Mode((o % 10000) / 1000),
		Mode((o % 1000) / 100),
		Instr(o % 100),
	}
}

func (st State) Clone() State {
	state := make([]int, len(st))
	copy(state, st)
	return state
}

func (st State) At(pos Position, mode Mode) int {
	switch mode {
	case PositionMode:
		return st[st[pos]]
	case ImmediateMode:
		return st[pos]
	}
	panic("ops")
}

func step(st State, pos Position) (Position, Success) {
	cur := st[pos]
	op := extractOp(cur)

	//fmt.Println(op)
	switch op.DE {
	case Add:
		a := st.At(pos+1, op.A)
		b := st.At(pos+2, op.B)
		st[st[pos+3]] = a + b
		return pos + 4, Running
	case Mul:
		a := st.At(pos+1, op.A)
		b := st.At(pos+2, op.B)
		st[st[pos+3]] = a * b
		return pos + 4, Running
	case In:
		var s string
		fmt.Scanln(&s)
		n, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		st[st[pos+1]] = n
		return pos + 2, Running
	case Out:
		v := st.At(pos+1, op.A)
		fmt.Print(v)
		return pos + 2, Running
	case JmpTrue:
		v := st.At(pos+1, op.A)
		if v != 0 {
			pos = Position(
				st.At(pos+2, op.B))
			return pos, Running
		}
		return pos + 3, Running
	case JmpFalse:
		v := st.At(pos+1, op.A)
		if v == 0 {
			pos = Position(
				st.At(pos+2, op.B))
			return pos, Running
		}
		return pos + 3, Running
	case LessThan:
		a := st.At(pos+1, op.A)
		b := st.At(pos+2, op.B)
		if a < b {
			st[st[pos+3]] = 1
		} else {
			st[st[pos+3]] = 0
		}
		return pos + 4, Running
	case Equals:
		a := st.At(pos+1, op.A)
		b := st.At(pos+2, op.B)
		if a == b {
			st[st[pos+3]] = 1
		} else {
			st[st[pos+3]] = 0
		}
		return pos + 4, Running
	case Halt:
		return pos + 1, Halted
	}
	panic("ops")
}

func loadProgram(path string) State {
	raw, _ := ioutil.ReadFile(path)
	input := strings.Split(strings.Trim(string(raw), "\n "), ",")

	prog := make(State, 0)
	for i := 0; i < len(input); i++ {
		n, err := strconv.Atoi(input[i])
		if err != nil {
			panic(err)
		}
		prog = append(prog, n)
	}
	return prog
}

func run(st State, pos Position) Position {
	flag := Running
	for flag == Running {
		if len(st) < int(pos+4) {
			//fmt.Printf("//%d %v\n", pos, st[pos:len(st)])

		} else {
			//fmt.Printf("//%d %v\n", pos, st[pos:pos+4])
		}
		pos, flag = step(st, pos)
	}
	return pos

}

func main() {
	//fmt.Println(extractOp(1002))
	//fmt.Println(extractOp(11111))
	//fmt.Println(extractOp(23456))

	prog := loadProgram(os.Args[1])
	//fmt.Println(prog)

	state := prog.Clone()
	run(state, Position(0))

	fmt.Println()
}
