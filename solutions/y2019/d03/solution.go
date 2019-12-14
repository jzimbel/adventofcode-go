package d03

import (
	"math"
	"strconv"
	"strings"

	"github.com/jzimbel/adventofcode-go/solutions"
)

var (
	origin = point{}
	dirs   = map[byte]*dir{
		'D': &dir{0, 1},
		'L': &dir{-1, 0},
		'R': &dir{1, 0},
		'U': &dir{0, -1},
	}
)

type point struct {
	x int
	y int
}

// dir represents a unit vector in one of the cardinal directions.
type dir point

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

// move a cursor 1 unit in the given direction.
func (c *cursor) move(d *dir) {
	c.p.x += d.x
	c.p.y += d.y
	c.dist++
}

// grid holds a map that records which points have been visited by which paths,
// as well as cursors for the two paths.
type grid struct {
	g map[point][2]record
	c [2]cursor
}

func newGrid() *grid {
	return &grid{
		make(map[point][2]record),
		[2]cursor{},
	}
}

func (g *grid) moveCursor(wireNum int, d *dir) {
	g.c[wireNum].move(d)
}

func (g *grid) getRecordsAtCursor(wireNum int) (r [2]record, ok bool) {
	r, ok = g.g[g.c[wireNum].p]
	return
}

func (g *grid) setRecordsAtCursor(wireNum int, r [2]record) {
	g.g[g.c[wireNum].p] = r
}

// draw a new point on the grid for the given wire by moving that wire's cursor in the given direction.
func (g *grid) draw(d *dir, wireNum int) {
	g.moveCursor(wireNum, d)
	records, ok := g.getRecordsAtCursor(wireNum)
	if !ok {
		records = [2]record{}
	}
	records[wireNum].visited = true
	records[wireNum].pathLength = g.c[wireNum].dist
	g.setRecordsAtCursor(wireNum, records)
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
func getMoves(input string) (moves [2][]*dir) {
	wires := strings.Split(input, "\n")

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

// solve finds the answers to parts 1 and 2 simultaneously.
// minDist = part 1 solution
// minPathLength = part 2 solution
func solve(moves [2][]*dir) (minDist, minPathLength int) {
	g := newGrid()
	for wireNum := range moves {
		g.drawPath(moves[wireNum], wireNum)
	}

	crosses := g.intersections()
	minDist = manhattan(origin, crosses[0].p)
	minPathLength = crosses[0].dist
	for _, c := range crosses[1:] {
		if dist := manhattan(origin, c.p); dist < minDist {
			minDist = dist
		}
		if c.dist < minPathLength {
			minPathLength = c.dist
		}
	}
	return
}

// Solve provides the day 3 puzzle solution.
func Solve(input string) (*solutions.Solution, error) {
	minDist, minPathLength := solve(getMoves(input))
	return &solutions.Solution{Part1: minDist, Part2: minPathLength}, nil
}
