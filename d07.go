package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Position int

type State int

const (
	Running State = 0
	Halted  State = 1
	Error   State = 2
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

type IntProg []int

type VM struct {
	p     IntProg
	pos   Position
	state State
	in    []int
	out   []int
}

type Op struct {
	C, B, A Mode // order not matching spec!
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

func (p IntProg) Clone() IntProg {
	newp := make([]int, len(p))
	copy(newp, p)
	return newp
}

func (p IntProg) At(pos Position, mode Mode) int {
	switch mode {
	case PositionMode:
		return p[p[pos]]
	case ImmediateMode:
		return p[pos]
	}
	panic("ops")
}

func (vm VM) step() VM {
	p := vm.p
	pos := vm.pos
	cur := p[pos]
	op := extractOp(cur)

	//fmt.Println(op)
	switch op.DE {
	case Add:
		a := p.At(pos+1, op.A)
		b := p.At(pos+2, op.B)
		p[p[pos+3]] = a + b
		vm.pos += 4
		return vm
	case Mul:
		a := p.At(pos+1, op.A)
		b := p.At(pos+2, op.B)
		p[p[pos+3]] = a * b
		vm.pos += 4
		return vm
	case In:
		var n int
		if len(vm.in) > 0 {
			n = vm.in[0]
			vm.in = vm.in[1:]
		} else {
			var s string
			fmt.Scanln(&s)
			n2, err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			}
			n = n2
		}
		p[p[pos+1]] = n
		vm.pos += 2
		return vm
	case Out:
		v := p.At(pos+1, op.A)
		//fmt.Print(v)
		vm.out = append(vm.out, v)
		vm.pos += 2
		return vm
	case JmpTrue:
		v := p.At(pos+1, op.A)
		if v != 0 {
			vm.pos = Position(
				p.At(pos+2, op.B))
			return vm
		}
		vm.pos += 3
		return vm
	case JmpFalse:
		v := p.At(pos+1, op.A)
		if v == 0 {
			vm.pos = Position(
				p.At(pos+2, op.B))
			return vm
		}
		vm.pos += 3
		return vm
	case LessThan:
		a := p.At(pos+1, op.A)
		b := p.At(pos+2, op.B)
		if a < b {
			p[p[pos+3]] = 1
		} else {
			p[p[pos+3]] = 0
		}
		vm.pos += 4
		return vm
	case Equals:
		a := p.At(pos+1, op.A)
		b := p.At(pos+2, op.B)
		if a == b {
			p[p[pos+3]] = 1
		} else {
			p[p[pos+3]] = 0
		}
		vm.pos += 4
		return vm
	case Halt:
		vm.pos += 1
		vm.state = Halted
		return vm
	}

	vm.state = Error
	return vm
}

func loadIntProg(path string) IntProg {
	raw, _ := ioutil.ReadFile(path)
	input := strings.Split(strings.Trim(string(raw), "\n "), ",")

	prog := make(IntProg, 0)
	for i := 0; i < len(input); i++ {
		n, err := strconv.Atoi(input[i])
		if err != nil {
			panic(err)
		}
		prog = append(prog, n)
	}
	return prog
}

func (vm VM) run() VM {
	for vm.state == Running {
		vm = vm.step()
	}

	if vm.state == Error {
		panic(vm.state)
	}

	return vm

}

func getPerms(v [5]int, i int) [][5]int {
	acc := make([][5]int, 0)
	if i > 4 {
		acc = append(acc, v)
		return acc
	}

	acc = append(acc, getPerms(v, i+1)...)
	for j := i + 1; j < len(v); j++ {
		w := v
		w[i] = v[j]
		w[j] = v[i]
		acc = append(acc, getPerms(w, i+1)...)
	}
	return acc
}

func main() {
	prog := loadIntProg(os.Args[1])

	perms := getPerms([5]int{0, 1, 2, 3, 4}, 0)
	maxSignal := 0
	for _, perm := range perms {
		signal := 0
		for i := 0; i < 5; i++ {
			phase := perm[i]

			in := make([]int, 0)
			in = append(in, phase)
			in = append(in, signal)
			vm := VM{prog.Clone(), Position(0), Running, in, make([]int, 0)}

			vm = vm.run()

			signal = vm.out[len(vm.out)-1]
		}

		if signal > maxSignal {
			maxSignal = signal
		}
	}

	fmt.Println("Result1:", maxSignal)

	perms = getPerms([5]int{0, 1, 2, 3, 4}, 0)
	maxSignal = 0
	for _, perm := range perms {
		signal := 0

		vms := make([]VM, 0, 5)
		for i := 0; i < 5; i++ {
			in := make([]int, 0)
			in = append(in, perm[i]+5) // phase
			out := make([]int, 0)
			vms = append(vms, VM{prog.Clone(), Position(0), Running, in, out})
		}

		vms[0].in = append(vms[0].in, signal) //input signal

		i := 0
		for {
			vm := vms[i]
			for vm.state == Running && len(vm.out) < 1 {
				vm = vm.step()
			}
			vms[i] = vm

			if vm.state == Halted {
				signal = vms[4].out[0]
				break
			} else if vm.state == Error {
				panic("ops")
			}

			j := (i + 1) % 5
			vms[j].in = append(vms[j].in, vm.out[len(vm.out)-1])
			vms[j].out = make([]int, 0)
			i = j
		}

		if signal > maxSignal {
			maxSignal = signal
		}
	}

	fmt.Println("Result2:", maxSignal)
}
