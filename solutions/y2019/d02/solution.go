package d02

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jzimbel/adventofcode-go/solutions"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/interpreter"
)

// original intcode interpreter, now replaced by the shared one
func compute(codes []int, noun int, verb int) (int, error) {
	codes[1] = noun
	codes[2] = verb
	ipt := 0

Run:
	for {
		switch codes[ipt] {
		case 1:
			codes[codes[ipt+3]] = codes[codes[ipt+1]] + codes[codes[ipt+2]]
		case 2:
			codes[codes[ipt+3]] = codes[codes[ipt+1]] * codes[codes[ipt+2]]
		case 99:
			break Run
		default:
			return 0, fmt.Errorf("unknown opcode encountered: %d", codes[ipt])
		}
		ipt += 4
	}
	return codes[0], nil
}

func part1(codes []int) (int, error) {
	return interpreter.NewWithNounVerb(codes, 12, 2, nil, nil).Run()
}

func part2(codes []int) (int, error) {
	for noun := 0; noun < 100; noun++ {
		for verb := 0; verb < 100; verb++ {
			result, err := interpreter.NewWithNounVerb(codes, noun, verb, nil, nil).Run()
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
	numbers := strings.Split(input, ",")
	codes := make([]int, len(numbers))
	for i, n := range numbers {
		intn, err := strconv.Atoi(n)
		if err != nil {
			return nil, err
		}
		codes[i] = intn
	}

	answer1, err := part1(codes)
	if err != nil {
		return nil, err
	}
	answer2, err := part2(codes)
	if err != nil {
		return nil, err
	}
	return &solutions.Solution{Part1: answer1, Part2: answer2}, nil
}
