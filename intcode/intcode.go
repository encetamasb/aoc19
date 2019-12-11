package intcode

import (
	//"fmt"
	"io/ioutil"
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

type IO interface {
	Send(int) bool
	Receive() (int, bool)
}

type VM struct {
	Prog  IntProg
	Pos   Position
	Rbase Position
	State State
	Io    IO
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
	if len(vm.Prog) <= int(pos) {
		newprog := make([]int, pos+5000)
		copy(newprog, vm.Prog)
		return newprog
	}
	return vm.Prog
}

func (vm *VM) Read(pos Position, mode Mode) int {
	prog := vm.Prog
	switch mode {
	case PositionMode:
		return prog.At(
			Position(prog.At(pos)))
	case ImmediateMode:
		return prog.At(pos)
	case RelativeMode:
		return prog.At(vm.Rbase + Position(prog.At(pos)))
	}
	panic("ops")
}

func (vm *VM) Write(pos Position, mode Mode, v int) {
	switch mode {
	case PositionMode:
		vm.Prog = vm.EnsureMemory(Position(vm.Prog.At(pos)))
		vm.Prog[vm.Prog.At(pos)] = v
		return
	case RelativeMode:
		vm.Prog = vm.EnsureMemory(vm.Rbase + Position(vm.Prog.At(pos)))
		vm.Prog[vm.Rbase+Position(vm.Prog.At(pos))] = v
		return
	}
	panic("ops")
}

func (vm *VM) Step() *VM {
	pos := vm.Pos
	cur := vm.Prog.At(pos)
	op := extractOp(cur)

	//next := [4]int{vm.Prog.At(pos), vm.Prog.At(pos+1), vm.Prog.At(pos+2), vm.Prog.At(pos+3) }
	//fmt.Printf("//%v %v %v [%v]\n", vm.Pos, vm.Rbase, next, vm.Prog.At(1000))
	//fmt.Printf("//%v %v(%v) %v(%v) %v(%v)\n", op.DE, op.A, next[1], op.B, next[2], op.C, next[3])
	switch op.DE {
	case Add:
		a := vm.Read(pos+1, op.A)
		b := vm.Read(pos+2, op.B)
		vm.Write(pos+3, op.C, a+b)
		vm.Pos += 4
		return vm
	case Mul:
		a := vm.Read(pos+1, op.A)
		b := vm.Read(pos+2, op.B)
		vm.Write(pos+3, op.C, a*b)
		vm.Pos += 4
		return vm
	case In:
		var n int
		n, ok := vm.Io.Receive()
		if !ok {
			vm.State = Error
			return vm
		}
		vm.Write(pos+1, op.A, n)
		vm.Pos += 2
		return vm
	case Out:
		v := vm.Read(pos+1, op.A)
		//fmt.Print(v)
		ok := vm.Io.Send(v)
		if !ok {
			vm.State = Error
			return vm
		}
		vm.Pos += 2
		return vm
	case JmpTrue:
		v := vm.Read(pos+1, op.A)
		if v != 0 {
			vm.Pos = Position(
				vm.Read(pos+2, op.B))
			return vm
		}
		vm.Pos += 3
		return vm
	case JmpFalse:
		v := vm.Read(pos+1, op.A)
		if v == 0 {
			vm.Pos = Position(
				vm.Read(pos+2, op.B))
			return vm
		}
		vm.Pos += 3
		return vm
	case LessThan:
		a := vm.Read(pos+1, op.A)
		b := vm.Read(pos+2, op.B)
		if a < b {
			vm.Write(pos+3, op.C, 1)
		} else {
			vm.Write(pos+3, op.C, 0)
		}
		vm.Pos += 4
		return vm
	case Equals:
		a := vm.Read(pos+1, op.A)
		b := vm.Read(pos+2, op.B)
		if a == b {
			vm.Write(pos+3, op.C, 1)
		} else {
			vm.Write(pos+3, op.C, 0)
		}
		vm.Pos += 4
		return vm
	case AdjRelBase:
		a := vm.Read(pos+1, op.A)
		vm.Rbase += Position(a)
		vm.Pos += 2
		return vm
	case Halt:
		vm.Pos += 1
		vm.State = Halted
		return vm
	}

	vm.State = Error
	return vm
}

func LoadIntProg(path string) IntProg {
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

func (vm *VM) Run() *VM {
	for vm.State == Running {
		vm = vm.Step()
	}

	if vm.State == Error {
		panic(vm.State)
	}

	return vm

}

type ChannelIO struct {
	In  <-chan int
	Out chan<- int
}

func (io ChannelIO) Receive() (int, bool) {
	n, ok := <-io.In
	return n, ok
}

func (io ChannelIO) Send(v int) bool {
	io.Out <- v
	// lazy :)
	return true
}

type ListIO struct {
	In  []int
	Out []int
}

func (io ListIO) Receive() (int, bool) {
	if len(io.In) < 1 {
		return 0, false
	}
	v := io.In[0]
	io.In = io.In[1:]
	return v, true
}

func (io ListIO) Send(v int) bool {
	io.Out = append(io.Out, v)
	return true
}
