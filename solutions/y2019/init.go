package y2019

import (
	"github.com/jzimbel/adventofcode-go/solutions"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/d01"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/d02"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/d03"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/d04"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/d05"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/d06"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/d07"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/d08"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/d09"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/d10"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/d11"
)

func init() {
	r, y := &solutions.Registry, 2019
	r.Register(y, 1, d01.Solve)
	r.Register(y, 2, d02.Solve)
	r.Register(y, 3, d03.Solve)
	r.Register(y, 4, d04.Solve)
	r.Register(y, 5, d05.Solve)
	r.Register(y, 6, d06.Solve)
	r.Register(y, 7, d07.Solve)
	r.Register(y, 8, d08.Solve)
	r.Register(y, 9, d09.Solve)
	r.Register(y, 10, d10.Solve)
	r.Register(y, 11, d11.Solve)
}
