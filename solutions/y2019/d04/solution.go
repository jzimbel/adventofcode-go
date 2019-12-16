package d04

import (
	"regexp"
	"strconv"
	"sync"

	"github.com/jzimbel/adventofcode-go/solutions"
	"github.com/jzimbel/adventofcode-go/solutions/common"
)

type nothing struct{}

var (
	inputPattern = regexp.MustCompile(`^(\d+)-(\d+)$`)
)

func isValidPart1(n int) bool {
	var d int
	var repeat bool

	prev := common.GetDigit(n, 5)
	for i := 4; i >= 0; prev, i = d, i-1 {
		d = common.GetDigit(n, i)
		switch {
		case d < prev:
			return false
		case d == prev:
			repeat = true
		}
	}
	return repeat
}

func isValidPart2(n int) bool {
	var d int
	var repeat bool

	prev := common.GetDigit(n, 5)
	repeatCount := 1
	for i := 4; i >= 0; prev, i = d, i-1 {
		d = common.GetDigit(n, i)

		// digit must be >= prev
		if d < prev {
			return false
		}

		// search for a discrete 2-digit repeat if one hasn't already been found
		if !repeat {
			switch {
			case d == prev:
				repeatCount++
			case d != prev:
				if repeatCount == 2 {
					repeat = true
				}
				repeatCount = 1
			}
		}
	}
	// handles case where the 2-digit repeat comes at the end of the number
	if repeatCount == 2 {
		repeat = true
	}
	return repeat
}

// solve concurrently checks numbers in the range [lower,upper] using given predicate function pred, and returns the number that passed.
func solve(lower, upper int, pred func(int) bool) (validCount int) {
	// range is inclusive, so we add 1
	nRange := upper - lower + 1

	var wg sync.WaitGroup
	wg.Add(nRange)
	c := make(chan nothing, nRange)

	for i := lower; i <= upper; i++ {
		go func(icpy int) {
			if pred(icpy) {
				c <- nothing{}
			}
			wg.Done()
		}(i)
	}

	// close the channel once all goroutines are done
	go func() {
		wg.Wait()
		close(c)
	}()

	// increment validCount for each time a value was sent to c
	for range c {
		validCount++
	}
	return
}

// Solve provides the day 3 puzzle solution.
func Solve(input string) (*solutions.Solution, error) {
	bounds := inputPattern.FindStringSubmatch(input)[1:]
	lower, _ := strconv.Atoi(bounds[0])
	upper, _ := strconv.Atoi(bounds[1])
	return &solutions.Solution{
		Part1: solve(lower, upper, isValidPart1),
		Part2: solve(lower, upper, isValidPart2),
	}, nil
}
