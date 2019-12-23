package main

import (
	. "./intcode"
	"fmt"
	"os"
	"sync"
	"time"
)

type PacketIO struct {
	sync.Mutex

	Addr int

	In  []int
	Out []int

	FromGW <-chan [3]int
	ToGW   chan<- [3]int
}

func NewPacketIO(addr int, fromGW, toGW chan [3]int) *PacketIO {
	in := make([]int, 0)
	in = append(in, addr)

	out := make([]int, 0)
	return &PacketIO{
		Addr:   addr,
		In:     in,
		Out:    out,
		FromGW: fromGW,
		ToGW:   toGW,
	}
}

func (io *PacketIO) Receive() (int, bool) {
	io.Lock()
	defer io.Unlock()
	if len(io.In) < 1 {
		return -1, true
	}

	//fmt.Println("XXX", io, io.In)
	v := io.In[0]
	io.In = io.In[1:]
	return v, true
}

func (io *PacketIO) Send(v int) bool {
	io.Lock()
	defer io.Unlock()
	io.Out = append(io.Out, v)

	if len(io.Out) > 2 {
		io.ToGW <- [3]int{io.Out[0], io.Out[1], io.Out[2]}
		io.Out = io.Out[3:]
	}
	return true
}

func (io *PacketIO) Finished() {
}

func (io *PacketIO) Go() {
	for {
		select {
		case packet := <-io.FromGW:
			io.Lock()
			io.In = append(io.In, packet[1], packet[2])
			io.Unlock()
		}

	}

}

type Network struct {
	N     int
	Froms []chan [3]int
	ToGW  chan [3]int
	VMs   []*VM
}

func NewNetwork(prog IntProg, n int) Network {
	froms := make([]chan [3]int, 0)
	toGW := make(chan [3]int, 1000)
	vms := make([]*VM, 0)
	for i := 0; i < n; i++ {
		fromGW := make(chan [3]int, 1000)
		froms = append(froms, fromGW)
		pio := NewPacketIO(i, fromGW, toGW)
		go pio.Go()
		vm := &VM{
			prog.Clone(),
			Position(0), Position(0),
			Running,
			pio}
		go vm.Run()
		vms = append(vms, vm)
	}
	return Network{
		n,
		froms,
		toGW,
		vms}

}

func main() {
	prog := LoadIntProg(os.Args[1])
	net := NewNetwork(prog, 50)

	var packet [3]int
	done := false
	for !done {
		select {
		case packet = <-net.ToGW:
			if packet[0] == 255 {
				done = true
				break
			}
			net.Froms[packet[0]] <- packet
		}
	}

	fmt.Println("Result1:", packet[2])

	net = NewNetwork(prog, 50)
	var nat, last [3]int
	cnt := 0
	done = false
	var v int
	for !done {
		select {
		case packet = <-net.ToGW:
			if packet[0] == 255 {
				fmt.Println("NAT", packet)
				nat = packet
				continue
			}
			net.Froms[packet[0]] <- packet

		case <-time.After(3 * time.Second):
			fmt.Println("REB", cnt, last, nat)
			if cnt > 0 && last[2] == nat[2] {
				v = last[2]
				done = true
				break
			}
			net.Froms[0] <- nat
			last = nat
			cnt += 1
		}
	}

	fmt.Println("Result2:", v)
}
