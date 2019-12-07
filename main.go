package main

import (
	"fmt"
	"os"
	"strconv"

	// must come before solutions
	_ "github.com/jzimbel/adventofcode-go/solutioninit"

	"github.com/jzimbel/adventofcode-go/input"
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
	result, err := solver(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
	}
	fmt.Println("Solution:", result.Part1, result.Part2)
}
