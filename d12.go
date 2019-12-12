package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var rex = regexp.MustCompile("x=([-\\d]+), y=([-\\d]+), z=([-\\d]+)")

type Vec [3]int

func pad(v int) string {
	if v > -1 {
		return " " + strconv.Itoa(v)
	}
	return strconv.Itoa(v)
}

func (v Vec) String() string {
	return fmt.Sprintf("<x=%v y=%v z=%v>", pad(v[0]), pad(v[1]), pad(v[2]))
}

func parseLine(line string) Vec {
	arr := rex.FindStringSubmatch(line)
	x, _ := strconv.Atoi(arr[1])
	z, _ := strconv.Atoi(arr[2])
	y, _ := strconv.Atoi(arr[3])

	return Vec{x, y, z}
}

func abs(x int) int {
	if x > 0 {
		return x
	}
	return -x
}

func sim(bodies []Vec, velocity []Vec, maxt int) ([]Vec, []Vec) {
	for t := 0; t < maxt; t++ {
		for i := 0; i < len(bodies); i++ {
			for j := i + 1; j < len(bodies); j++ {
				bi, bj := bodies[i], bodies[j]
				for z := 0; z < 3; z++ {
					if bi[z] < bj[z] {
						velocity[i][z] += 1
						velocity[j][z] -= 1
					} else if bi[z] > bj[z] {
						velocity[i][z] -= 1
						velocity[j][z] += 1
					}
				}
			}
		}

		for i := 0; i < len(bodies); i++ {
			bodies[i] = Vec{bodies[i][0] + velocity[i][0], bodies[i][1] + velocity[i][1], bodies[i][2] + velocity[i][2]}
		}
	}
	return bodies, velocity
}

func matching(d []int) bool {
	if len(d)%2 != 0 {
		return false
	}
	off := len(d) / 2
	for i := 0; i < off; i++ {
		if d[i] != d[off+i] {
			return false
		}
	}
	return true
}

func findpatt(bodies []Vec, a, b int) int {
	hist := make([]int, 0)
	hist = append(hist, bodies[a][b])
	velocity := make([]Vec, len(bodies))

	t := 0
	for {
		for i := 0; i < len(bodies); i++ {
			for j := i + 1; j < len(bodies); j++ {
				bi, bj := bodies[i], bodies[j]
				for z := 0; z < 3; z++ {
					if bi[z] < bj[z] {
						velocity[i][z] += 1
						velocity[j][z] -= 1
					} else if bi[z] > bj[z] {
						velocity[i][z] -= 1
						velocity[j][z] += 1
					}
				}
			}
		}

		for i := 0; i < len(bodies); i++ {
			bodies[i] = Vec{bodies[i][0] + velocity[i][0], bodies[i][1] + velocity[i][1], bodies[i][2] + velocity[i][2]}
		}

		hist = append(hist, bodies[a][b])

		if len(hist) > 3 && matching(hist) {
			return t/2 + 1
		}

		t += 1
	}
}

func gcd(a, b int) int {
	for {
		if b == 0 {
			break
		} else if a > b {
			a, b = b, a%b
		} else {
			a, b = a, b%a
		}
	}
	return a
}
func lcm(a, b int) int {
	return (a / gcd(a, b)) * b
}

func main() {
	raw, _ := ioutil.ReadFile(os.Args[1])
	lines := strings.Split(strings.Trim(string(raw), "\n "), "\n")

	bodies := make([]Vec, 0, len(lines))
	for _, line := range lines {
		bodies = append(bodies, parseLine(line))
	}

	velocity := make([]Vec, len(bodies))
	bodies, velocity = sim(bodies, velocity, 1000)

	total := 0
	for i := 0; i < len(bodies); i++ {
		total += (abs(bodies[i][0]) + abs(bodies[i][1]) + abs(bodies[i][2])) * (abs(velocity[i][0]) + abs(velocity[i][1]) + abs(velocity[i][2]))
	}

	fmt.Println("Result1:", total)

	v := 1
	for i := 0; i < len(bodies); i++ {
		for j := 0; j < 3; j++ {
			bodies = make([]Vec, 0, len(lines))
			for _, line := range lines {
				bodies = append(bodies, parseLine(line))
			}

			a := findpatt(bodies, i, j)
			v = lcm(v, a)
		}
	}
	fmt.Println("Result2:", v)
}
