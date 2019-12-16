package d05

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jzimbel/adventofcode-go/color"
	"github.com/jzimbel/adventofcode-go/solutions"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/interpreter"
)

func input() (v int) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(color.R("> "))
		b, err := reader.ReadBytes('\n')
		s := string(bytes.TrimSpace(b))
		if err != nil {
			fmt.Println(color.R("Rejected"))
			continue
		}
		v, err = strconv.Atoi(s)
		if err != nil {
			fmt.Println(color.R("Rejected"))
			continue
		}
		fmt.Println(color.G("Accepted"))
		break
	}
	return
}

func output(n int) {
	fmt.Println(n)
}

func part1(codes []int) (int, error) {
	var lastOutput int
	myOutput := func(n int) {
		lastOutput = n
		output(n)
	}

	_, err := interpreter.New(codes, input, myOutput).Run()
	if err != nil {
		return 0, err
	}
	return lastOutput, nil
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

	return &solutions.Solution{Part1: answer1}, nil
}
