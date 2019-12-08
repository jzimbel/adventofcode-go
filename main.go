package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jzimbel/adventofcode-go/color"
	"github.com/jzimbel/adventofcode-go/input"
	_ "github.com/jzimbel/adventofcode-go/solutioninit"
	"github.com/jzimbel/adventofcode-go/solutions"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <year> <day>\n", os.Args[0])
}

func getArgs() (int, int, bool) {
	args := os.Args[1:]
	if len(args) != 2 {
		return 0, 0, false
	}
	year, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Year argument must be an integer.")
		return 0, 0, false
	}
	day, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Day argument must be an integer.")
		return 0, 0, false
	}
	return year, day, true
}

func printSolution(s *solutions.Solution, year int, day int) {
	if s != nil {
		fmt.Println("Answer to part 1:", color.G(s.Part1))
		fmt.Println("Answer to part 2:", color.G(s.Part2))
	} else {
		fmt.Fprintf(os.Stderr, "No solution for year %d, day %d yet.\n", year, day)
	}
}

func main() {
	year, day, ok := getArgs()
	if !ok {
		usage()
		os.Exit(1)
	}
	input, err := input.Get(year, day)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load puzzle input: %v.\n", err)
		os.Exit(1)
	}
	solver, ok := solutions.Registry.Get(year, day)
	if !ok {
		fmt.Fprintf(os.Stderr, "Could not find solution code for year %d, day %d.\n", year, day)
		os.Exit(1)
	}
	s, err := solver(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	printSolution(s, year, day)
}
