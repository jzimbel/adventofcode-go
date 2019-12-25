// Package interpreter provides an intcode interpreter for day 2, 5, and other puzzles.
package interpreter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jzimbel/adventofcode-go/solutions/common"
)

// Program code / memory of a running Program.
type Program map[uint]int

// Interpreter executes a set of instructions.
// Instantiate with one of the `New*` functions.
type Interpreter struct {
	mem    Program
	ipt    uint
	rpt    uint
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
	mRelative
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
	oMoveRpt
	oHalt op = 99
)

var procedures map[op]*procedure

// New interpreter.
func New(initMem Program, input func() int, output func(int)) (i *Interpreter) {
	mem := make(Program, len(initMem))
	for i, n := range initMem {
		mem[i] = n
	}
	return &Interpreter{mem: mem, input: input, output: output}
}

// NewWithNounVerb produces an interpreter with memory at addresses 1 and 2 replaced by noun and verb.
func NewWithNounVerb(initMem Program, noun int, verb int, input func() int, output func(int)) *Interpreter {
	mem := make(Program, len(initMem))
	for i, n := range initMem {
		mem[i] = n
	}
	mem[1] = noun
	mem[2] = verb
	return &Interpreter{mem: mem, input: input, output: output}
}

// ParseMem parses the initial memory/instructions of an intcode program from a puzzle input.
func ParseMem(input string) (mem Program) {
	numbers := strings.Split(input, ",")
	mem = make(Program, len(numbers))
	for i, n := range numbers {
		intn, _ := strconv.Atoi(n)
		mem[uint(i)] = intn
	}
	return
}

// Run runs the interpreter until it encounters a halt opcode.
func (i *Interpreter) Run() (int, error) {
	var o *opDesc
	var skipIncrement bool
	for {
		o = getOpDesc(i.get(i.ipt))
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
	return i.get(0), nil
}

func (i *Interpreter) get(addr uint) (val int) {
	val, ok := i.mem[addr]
	if !ok {
		val = 0
	}
	return
}

func (i *Interpreter) set(addr uint, val int) {
	i.mem[addr] = val
}

func (i *Interpreter) getParam(offset uint, m mode, isWriteArg bool) (v int) {
	if isWriteArg {
		switch m {
		case mPosition:
			v = i.codeAt(offset)
		case mRelative:
			v = int(i.rpt) + i.codeAt(offset)
		default:
			panic("immediate mode used for write arg")
		}
	} else {
		switch m {
		case mPosition:
			v = i.get(uint(i.codeAt(offset)))
		case mImmediate:
			v = i.codeAt(offset)
		case mRelative:
			v = i.get(i.rpt + uint(i.codeAt(offset)))
		default:
			panic(fmt.Sprintf("unknown parameter mode encountered: %d", m))
		}
	}
	return
}

func (i *Interpreter) codeAt(offset uint) int {
	return i.get(i.ipt + offset)
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
	i.set(uint(args[2]), args[0]+args[1])
	return
}

func mul(i *Interpreter, args ...int) (skipIncrement bool) {
	i.set(uint(args[2]), args[0]*args[1])
	return
}

func input(i *Interpreter, args ...int) (skipIncrement bool) {
	i.set(uint(args[0]), i.input())
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
		i.set(uint(args[2]), 1)
	} else {
		i.set(uint(args[2]), 0)
	}
	return
}

func equals(i *Interpreter, args ...int) (skipIncrement bool) {
	if args[0] == args[1] {
		i.set(uint(args[2]), 1)
	} else {
		i.set(uint(args[2]), 0)
	}
	return
}

func moveRpt(i *Interpreter, args ...int) (skipIncrement bool) {
	// might underflow if arg is negative, but we're going
	// to trust that none of the intcode programs will do that
	i.rpt = uint(int(i.rpt) + args[0])
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
		oMoveRpt: {
			f:         moveRpt,
			arity:     1,
			writeArgs: makeSet(),
		},
		oHalt: {
			f:         nil,
			arity:     0,
			writeArgs: makeSet(),
		},
	}
}
