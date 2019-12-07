package solutions

import (
	"fmt"
)

// Solution for parts 1 and 2 of a daily puzzle.
type Solution struct {
	Part1 interface{}
	Part2 interface{}
}

// Solver is a puzzle solver function type. Takes a puzzle input and returns a solution struct or an error.
type Solver func(string) (*Solution, error)

type registry map[string]Solver

func getKey(year int, day int) string {
	return fmt.Sprintf("%d-%02d", year, day)
}

// Register adds a new solver to the solution registry.
func (r registry) Register(year int, day int, f Solver) {
	r[getKey(year, day)] = f
}

// Get looks up a solver in the registry and returns it if found, as well as a bool indicating success/failure.
func (r registry) Get(year int, day int) (s Solver, ok bool) {
	s, ok = r[getKey(year, day)]
	return
}

// Registry of solver functions.
var Registry registry

func init() {
	Registry = make(map[string]Solver)
}
