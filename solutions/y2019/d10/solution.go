package d10

import (
	"math"
	"strings"

	"github.com/jzimbel/adventofcode-go/solutions"
)

const (
	// width and height of my puzzle input, used for some slight optimizations
	width  = 24
	height = 24
)

type point struct {
	x int
	y int
}

type grid map[point]struct{}

// memoizedGCD returns a function that gives the greatest common
// denominator of two ints.
func memoizedGCD() (gcd func(int, int) int) {
	// stores memoized results
	gcdCache := make(map[[2]int]int, width*height)

	gcd = func(a, b int) (n int) {
		var key [2]int
		if a < b {
			key = [...]int{a, b}
		} else {
			key = [...]int{b, a}
		}

		var ok bool
		if n, ok = gcdCache[key]; !ok {
			if b != 0 {
				n = gcd(b, a%b)
			} else {
				n = a
			}
			gcdCache[key] = n
		}
		return
	}
	return
}

// axisDistances returns the x and y distances between two points.
func axisDistances(p1, p2 point) (int, int) {
	x, y := int(math.Abs(float64(p1.x-p2.x))), int(math.Abs(float64(p1.y-p2.y)))
	return x, y
}

func part1(g grid, gcd func(int, int) int) (maxVisibleCount int) {
	for p := range g {
		var visibleCount int
		for other := range g {
			if p == other {
				continue
			}
			denom := gcd(axisDistances(p, other))
			xStepSize, yStepSize := (other.x-p.x)/denom, (other.y-p.y)/denom
			var blocked bool
			for i := 1; i < denom; i++ {
				if _, ok := g[point{p.x + i*xStepSize, p.y + i*yStepSize}]; ok {
					blocked = true
					break
				}
			}
			if !blocked {
				visibleCount++
			}
		}
		if visibleCount > maxVisibleCount {
			maxVisibleCount = visibleCount
		}
	}
	return
}

// Solve provides the day 10 puzzle solution.
func Solve(input string) (*solutions.Solution, error) {
	g := make(grid, width*height)
	rows := strings.Split(input, "\n")
	for y := range rows {
		for x := range rows[y] {
			if rows[y][x] == '#' {
				g[point{x, y}] = struct{}{}
			}
		}
	}
	gcd := memoizedGCD()

	return &solutions.Solution{Part1: part1(g, gcd), Part2: nil}, nil
}
