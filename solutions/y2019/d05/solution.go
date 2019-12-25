package d05

import (
	"fmt"

	"github.com/jzimbel/adventofcode-go/solutions"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/interpreter"
)

func run(initMem interpreter.Program, systemID int) (int, error) {
	var lastOutput int

	_, err := interpreter.New(
		initMem,
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
	initMem := interpreter.ParseMem(input)

	answer1, err := run(initMem, 1)
	if err != nil {
		return nil, err
	}
	answer2, err := run(initMem, 5)
	if err != nil {
		return nil, err
	}

	return &solutions.Solution{Part1: answer1, Part2: answer2}, nil
}
