package d03

import (
	"math"
	"strconv"
	"strings"

	"github.com/jzimbel/adventofcode-go/solutions"
)

type dir [2]int

type point struct {
	x int
	y int
}

type vec struct {
	dir dir
	mag int
}

type drawingGrid struct {
	c point          // cursor
	g map[point]bool // grid
}

var (
	origin = point{}
	dirs   = map[byte]dir{
		'D': dir{0, 1},
		'L': dir{-1, 0},
		'R': dir{1, 0},
		'U': dir{0, -1},
	}
)

func newVec(descriptor string) *vec {
	mag, _ := strconv.Atoi(descriptor[1:])
	return &vec{
		dir: dirs[descriptor[0]],
		mag: mag,
	}
}

func newDrawingGrid() *drawingGrid {
	return &drawingGrid{g: make(map[point]bool)}
}

func (g *drawingGrid) draw(v *vec) {
	xInc, yInc := v.dir[0], v.dir[1]
	x, y := g.c.x+xInc, g.c.y+yInc
	for i := 0; i < v.mag; i, x, y = i+1, x+xInc, y+yInc {
		g.g[point{x, y}] = true
	}
	g.c = point{x - xInc, y - yInc}
}

func (g *drawingGrid) intersect(g2 *drawingGrid) (shared []point) {
	for p1 := range g.g {
		if _, ok := g2.g[p1]; ok {
			shared = append(shared, p1)
		}
	}
	return
}

func getSteps(input string) [][]*vec {
	steps := make([][]*vec, 2)
	wires := strings.Split(input, "\n")
	for i, wire := range wires {
		descriptors := strings.Split(wire, ",")
		for _, d := range descriptors {
			steps[i] = append(steps[i], newVec(d))
		}
	}
	return steps
}

func drawPath(c chan *drawingGrid, steps []*vec) {
	defer close(c)
	grid := newDrawingGrid()
	for _, step := range steps {
		grid.draw(step)
	}
	c <- grid
}

func manhattan(p1 point, p2 point) int {
	return int(
		math.Abs(float64(p1.x-p2.x)) +
			math.Abs(float64(p1.y-p2.y)),
	)
}

func part1(steps [][]*vec) (minDist int) {
	c1, c2 := make(chan *drawingGrid), make(chan *drawingGrid)
	go drawPath(c1, steps[0])
	go drawPath(c2, steps[1])
	grid1 := <-c1
	grid2 := <-c2
	crosses := grid1.intersect(grid2)
	minDist = manhattan(origin, crosses[0])
	for _, p := range crosses[1:] {
		if dist := manhattan(origin, p); dist < minDist {
			minDist = dist
		}
	}
	return
}

func part2(steps [][]*vec) (stepCount int) {
	return
}

// Solve provides the day 3 puzzle solution.
func Solve(input string) (*solutions.Solution, error) {
	steps := getSteps(input)
	return &solutions.Solution{Part1: part1(steps), Part2: part2(steps)}, nil
}
