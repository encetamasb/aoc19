package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Node struct {
	id   string
	to   []*Node
	from []*Node
}

type Head struct {
	node *Node
	v    int
}

func calc1(m map[string]*Node) int {
	sum := 0
	q := make([]Head, 0)
	q = append(q, Head{m["COM"], 0})
	for len(q) > 0 {
		sum += 1 // direct

		h := q[0]
		q = q[1:]

		sum += (h.v - 1) // indirect

		for i := 0; i < len(h.node.to); i++ {
			q = append(q, Head{h.node.to[i], h.v + 1})

		}
	}
	return sum
}

func calc2(m map[string]*Node) int {
	q := make([]Head, 0)
	q = append(q, Head{m["YOU"].from[0], 0})
	visited := make(map[string]bool)
	for len(q) > 0 {
		h := q[0]
		q = q[1:]

		visited[h.node.id] = true

		for i := 0; i < len(h.node.to); i++ {
			if h.node.to[i].id == "SAN" {
				return h.v
			}
			if _, ok := visited[h.node.to[i].id]; !ok {
				q = append(q, Head{h.node.to[i], h.v + 1})
			}
		}

		for i := 0; i < len(h.node.from); i++ {
			if h.node.from[i].id == "SAN" {
				return h.v
			}

			if _, ok := visited[h.node.from[i].id]; !ok {
				q = append(q, Head{h.node.from[i], h.v + 1})
			}
		}
	}
	panic("ops")
}

func load(fn string) map[string]*Node {
	raw, _ := ioutil.ReadFile(fn)
	lines := strings.Split(strings.Trim(string(raw), "\n "), "\n")
	m := make(map[string]*Node)
	for _, line := range lines {
		parts := strings.Split(line, ")")

		to, ok := m[parts[1]]
		if !ok {
			to = &Node{parts[1], make([]*Node, 0), make([]*Node, 0)}
			m[parts[1]] = to
		}

		from, ok := m[parts[0]]
		if !ok {
			from = &Node{parts[0], make([]*Node, 0), make([]*Node, 0)}
			m[parts[0]] = from
		}

		from.to = append(from.to, to)
		to.from = append(to.from, from)
	}
	return m
}

func main() {
	m := load(os.Args[1])
	fmt.Println("Result1:", calc1(m))
	fmt.Println("Result2:", calc2(m))
}
