package d10

import (
	"math"
	"sort"
	"strings"

	"github.com/jzimbel/adventofcode-go/solutions"
)

const (
	// width and height of my puzzle input, used for some slight optimizations
	width  = 24
	height = 24
)

var epsilon float64

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

// axisDistances returns (separately) the x and y distances between two points.
func axisDistances(p1, p2 *point) (xDist int, yDist int) {
	xDist, yDist = int(math.Abs(float64(p1.x-p2.x))), int(math.Abs(float64(p1.y-p2.y)))
	return
}

func isBlocked(p1, p2 *point, g grid, gcd func(int, int) int) (blocked bool) {
	denom := gcd(axisDistances(p1, p2))
	xStepSize, yStepSize := (p2.x-p1.x)/denom, (p2.y-p1.y)/denom
	for i := 1; i < denom; i++ {
		if _, ok := g[point{p1.x + i*xStepSize, p1.y + i*yStepSize}]; ok {
			blocked = true
			break
		}
	}
	return
}

func part1(g grid, gcd func(int, int) int) (maxVisibleCount int, optimalPoint point) {
	for p1 := range g {
		var visibleCount int
		for p2 := range g {
			if p1 == p2 {
				continue
			}
			if !isBlocked(&p1, &p2, g, gcd) {
				visibleCount++
			}
		}
		if visibleCount > maxVisibleCount {
			maxVisibleCount = visibleCount
			optimalPoint = p1
		}
	}
	return
}

// getCenteredGrid creates and returns a new grid with all
// points translated such that p at the origin.
func getCenteredGrid(g grid, p point) (cg grid) {
	cg = make(grid, len(g))
	for p2 := range g {
		cg[point{p2.x - p.x, p2.y - p.y}] = struct{}{}
	}
	return cg
}

type rPoint struct {
	r    float64
	θ    float64 // #codegolfing
	orig point
}

// rPoints is an ordered list of radial (or polar) coordinates.
// Polar axis (where θ = 0) is up.
type rPoints []rPoint

func (rg rPoints) Len() int {
	return len(rg)
}

func (rg rPoints) Less(i, j int) bool {
	if math.Abs(rg[i].θ-rg[j].θ) < epsilon {
		return rg[i].r < rg[j].r
	}
	return rg[i].θ < rg[j].θ
}

func (rg rPoints) Swap(i, j int) {
	rg[i], rg[j] = rg[j], rg[i]
}

// dist returns the Euclidean distance between two points. ( sqrt(a**2 + b**2) )
func dist(p1, p2 *point) float64 {
	return math.Sqrt(math.Pow(math.Abs(float64(p1.x-p2.x)), 2) + math.Pow(math.Abs(float64(p1.y-p2.y)), 2))
}

// clockwiseAngleFromUp calculates the angle from -y in radians, moving clockwise.
// Atan2 normally takes arguments as (y,x), but we reverse them and negate x in order to
// have polar θ = 0 be Cartesian (0,-1) and increasing θ move in a clockwise direction.
// Atan2 also produces values in range [-π, π], but we want them to be [0, 2π],
// so when it would normally produce a negative, we use 2π + atan2Result.
func clockwiseAngleFromUp(p *point) (θ float64) {
	newX, newY := -p.y, p.x
	θ = math.Atan2(float64(newY), float64(newX))
	if θ < 0 {
		θ = 2*math.Pi + θ
	}
	return
}

func part2(cg grid, gcd func(int, int) int, optimalPoint point) int {
	// record angle and distance from center of each asteroid in a sorted slice of struct {rad float64; dist float64}
	rp := make(rPoints, 0, len(cg))
	origin := point{}
	for p := range cg {
		if p == origin {
			continue
		}
		rp = append(rp, rPoint{r: dist(&origin, &p), θ: clockwiseAngleFromUp(&p), orig: p})
	}
	sort.Sort(rp)

	var vaporizedCount int
	for len(rp) > 0 {
		nextRp := make(rPoints, 0, len(rp))
		remove := make([]*point, 0, len(rp))
		for i := range rp {
			if isBlocked(&origin, &rp[i].orig, cg, gcd) {
				nextRp = append(nextRp, rp[i])
			} else {
				vaporizedCount++
				if vaporizedCount == 200 {
					return (optimalPoint.x+rp[i].orig.x)*100 + (optimalPoint.y + rp[i].orig.y)
				}
				remove = append(remove, &rp[i].orig)
			}
		}

		for i := range remove {
			delete(cg, *remove[i])
		}
		rp = nextRp
	}
	// unreachable as long as there are at least 200 asteroids on the grid
	return 0
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

	maxVisibleCount, optimalPoint := part1(g, gcd)

	// now center the grid on the optimal point for part 2
	g = getCenteredGrid(g, optimalPoint)
	return &solutions.Solution{Part1: maxVisibleCount, Part2: part2(g, gcd, optimalPoint)}, nil
}

func init() {
	epsilon = math.Nextafter(1, 2) - 1
}
