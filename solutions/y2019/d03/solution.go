package d03

import (
	"math"
	"strconv"
	"strings"

	"github.com/jzimbel/adventofcode-go/solutions"
)

// dir represents a unit vector in one of the cardinal directions.
// dir[0] = x component, dir[y] = y component
type dir [2]int

type point struct {
	x int
	y int
}

// record of whether a point has been visited and how long the path was when it was first visited.
type record struct {
	visited    bool
	pathLength int
}

// cursor records a point and the path distance traveled to reach that point.
type cursor struct {
	p    point
	dist int
}

// grid holds a map that records which points have been visited by which paths,
// as well as cursors for the two paths.
type grid struct {
	g map[point][2]record
	c [2]cursor
}

var (
	origin = point{}
	dirs   = map[byte]*dir{
		'D': &dir{0, 1},
		'L': &dir{-1, 0},
		'R': &dir{1, 0},
		'U': &dir{0, -1},
	}
)

func newGrid() *grid {
	return &grid{
		make(map[point][2]record),
		[2]cursor{},
	}
}

// draw a new point on the grid for the given wire by moving that wire's cursor in the given direction.
func (g *grid) draw(d *dir, wireNum int) {
	g.c[wireNum].p.x += d[0]
	g.c[wireNum].p.y += d[1]
	g.c[wireNum].dist++
	records, ok := g.g[g.c[wireNum].p]
	if ok {
		if !records[wireNum].visited {
			records[wireNum].visited = true
			records[wireNum].pathLength = g.c[wireNum].dist
			g.g[g.c[wireNum].p] = records
		}
	} else {
		records = [2]record{}
		records[wireNum].visited = true
		records[wireNum].pathLength = g.c[wireNum].dist
		g.g[g.c[wireNum].p] = records
	}
}

// draw a full wire on the grid based on the given move set.
func (g *grid) drawPath(moveSet []*dir, wireNum int) {
	for stepNum := range moveSet {
		g.draw(moveSet[stepNum], wireNum)
	}
}

// intersections finds all points in the grid where the wires crossed,
// and returns a slice of cursors giving the intersection points and
// the total distances traveled by the wires at the time they crossed.
func (g *grid) intersections() []cursor {
	shared := make([]cursor, 0, len(g.g))
	for p, records := range g.g {
		if records[0].visited && records[1].visited {
			shared = append(shared, cursor{p, records[0].pathLength + records[1].pathLength})
		}
	}
	return shared
}

// getMoves decomposes the input into slices of 1-step movements.
func getMoves(input string) [][]*dir {
	wires := strings.Split(input, "\n")
	moves := make([][]*dir, len(wires))

	for i, wire := range wires {
		vecs := strings.Split(wire, ",")

		for _, vec := range vecs {
			mag, _ := strconv.Atoi(vec[1:])
			d := dirs[vec[0]]
			stroke := make([]*dir, mag)

			for j := 0; j < mag; j++ {
				stroke[j] = d
			}

			moves[i] = append(moves[i], stroke...)
		}
	}
	return moves
}

// manhattan calculates the Manhattan distance between two points.
func manhattan(p1 point, p2 point) int {
	return int(
		math.Abs(float64(p1.x-p2.x)) +
			math.Abs(float64(p1.y-p2.y)),
	)
}

func part1(moves [][]*dir) (minDist int) {
	g := newGrid()
	for wireNum := range moves {
		g.drawPath(moves[wireNum], wireNum)
	}

	crosses := g.intersections()
	minDist = manhattan(origin, crosses[0].p)
	for _, c := range crosses[1:] {
		if dist := manhattan(origin, c.p); dist < minDist {
			minDist = dist
		}
	}
	return
}

func part2(moves [][]*dir) (minPathLength int) {
	g := newGrid()
	for wireNum := range moves {
		g.drawPath(moves[wireNum], wireNum)
	}

	crosses := g.intersections()
	minPathLength = crosses[0].dist
	for _, c := range crosses[1:] {
		if c.dist < minPathLength {
			minPathLength = c.dist
		}
	}
	return
}

// Solve provides the day 3 puzzle solution.
func Solve(input string) (*solutions.Solution, error) {
	moves := getMoves(input)
	return &solutions.Solution{Part1: part1(moves), Part2: part2(moves)}, nil
}
