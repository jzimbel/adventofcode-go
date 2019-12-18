// Package interpreter provides an intcode interpreter for day 2, 5, and other puzzles.
package interpreter

import (
	"fmt"

	"github.com/jzimbel/adventofcode-go/solutions/common"
)

// Interpreter executes a set of instructions.
type Interpreter struct {
	codes  []int
	ipt    uint
	input  func() int
	output func(int)
}

// parameter mode
type mode uint

// opcode
type op uint

// pairs a function that executes an instruction with its arity
// and indices of args (if any) that are used for writing values to memory
type procedure struct {
	f         func(i *Interpreter, args ...int) (skipIncrement bool)
	arity     uint
	writeArgs set
}

// stores the information found in an opcode
type opDesc struct {
	proc *procedure
	m    [3]mode
}

// Set type allows for efficient existence checking
type set map[uint]struct{}

const (
	mPosition mode = iota
	mImmediate
)

const (
	_ op = iota
	oAdd
	oMul
	oIn
	oOut
	oJumpIfNonZero
	oJumpIfZero
	oLessThan
	oEquals
	oHalt op = 99
)

var procedures map[op]*procedure

// New interpreter.
func New(codes []int, input func() int, output func(int)) *Interpreter {
	c := make([]int, len(codes))
	copy(c, codes)
	return &Interpreter{codes: c, input: input, output: output}
}

// NewWithNounVerb produces an interpreter with codes at indices 1 and 2 replaced by noun and verb.
func NewWithNounVerb(codes []int, noun int, verb int, input func() int, output func(int)) *Interpreter {
	c := make([]int, len(codes))
	copy(c, codes)
	c[1] = noun
	c[2] = verb
	return &Interpreter{codes: c, input: input, output: output}
}

// Run runs the interpreter until it encounters a halt opcode.
func (i *Interpreter) Run() (int, error) {
	var o *opDesc
	var skipIncrement bool
	for {
		o = getOpDesc(i.codes[i.ipt])
		if o.proc.f == nil {
			// halt code encountered
			break
		}

		args := make([]int, o.proc.arity)
		var j uint
		for j = 0; j < o.proc.arity; j++ {
			args[j] = i.getParam(j+1, o.m[j], o.proc.writeArgs.has(j))
		}

		skipIncrement = o.proc.f(i, args...)
		if !skipIncrement {
			i.ipt += o.proc.arity + 1
		}
	}
	return i.codes[0], nil
}

func (i *Interpreter) getParam(offset uint, m mode, isWriteArg bool) (v int) {
	if isWriteArg {
		v = i.codeAt(offset)
	} else {
		switch m {
		case mPosition:
			v = i.codes[i.codeAt(offset)]
		case mImmediate:
			v = i.codeAt(offset)
		default:
			panic(fmt.Sprintf("unknown parameter mode encountered: %d", m))
		}
	}
	return
}

func (i *Interpreter) codeAt(offset uint) int {
	return i.codes[i.ipt+offset]
}

func makeSet(values ...uint) (s set) {
	s = make(set, len(values))
	for _, value := range values {
		s[value] = struct{}{}
	}
	return
}

func (s set) has(value uint) (ok bool) {
	_, ok = s[value]
	return
}

func getOpDesc(o int) *opDesc {
	return &opDesc{
		proc: procedures[op(common.GetDigits(o, 0, 2))],
		m: [...]mode{
			mode(common.GetDigit(o, 2)),
			mode(common.GetDigit(o, 3)),
			mode(common.GetDigit(o, 4)),
		},
	}
}

func add(i *Interpreter, args ...int) (skipIncrement bool) {
	i.codes[args[2]] = args[0] + args[1]
	return
}

func mul(i *Interpreter, args ...int) (skipIncrement bool) {
	i.codes[args[2]] = args[0] * args[1]
	return
}

func input(i *Interpreter, args ...int) (skipIncrement bool) {
	i.codes[args[0]] = i.input()
	return
}

func output(i *Interpreter, args ...int) (skipIncrement bool) {
	i.output(args[0])
	return
}

func jumpIfNonZero(i *Interpreter, args ...int) (skipIncrement bool) {
	if args[0] != 0 {
		i.ipt = uint(args[1])
		skipIncrement = true
	}
	return
}

func jumpIfZero(i *Interpreter, args ...int) (skipIncrement bool) {
	if args[0] == 0 {
		i.ipt = uint(args[1])
		skipIncrement = true
	}
	return
}

func lessThan(i *Interpreter, args ...int) (skipIncrement bool) {
	if args[0] < args[1] {
		i.codes[args[2]] = 1
	} else {
		i.codes[args[2]] = 0
	}
	return
}

func equals(i *Interpreter, args ...int) (skipIncrement bool) {
	if args[0] == args[1] {
		i.codes[args[2]] = 1
	} else {
		i.codes[args[2]] = 0
	}
	return
}

func init() {
	procedures = map[op]*procedure{
		oAdd: {
			f:         add,
			arity:     3,
			writeArgs: makeSet(2),
		},
		oMul: {
			f:         mul,
			arity:     3,
			writeArgs: makeSet(2),
		},
		oIn: {
			f:         input,
			arity:     1,
			writeArgs: makeSet(0),
		},
		oOut: {
			f:         output,
			arity:     1,
			writeArgs: makeSet(),
		},
		oHalt: {
			f:         nil,
			arity:     0,
			writeArgs: makeSet(),
		},
		oJumpIfNonZero: {
			f:         jumpIfNonZero,
			arity:     2,
			writeArgs: makeSet(),
		},
		oJumpIfZero: {
			f:         jumpIfZero,
			arity:     2,
			writeArgs: makeSet(),
		},
		oLessThan: {
			f:         lessThan,
			arity:     3,
			writeArgs: makeSet(2),
		},
		oEquals: {
			f:         equals,
			arity:     3,
			writeArgs: makeSet(2),
		},
	}
}
