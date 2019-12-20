package d07

import (
	"strconv"
	"strings"
	"sync"

	"modernc.org/mathutil"
	"github.com/jzimbel/adventofcode-go/solutions"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/interpreter"
)

const (
	ampCount          = 5
	initialInput      = 0
	minPhase     uint = 0
	maxPhase     uint = 4
)

// Implements sort.Interface to take advantage of Permutation functions
type phaseSettings [ampCount]uint

func (ps *phaseSettings) Len() int {
	return len(ps)
}

func (ps *phaseSettings) Less(i, j int) bool {
	return ps[i] < ps[j]
}

func (ps *phaseSettings) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func phaseSettingsGenerator() <-chan *phaseSettings {
	ch := make(chan *phaseSettings)

	go func() {
		defer close(ch)
		ps := &phaseSettings{}
		for i := uint(0); i < ampCount; i++ {
			ps[i] = i
		}
		mathutil.PermutationFirst(ps)

		var done bool
		for !done {
			psCopy := *ps
			ch <- &psCopy
			done = !mathutil.PermutationNext(ps)
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

func runAmplifiers(codes []int, settings *phaseSettings) (signal int) {
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

func part1(codes []int) (maxSignal int) {
	ch := make(chan int)
	wg := sync.WaitGroup{}
	for settings := range phaseSettingsGenerator() {
		wg.Add(1)
		go func(settings *phaseSettings) {
			defer wg.Done()
			ch <- runAmplifiers(codes, settings)
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
	codes := make([]int, len(numbers))
	for i, n := range numbers {
		intn, err := strconv.Atoi(n)
		if err != nil {
			return nil, err
		}
		codes[i] = intn
	}
	return &solutions.Solution{Part1: part1(codes)}, nil
}
