// Package solutioninit imports all solution packages
// so that they can register themselves in the solution registry defined in registry.go.
package solutioninit

import (
	// register solutions for 2019 puzzles
	_ "github.com/jzimbel/adventofcode-go/solutions/y2019"
)
