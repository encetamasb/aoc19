package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Dir int

const (
	FW Dir = 0
	BW Dir = 1
)

type Card struct {
	v    int
	up   *Card
	down *Card
}

type Cards struct {
	n   int
	dir Dir
	top *Card
}

func (cards Cards) Up(c *Card) *Card {
	if cards.dir == FW {
		return c.up
	}
	return c.down
}

func (cards Cards) Down(c *Card) *Card {
	if cards.dir == FW {
		return c.down
	}
	return c.up
}

func NewCards(n int) Cards {
	top := &Card{0, nil, nil}
	cur := top
	for i := 1; i < n; i++ {
		next := &Card{i, cur, nil}
		cur.down = next
		cur = next
	}
	cur.down = top
	top.up = cur
	return Cards{n, FW, top}
}

func (cards Cards) DealIntoNewStack() Cards {
	cards.top = cards.Up(cards.top)
	cards.dir = Dir(1 - cards.dir)
	return cards
}

func (cards Cards) CutN(n int) Cards {
	if n > 0 {
		cur := cards.top
		for i := 0; i < n; i++ {
			cur = cards.Down(cur)
		}
		cards.top = cur
	} else {
		cur := cards.top
		for i := 0; i < -n; i++ {
			cur = cards.Up(cur)
		}
		cards.top = cur
	}
	return cards
}

func (cards Cards) DealIncrement(n int) Cards {
	t := make([]int, cards.n)
	cur := cards.top
	t[0] = cur.v
	for i := 1; i < cards.n; i++ {
		cur = cards.Down(cur)
		t[(i*n)%cards.n] = cur.v
	}

	cur = cards.top
	for i := 0; i < cards.n; i++ {
		cur.v = t[i]
		cur = cards.Down(cur)
	}

	return cards
}

func (cards Cards) DealIncrement2(n int) Cards {
	cur := cards.top
	curv := cur.v
	curp := 0
	remaining := cards.n
	for remaining > 0 {
		nextp := (curp * n) % cards.n

		if curp == nextp {
			cur.v = curv
			remaining -= 1
			cards.Show()
			cur = cards.Down(cur)
			curv = cur.v
			curp += 1
			continue
		}

		for nextp != curp {
			if nextp < curp {
				cur = cards.Up(cur)
				curp -= 1
			} else {
				cur = cards.Down(cur)
				curp += 1
			}
		}

		curv2 := cur.v
		cur.v = curv
		curv = curv2
		remaining -= 1
		cards.Show()
	}
	return cards
}

func (cards Cards) Show() {
	fmt.Print(cards.top.v, " ")
	cur := cards.Down(cards.top)
	for cur != cards.top {
		fmt.Print(cur.v, " ")
		cur = cards.Down(cur)
	}
	fmt.Println()
}

func test() {
	cards := NewCards(10)
	cards.Show()
	cards = cards.DealIntoNewStack()
	cards.Show()
	cards = cards.DealIntoNewStack()
	cards.Show()
	cards = cards.CutN(3)
	cards.Show()
	cards = cards.CutN(4)
	cards.Show()
	cards = cards.CutN(-4)
	cards.Show()
	cards = cards.CutN(-3)
	cards.Show()
	cards = cards.DealIncrement(3)
	cards.Show()

	shuffles := make([]Shuffle, 0)
	shuffles = append(shuffles, Shuffle{"cut", -8})
	fmt.Println(Fast(0, 10, shuffles))
	fmt.Println(Fast(1, 10, shuffles))
	fmt.Println(Fast(2, 10, shuffles))
	fmt.Println(Fast(3, 10, shuffles))
}


var rex = regexp.MustCompile("([-\\d]+)")

type Shuffle struct {
	t     string
	param int
}

func (s Shuffle) Apply(cards Cards) Cards {
	switch s.t {
	case "cut":
		return cards.CutN(s.param)
	case "stack":
		return cards.DealIntoNewStack()
	case "increment":
		return cards.DealIncrement(s.param)
	}
	panic("ops")
}

func Slow(curpos int, n int, shuffles []Shuffle) int {
	cards := NewCards(n)
	for _, s := range shuffles {
		cards = s.Apply(cards)
	}

	cur := cards.top
	i := 0
	for {
		if cur.v == curpos {
			break
		}
		i += 1
		cur = cards.Down(cur)
	}
	return i
}

func Fast(curpos int, n int, shuffles []Shuffle) int {
	for _, s := range shuffles {
		switch s.t {
		case "cut":
			if s.param > 0 {
				if curpos < s.param {
					curpos = n - s.param + curpos
				} else {
					curpos = curpos - s.param
				}
			} else {
				if curpos >= n+s.param {
					curpos = curpos - (n + s.param)
				} else {
					curpos = curpos - s.param
				}
			}
		case "stack":
			curpos = n - curpos - 1
		case "increment":
			curpos = (curpos * s.param) % n
		}
	}
	return curpos
}

func LoadShuffles(path string) []Shuffle {
	raw, _ := ioutil.ReadFile(path)
	lines := strings.Split(strings.Trim(string(raw), "\n "), "\n")

	shuffles := make([]Shuffle, 0)
	for _, line := range lines {
		//fmt.Println(line)
		if strings.Contains(line, "cut") {
			d := rex.FindString(line)
			n, err := strconv.Atoi(d)
			if err != nil {
				panic(err)
			}
			shuffles = append(shuffles, Shuffle{"cut", n})
		} else if strings.Contains(line, "stack") {
			shuffles = append(shuffles, Shuffle{"stack", 0})
		} else if strings.Contains(line, "increment") {
			d := rex.FindString(line)
			n, err := strconv.Atoi(d)
			if err != nil {
				panic(err)
			}
			shuffles = append(shuffles, Shuffle{"increment", n})
		}
	}

	return shuffles
}



func main() {
	shuffles := LoadShuffles(os.Args[1])
	fmt.Println("Result1:", Slow(2019, 10007, shuffles), Fast(2019, 10007, shuffles))

	n := 119315717514047
	curpos := 2020
	t := 101741582076661
	m := make(map[int]int)
	m[curpos] = 0
	for i := 0; i < t; i++ {
		curpos = Fast(curpos, n, shuffles)
		t0, exists := m[curpos]
		if exists {
			fmt.Println("EXISTS", t0)
			break
		}
		m[curpos] = i + 1
	}
	fmt.Println(curpos)
}
