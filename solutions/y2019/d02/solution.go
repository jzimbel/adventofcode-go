package d02

import (
	"github.com/jzimbel/adventofcode-go/solutions"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/interpreter"
)

func part1(initMem interpreter.Program) (int, error) {
	return interpreter.NewWithNounVerb(initMem, 12, 2, nil, nil).Run()
}

func part2(initMem interpreter.Program) (int, error) {
	for noun := 0; noun < 100; noun++ {
		for verb := 0; verb < 100; verb++ {
			result, err := interpreter.NewWithNounVerb(initMem, noun, verb, nil, nil).Run()
			if err != nil {
				return 0, err
			}
			if result == 19690720 {
				return 100*noun + verb, nil
			}
		}
	}
	return 0, nil
}

// Solve provides the day 2 puzzle solution.
func Solve(input string) (*solutions.Solution, error) {
	initMem := interpreter.ParseMem(input)

	answer1, err := part1(initMem)
	if err != nil {
		return nil, err
	}
	answer2, err := part2(initMem)
	if err != nil {
		return nil, err
	}
	return &solutions.Solution{Part1: answer1, Part2: answer2}, nil
}
