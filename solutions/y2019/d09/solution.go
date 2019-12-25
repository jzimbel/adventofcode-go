package d09

import (
	"fmt"

	"github.com/jzimbel/adventofcode-go/solutions"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/interpreter"
)

func part1(initMem interpreter.Program) (result int) {
	input := func() int {
		return 1
	}
	output := func(n int) {
		result = n
		fmt.Println(n)
	}

	interpreter.New(initMem, input, output).Run()
	return
}

func part2(initMem interpreter.Program) (result int) {
	input := func() int {
		return 2
	}
	output := func(n int) {
		result = n
		fmt.Println(n)
	}

	interpreter.New(initMem, input, output).Run()
	return
}

// Solve provides the day 9 puzzle solution.
func Solve(input string) (*solutions.Solution, error) {
	initMem := interpreter.ParseMem(input)

	return &solutions.Solution{Part1: part1(initMem), Part2: part2(initMem)}, nil
}
