package d01

import (
	"strconv"
	"strings"
	"sync"

	"github.com/jzimbel/adventofcode-go/solutions"
)

func part1(masses []int) (sum int) {
	// unnecessary concurrency woo!
	c := make(chan int)
	defer close(c)
	calculate := func(n int) {
		c <- n/3 - 2
	}

	for i := range masses {
		go calculate(masses[i])
	}
	for range masses {
		sum += <-c
	}
	return
}

func part2(masses []int) (sum int) {
	// slightly less unnecessary concurrency?
	c := make(chan int)
	wg := sync.WaitGroup{}

	var calculate func(int)
	calculate = func(n int) {
		fuel := n/3 - 2
		if fuel > 0 {
			c <- fuel
			calculate(fuel)
		} else {
			wg.Done()
		}
	}

	for i := range masses {
		wg.Add(1)
		go calculate(masses[i])
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for n := range c {
		sum += n
	}
	return
}

// Solve provides the day 1 puzzle solution.
func Solve(input string) (*solutions.Solution, error) {
	lines := strings.Split(input, "\n")
	var masses []int
	for _, line := range lines {
		mass, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		masses = append(masses, mass)
	}
	return &solutions.Solution{Part1: part1(masses), Part2: part2(masses)}, nil
}
