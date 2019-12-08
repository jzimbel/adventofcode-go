package d02

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jzimbel/adventofcode-go/solutions"
)

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
	return compute(codes, 12, 2)
}

func part2(codesImmutable []int) (int, error) {
	codes := make([]int, len(codesImmutable))
	for noun := 0; noun < 100; noun++ {
		for verb := 0; verb < 100; verb++ {
			copy(codes, codesImmutable)
			result, err := compute(codes, noun, verb)
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
	codes1 := make([]int, len(numbers))
	codes2 := make([]int, len(numbers))
	for i, n := range numbers {
		intn, err := strconv.Atoi(n)
		if err != nil {
			return nil, err
		}
		codes1[i] = intn
	}
	copy(codes2, codes1)

	answer1, err := part1(codes1)
	if err != nil {
		return nil, err
	}
	answer2, err := part2(codes2)
	if err != nil {
		return nil, err
	}
	return &solutions.Solution{Part1: answer1, Part2: answer2}, nil
}
