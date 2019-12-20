package d07

import (
	"strconv"
	"strings"
	"sync"

	"github.com/jzimbel/adventofcode-go/solutions"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/interpreter"
	"modernc.org/mathutil"
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

func phaseSettingsGenerator(offset uint) <-chan *phaseSettings {
	ch := make(chan *phaseSettings)

	go func() {
		defer close(ch)
		ps := &phaseSettings{}
		for i := uint(0); i < ampCount; i++ {
			ps[i] = i + offset
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
	return func() (n int) {
		defer func() { callCount++ }()
		switch callCount {
		case 0:
			n = int(phaseSetting)
		default:
			n = <-ch
		}
		return
	}
}

func makeOutputDevice(ch chan<- int) func(int) {
	return func(n int) {
		ch <- n
	}
}

func makeLoopingOutputDevice(loop chan<- int, output chan<- int) func(int) {
	return func(n int) {
		loop <- n
		output <- n
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

func runAmplifierLoop(codes []int, settings *phaseSettings) (signal int) {
	// 0 -> Amp A -> Amp B -> Amp C -> Amp D -> Amp E -> (to thrusters upon Amp E halt)
	//    0        1        2        3        4        0
	// 5 amps, 5 channels
	chs := [ampCount]chan int{}
	for i := range chs {
		chs[i] = make(chan int)
	}
	// final amplifier also sends to this channel so that we can capture outputs
	output := make(chan int)

	wg := sync.WaitGroup{}
	for i := 0; i < ampCount; i++ {
		wg.Add(1)
		go func(icpy int) {
			defer func() {
				wg.Done()
			}()
			var outDevice func(int)
			if icpy == ampCount-1 {
				outDevice = makeLoopingOutputDevice(chs[(icpy+1)%ampCount], output)
			} else {
				outDevice = makeOutputDevice(chs[(icpy+1)%ampCount])
			}
			interpreter.New(codes, makeInputDevice(settings[icpy], chs[icpy]), outDevice).Run()
		}(i)
	}

	go func() {
		defer func() {
			for i := range chs {
				close(chs[i])
			}
			close(output)
		}()
		wg.Wait()
	}()

	chs[0] <- initialInput
	var finalSignal int
	for n := range output {
		finalSignal = n
	}
	return finalSignal
}

func part1(codes []int) (maxSignal int) {
	ch := make(chan int)
	wg := sync.WaitGroup{}
	for settings := range phaseSettingsGenerator(0) {
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

func part2(codes []int) (maxSignal int) {
	ch := make(chan int)
	wg := sync.WaitGroup{}
	for settings := range phaseSettingsGenerator(5) {
		wg.Add(1)
		go func(settings *phaseSettings) {
			defer wg.Done()
			ch <- runAmplifierLoop(codes, settings)
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
	return &solutions.Solution{ /*Part1: part1(codes),*/ Part2: part2(codes)}, nil
}
