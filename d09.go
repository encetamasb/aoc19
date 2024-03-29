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
	RelativeMode  Mode = 2
)

func (mode Mode) String() string {
	return [3]string{"Pos", "Imm", "Rel"}[mode]
}

type Instr int

const (
	Add        Instr = 1
	Mul        Instr = 2
	In         Instr = 3
	Out        Instr = 4
	JmpTrue    Instr = 5
	JmpFalse   Instr = 6
	LessThan   Instr = 7
	Equals     Instr = 8
	AdjRelBase Instr = 9
	Halt       Instr = 99
)

func (instr Instr) String() string {
	if instr == 99 {
		return "Halt"
	}
	return [11]string{"???", "Add", "Mul", "In", "Out", "JmpTrue", "JmpFalse", "LessThan", "Equals", "AdjRelBase"}[instr]
}

type IntProg []int

type VM struct {
	prog  IntProg
	pos   Position
	rbase Position
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

func (prog IntProg) At(pos Position) int {
	if len(prog) <= int(pos) {
		return 0
	}
	return prog[pos]
}

func (vm *VM) EnsureMemory(pos Position) IntProg {
	if len(vm.prog) <= int(pos) {
		newprog := make([]int, pos+5000)
		copy(newprog, vm.prog)
		return newprog
	}
	return vm.prog
}

func (vm *VM) Read(pos Position, mode Mode) int {
	prog := vm.prog
	switch mode {
	case PositionMode:
		return prog.At(
			Position(prog.At(pos)))
	case ImmediateMode:
		return prog.At(pos)
	case RelativeMode:
		return prog.At(vm.rbase + Position(prog.At(pos)))
	}
	panic("ops")
}

func (vm *VM) Write(pos Position, mode Mode, v int) {
	switch mode {
	case PositionMode:
		vm.prog = vm.EnsureMemory(Position(vm.prog.At(pos)))
		vm.prog[vm.prog.At(pos)] = v
		return
	case RelativeMode:
		vm.prog = vm.EnsureMemory(vm.rbase + Position(vm.prog.At(pos)))
		vm.prog[vm.rbase+Position(vm.prog.At(pos))] = v
		return
	}
	panic("ops")
}

func (vm *VM) step() *VM {
	pos := vm.pos
	cur := vm.prog.At(pos)
	op := extractOp(cur)

	//next := [4]int{vm.prog.At(pos), vm.prog.At(pos+1), vm.prog.At(pos+2), vm.prog.At(pos+3) }
	//fmt.Printf("//%v %v %v [%v]\n", vm.pos, vm.rbase, next, vm.prog.At(1000))
	//fmt.Printf("//%v %v(%v) %v(%v) %v(%v)\n", op.DE, op.A, next[1], op.B, next[2], op.C, next[3])
	switch op.DE {
	case Add:
		a := vm.Read(pos+1, op.A)
		b := vm.Read(pos+2, op.B)
		vm.Write(pos+3, op.C, a+b)
		vm.pos += 4
		return vm
	case Mul:
		a := vm.Read(pos+1, op.A)
		b := vm.Read(pos+2, op.B)
		vm.Write(pos+3, op.C, a*b)
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
		vm.Write(pos+1, op.A, n)
		vm.pos += 2
		return vm
	case Out:
		v := vm.Read(pos+1, op.A)
		fmt.Print(v)
		vm.out = append(vm.out, v)
		vm.pos += 2
		return vm
	case JmpTrue:
		v := vm.Read(pos+1, op.A)
		if v != 0 {
			vm.pos = Position(
				vm.Read(pos+2, op.B))
			return vm
		}
		vm.pos += 3
		return vm
	case JmpFalse:
		v := vm.Read(pos+1, op.A)
		if v == 0 {
			vm.pos = Position(
				vm.Read(pos+2, op.B))
			return vm
		}
		vm.pos += 3
		return vm
	case LessThan:
		a := vm.Read(pos+1, op.A)
		b := vm.Read(pos+2, op.B)
		if a < b {
			vm.Write(pos+3, op.C, 1)
		} else {
			vm.Write(pos+3, op.C, 0)
		}
		vm.pos += 4
		return vm
	case Equals:
		a := vm.Read(pos+1, op.A)
		b := vm.Read(pos+2, op.B)
		if a == b {
			vm.Write(pos+3, op.C, 1)
		} else {
			vm.Write(pos+3, op.C, 0)
		}
		vm.pos += 4
		return vm
	case AdjRelBase:
		a := vm.Read(pos+1, op.A)
		vm.rbase += Position(a)
		vm.pos += 2
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

func (vm *VM) run() *VM {
	for vm.state == Running {
		vm = vm.step()
	}

	if vm.state == Error {
		panic(vm.state)
	}

	return vm

}

func main() {
	prog := loadIntProg(os.Args[1])
	in := make([]int, 1)
	in[0] = 1
	vm := &VM{prog.Clone(), Position(0), Position(0), Running, in, make([]int, 0)}
	vm.run()
	fmt.Println("\nResult1:", vm.out[len(vm.out)-1])

	in[0] = 2
	vm = &VM{prog.Clone(), Position(0), Position(0), Running, in, make([]int, 0)}
	vm.run()
	fmt.Println("\nResult2:", vm.out[len(vm.out)-1])
}
