// Package interpreter provides an intcode interpreter for day 2, 5, and other puzzles.
package interpreter

import (
	"fmt"

	"github.com/jzimbel/adventofcode-go/solutions/common"
)

type nothing struct{}

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
type procedure struct {
	f         func(i *Interpreter, args ...int)
	arity     uint
	writeArgs map[uint]nothing
}

// stores the information found in an opcode
type opDesc struct {
	proc *procedure
	m    [3]mode
}

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
	for {
		o = getOpDesc(i.codes[i.ipt])
		if o.proc.f == nil {
			// halt code encountered
			break
		}

		args := make([]int, o.proc.arity)
		var j uint
		var isWriteArg bool
		for j = 0; j < o.proc.arity; j++ {
			_, isWriteArg = o.proc.writeArgs[j]
			args[j] = i.getParam(j+1, o.m[j], isWriteArg)
		}

		o.proc.f(i, args...)
		i.ipt += o.proc.arity + 1
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

func add(i *Interpreter, args ...int) {
	i.codes[args[2]] = args[0] + args[1]
}

func mul(i *Interpreter, args ...int) {
	i.codes[args[2]] = args[0] * args[1]
}

func input(i *Interpreter, args ...int) {
	i.codes[args[0]] = i.input()
}

func output(i *Interpreter, args ...int) {
	i.output(args[0])
}

func init() {
	procedures = map[op]*procedure{
		oAdd: {
			f:         add,
			arity:     3,
			writeArgs: map[uint]nothing{2: nothing{}},
		},
		oMul: {
			f:         mul,
			arity:     3,
			writeArgs: map[uint]nothing{2: nothing{}},
		},
		oIn: {
			f:         input,
			arity:     1,
			writeArgs: map[uint]nothing{0: nothing{}},
		},
		oOut: {
			f:         output,
			arity:     1,
			writeArgs: map[uint]nothing{},
		},
		oHalt: {
			f:         nil,
			arity:     0,
			writeArgs: map[uint]nothing{},
		},
	}
}
