package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var rex = regexp.MustCompile("(\\d+) (\\w+)")

type Rec struct {
	name string
	q    int
}

type State struct {
	Spent     int
	Remaining map[string]int
	Missing   []Rec
}

func calc1(recipies map[string][]Rec) int {
	st := State{
		0,
		make(map[string]int),
		recipies["FUEL"][:len(recipies["FUEL"])-1]}

	for len(st.Missing) > 0 {
		cur := st.Missing[0]
		st.Missing = st.Missing[1:]

		if st.Remaining[cur.name] >= cur.q {
			st.Remaining[cur.name] -= cur.q
			continue
		}

		cur.q -= st.Remaining[cur.name]
		st.Remaining[cur.name] = 0

		parts := recipies[cur.name]
		unitq := parts[len(parts)-1].q

		q := int(math.Ceil(float64(cur.q) / float64(unitq)))

		for i := 0; i < len(parts)-1; i++ {
			sub := parts[i]
			if sub.name == "ORE" {
				st.Spent += sub.q * q
			} else {
				st.Missing = append(
					st.Missing,
					Rec{sub.name, sub.q * q})
			}
		}

		st.Remaining[cur.name] += (q * unitq) - cur.q
	}

	return st.Spent
}

func calc2(recipies map[string][]Rec, maxSpent int) int {
	st := State{
		0,
		make(map[string]int),
		make([]Rec, 0)}
	cnt := 0
	for {
		lastSpent := st.Spent
		for i := 0; i < len(recipies["FUEL"])-1; i++ {
			st.Missing = append(st.Missing, recipies["FUEL"][i])
		}

		for len(st.Missing) > 0 {
			cur := st.Missing[0]
			st.Missing = st.Missing[1:]

			if st.Remaining[cur.name] >= cur.q {
				st.Remaining[cur.name] -= cur.q
				continue
			}

			cur.q -= st.Remaining[cur.name]
			st.Remaining[cur.name] = 0

			parts := recipies[cur.name]
			unitq := parts[len(parts)-1].q

			q := int(math.Ceil(float64(cur.q) / float64(unitq)))

			for i := 0; i < len(parts)-1; i++ {
				sub := parts[i]
				if sub.name == "ORE" {
					st.Spent += sub.q * q
				} else {
					st.Missing = append(
						st.Missing,
						Rec{sub.name, sub.q * q})
				}
			}

			st.Remaining[cur.name] += (q * unitq) - cur.q
		}

		if st.Spent > maxSpent {
			return cnt

		}
		cnt += 1
	}
}

func parse(lines []string) map[string][]Rec {
	recs := make(map[string][]Rec)

	for _, line := range lines {
		arr := rex.FindAllStringSubmatch(line, -1)
		row := make([]Rec, 0)
		for i := 0; i < len(arr); i++ {
			n, _ := strconv.Atoi(arr[i][1])
			rec := Rec{arr[i][2], n}
			row = append(row, rec)
		}

		recs[row[len(row)-1].name] = row
	}
	return recs
}

func main() {
	raw, _ := ioutil.ReadFile(os.Args[1])
	lines := strings.Split(strings.Trim(string(raw), "\n "), "\n")
	recs := parse(lines)
	fmt.Println("Result1:", calc1(recs))
	fmt.Println("Result2:", calc2(recs, 1000000000000))
}
