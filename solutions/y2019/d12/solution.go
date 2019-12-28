package d12

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/jzimbel/adventofcode-go/solutions"
)

const (
	axisCount = 3
	steps     = 1000
)

type point [3]int
type vec struct {
	pos point
	vel point
}
type system []vec

type axisVec struct {
	pPos *int
	pVel *int
}
type axis []axisVec
type axisSystem [axisCount]axis

func newAxisSystem(s system) *axisSystem {
	axs := axisSystem{}
	for i := range axs {
		axs[i] = make(axis, len(s))
	}

	for i := range s {
		for j := 0; j < axisCount; j++ {
			axs[j][i].pPos = &s[i].pos[j]
			axs[j][i].pVel = &s[i].vel[j]
		}
	}
	return &axs
}

func (av *axisVec) attract(av2 *axisVec) {
	switch {
	case *av.pPos < *av2.pPos:
		*av.pVel++
	case *av.pPos > *av2.pPos:
		*av.pVel--
	}
}

func (av *axisVec) move() {
	*av.pPos += *av.pVel
}

func (av *axisVec) attractAll(a axis) {
	for i := range a {
		av.attract(&a[i])
	}
}

func (av *axisVec) equals(av2 *axisVec) bool {
	return *av.pPos == *av2.pPos && *av.pVel == *av2.pVel
}

func (a axis) gravitate() {
	for i := range a {
		a[i].attractAll(a)
	}
}

func (a axis) move() {
	for i := range a {
		a[i].move()
	}
}

func (a axis) equals(a2 axis) (eq bool) {
	eq = true
	for i := range a {
		if eq = a[i].equals(&a2[i]); !eq {
			break
		}
	}
	return
}

func (a axis) simulate() {
	for j := 0; j < steps; j++ {
		a.gravitate()
		a.move()
	}
}

func (a axis) copy() (acpy axis) {
	acpy = make(axis, len(a))
	for i := range a {
		pos, vel := *a[i].pPos, *a[i].pVel
		acpy[i] = axisVec{pPos: &pos, pVel: &vel}
	}
	return
}

func (a axis) findRepeat() (count int) {
	ref := a.copy()
	for done := false; !done; done = a.equals(ref) {
		count++
		a.gravitate()
		a.move()
	}
	return
}

func (axs *axisSystem) simulate() {
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wg.Add(axisCount)
	for i := range axs {
		go func(icpy int) {
			defer wg.Done()
			axs[icpy].simulate()
		}(i)
	}
}

func (axs *axisSystem) findRepeat() uint64 {
	wg := sync.WaitGroup{}
	wg.Add(axisCount)
	ch := make(chan uint64)
	for i := range axs {
		go func(icpy int) {
			defer wg.Done()
			ch <- uint64(axs[icpy].findRepeat())
		}(i)
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	results := make([]uint64, 0, axisCount)
	for n := range ch {
		results = append(results, n)
	}
	return lcm(results[0], results[1], results[2:]...)
}

func (p *point) energy() (e int) {
	for _, v := range p {
		e += int(math.Abs(float64(v)))
	}
	return
}

func (v *vec) energy() (e int) {
	return v.pos.energy() * v.vel.energy()
}

func (s system) energy() (e int) {
	for i := range s {
		e += s[i].energy()
	}
	return
}

var pattern = regexp.MustCompile(`^<x=(-?\d+), y=(-?\d+), z=(-?\d+)>$`)

func parse(input string) (s system) {
	lines := strings.Split(input, "\n")
	s = make(system, len(lines))
	for i := range lines {
		matches := pattern.FindStringSubmatch(lines[i])[1:]
		intMatches := make([]int, len(matches))
		for j := range matches {
			n, _ := strconv.Atoi(matches[j])
			intMatches[j] = n
		}
		s[i] = vec{pos: point{intMatches[0], intMatches[1], intMatches[2]}}
	}
	return
}

func gcd(a, b uint64) (n uint64) {
	if b != 0 {
		n = gcd(b, a%b)
	} else {
		n = a
	}
	return
}

func lcm(a, b uint64, integers ...uint64) uint64 {
	result := a * b / gcd(a, b)
	for i := 0; i < len(integers); i++ {
		result = lcm(result, integers[i])
	}
	return result
}

func part1(initS system) int {
	s := make(system, len(initS))
	copy(s, initS)
	axs := newAxisSystem(s)
	axs.simulate()
	return s.energy()
}

func part2(initS system) uint64 {
	s := make(system, len(initS))
	copy(s, initS)
	axs := newAxisSystem(s)
	return axs.findRepeat()
}

// Solve provides the day 12 puzzle solution.
func Solve(input string) (*solutions.Solution, error) {
	s := parse(input)
	return &solutions.Solution{Part1: part1(s), Part2: part2(s)}, nil
}
