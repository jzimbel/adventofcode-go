package d12

import (
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/jzimbel/adventofcode-go/solutions"
)

const steps = 1000

type point [3]int

type vec struct {
	pos point
	vel point
}

type system []vec

var pattern = regexp.MustCompile(`^<x=(-?\d+), y=(-?\d+), z=(-?\d+)>$`)

func (p *point) energy() (e int) {
	for _, v := range p {
		e += int(math.Abs(float64(v)))
	}
	return
}

func (v *vec) attract(v2 *vec) {
	for i := range v.pos {
		switch {
		case v.pos[i] < v2.pos[i]:
			v.vel[i]++
		case v.pos[i] > v2.pos[i]:
			v.vel[i]--
		}
	}
}

// attract simulates attracting v to the other vecs in its system.
func (v *vec) attractAll(s system) {
	for i := range s {
		v.attract(&s[i])
	}
}

func (v *vec) move() {
	for i := range v.pos {
		v.pos[i] += v.vel[i]
	}
}

func (v *vec) energy() (e int) {
	return v.pos.energy() * v.vel.energy()
}

func (s system) gravitate() {
	for i := range s {
		s[i].attractAll(s)
	}
}

func (s system) translate() {
	for i := range s {
		s[i].move()
	}
}

func (s system) energy() (e int) {
	for i := range s {
		e += s[i].energy()
	}
	return
}

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

func part1(initS system) int {
	s := make(system, len(initS))
	copy(s, initS)

	for i := 0; i < steps; i++ {
		s.gravitate()
		s.translate()
	}
	return s.energy()
}

// Solve provides the day 12 puzzle solution.
func Solve(input string) (*solutions.Solution, error) {
	s := parse(input)
	return &solutions.Solution{Part1: part1(s), Part2: nil}, nil
}
