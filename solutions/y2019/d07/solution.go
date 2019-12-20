package d07

import (
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/jzimbel/adventofcode-go/solutions"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/interpreter"
)

const (
	ampCount          = 5
	initialInput      = 0
	minPhase     uint = 0
	maxPhase     uint = 4
)

var codes []int

func pow2(p uint) uint {
	return uint(math.Pow(2, float64(p)))
}

// phaseSettingsGenerator generates permutations of phase settings.
func phaseSettingsGenerator() <-chan [ampCount]uint {
	var uniqueMask uint
	for i := 0; i < ampCount; i++ {
		uniqueMask <<= 1
		uniqueMask++
	}

	var compare uint
	ch := make(chan [ampCount]uint)

	go func() {
		defer close(ch)
		// there are nice algorithms to do this, but I'm lazy
		for a := minPhase; a <= maxPhase; a++ {
			for b := minPhase; b <= maxPhase; b++ {
				for c := minPhase; c <= maxPhase; c++ {
					for d := minPhase; d <= maxPhase; d++ {
						for e := minPhase; e <= maxPhase; e++ {
							compare = pow2(a) | pow2(b) | pow2(c) | pow2(d) | pow2(e)
							if uniqueMask^compare == 0 {
								ch <- [...]uint{a, b, c, d, e}
							}
						}
					}
				}
			}
		}
	}()
	return ch
}

func makeInputDevice(phaseSetting uint, ch <-chan int) func() int {
	callCount := 0
	return func() int {
		defer func() { callCount++ }()
		switch callCount {
		case 0:
			return int(phaseSetting)
		case 1:
			return <-ch
		default:
			panic("Input called more times than expected")
		}
	}
}

func makeOutputDevice(ch chan<- int) func(int) {
	return func(n int) {
		ch <- n
	}
}

func runAmplifiers(settings [ampCount]uint) (signal int) {
	// 0 -> Amp A -> Amp B -> Amp C -> Amp D -> Amp E -> (to thrusters)
	// 5 amps, 6 channels
	chs := [ampCount + 1]chan int{}
	for i := range chs {
		chs[i] = make(chan int)
	}
	for i := 0; i < ampCount; i++ {
		go func(icpy int) {
			interpreter.New(codes, makeInputDevice(settings[icpy], chs[icpy]), makeOutputDevice(chs[icpy+1])).Run()
		}(i)
	}
	chs[0] <- initialInput
	return <-chs[ampCount]
}

func part1() (maxSignal int) {
	ch := make(chan int)
	wg := sync.WaitGroup{}
	for settings := range phaseSettingsGenerator() {
		wg.Add(1)
		go func(settings [ampCount]uint) {
			defer wg.Done()
			ch <- runAmplifiers(settings)
		}(settings)
	}
	go func() {
		defer close(ch)
		wg.Wait()
	}()

	for signal := range ch {

		if signal > maxSignal {
			maxSignal = signal
		}
	}
	return
}

// Solve provides the day 7 puzzle solution.
func Solve(input string) (*solutions.Solution, error) {
	numbers := strings.Split(input, ",")
	// this assigns the intcode source to a package-level variable for ease of access
	codes = make([]int, len(numbers))
	for i, n := range numbers {
		intn, err := strconv.Atoi(n)
		if err != nil {
			return nil, err
		}
		codes[i] = intn
	}
	return &solutions.Solution{Part1: part1()}, nil
}
