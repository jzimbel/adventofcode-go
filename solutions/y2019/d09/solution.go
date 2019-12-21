package d09

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jzimbel/adventofcode-go/solutions"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/interpreter"
)

func part1(codes []int) (result int) {
	input := func() int {
		return 1
	}
	output := func(n int) {
		result = n
		fmt.Println(n)
	}

	interpreter.New(codes, input, output).Run()
	return
}

func part2(codes []int) (result int) {
	input := func() int {
		return 2
	}
	output := func(n int) {
		result = n
		fmt.Println(n)
	}

	interpreter.New(codes, input, output).Run()
	return
}

// Solve provides the day 9 puzzle solution.
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

	return &solutions.Solution{Part1: part1(codes), Part2: part2(codes)}, nil
}
