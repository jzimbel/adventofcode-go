package d05

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jzimbel/adventofcode-go/solutions"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/interpreter"
)

func run(codes []int, systemID int) (int, error) {
	var lastOutput int

	_, err := interpreter.New(
		codes,
		func() int {
			return systemID
		},
		func(n int) {
			lastOutput = n
			fmt.Println(n)
		},
	).Run()
	if err != nil {
		return 0, err
	}
	return lastOutput, nil
}

// Solve provides the day 5 puzzle solution.
func Solve(input string) (*solutions.Solution, error) {
	numbers := strings.Split(input, ",")
	codes := make([]int, len(numbers))
	for i, n := range numbers {
		intn, err := strconv.Atoi(n)
		if err != nil {
			return nil, err
		}
		codes[i] = intn
	}

	answer1, err := run(codes, 1)
	if err != nil {
		return nil, err
	}
	answer2, err := run(codes, 5)
	if err != nil {
		return nil, err
	}

	return &solutions.Solution{Part1: answer1, Part2: answer2}, nil
}
