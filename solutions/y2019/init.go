package y2019

import (
	"github.com/jzimbel/adventofcode-go/solutions"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/d01"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/d02"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/d03"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/d04"
)

func init() {
	r, y := &solutions.Registry, 2019
	r.Register(y, 1, d01.Solve)
	r.Register(y, 2, d02.Solve)
	r.Register(y, 3, d03.Solve)
	r.Register(y, 4, d04.Solve)
}
