package y2019

import (
	"github.com/jzimbel/adventofcode-go/solutions"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/d01"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/d02"
)

func init() {
	r, y := &solutions.Registry, 2019
	r.Register(y, 1, d01.Solve)
	r.Register(y, 2, d02.Solve)
}
