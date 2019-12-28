package d13

import (
	"math"

	"github.com/jzimbel/adventofcode-go/solutions"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/interpreter"
)

const (
	outputTypeX = iota
	outputTypeY
	outputTypeDraw
)

const (
	_ = iota
	_
	block
	paddle
	ball
)

type point struct {
	x int
	y int
}

type grid map[point]int

func part1(initMem interpreter.Program) int {
	g := make(grid)
	var x, y, outputType int

	output := func(n int) {
		switch outputType {
		case outputTypeX:
			x = n
		case outputTypeY:
			y = n
		case outputTypeDraw:
			g[point{x, y}] = n
		}
		outputType = (outputType + 1) % 3
	}

	interpreter.New(initMem, nil, output).Run()

	var blockCount int
	for _, tile := range g {
		if tile == block {
			blockCount++
		}
	}
	return blockCount
}

func part2(initMem interpreter.Program) int {
	initMem[0] = 2
	g := make(grid)
	var x, y, outputType, ballX, paddleX, score int

	input := func() int {
		diff := ballX - paddleX
		if diff == 0 {
			return diff
		}
		return diff / int(math.Abs(float64(diff)))
	}

	output := func(n int) {
		switch outputType {
		case outputTypeX:
			x = n
		case outputTypeY:
			y = n
		case outputTypeDraw:
			if x == -1 && y == 0 {
				score = n
			} else {
				g[point{x, y}] = n
				switch n {
				case ball:
					ballX = x
				case paddle:
					paddleX = x
				}
			}
		}
		outputType = (outputType + 1) % 3
	}

	interpreter.New(initMem, input, output).Run()
	return score
}

// Solve provides the day 13 puzzle solution.
func Solve(input string) (*solutions.Solution, error) {
	initMem := interpreter.ParseMem(input)
	return &solutions.Solution{Part1: part1(initMem), Part2: part2(initMem)}, nil
}
