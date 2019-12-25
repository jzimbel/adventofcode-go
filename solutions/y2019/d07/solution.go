package d07

import (
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

// Implements sort.Interface to take advantage of mathutil Permutation functions
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

// phaseSettingsGenerator returns a channel that receives all permutations of phase settings and then closes.
// ch1 := phaseSettingsGenerator(0)
// ch1 will receive [0 1 2 3 4], [0 1 2 4 3], ...
// ch2 := phaseSettingsGenerator(5)
// ch2 will receive [5 6 7 8 9], [5 6 7 9 8], ...
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

// makeInputDevice returns an input function to be used by the amplifier intcode program.
// The first time the input is called, it returns the phase setting for the amplifier.
// All future calls return values received from the given channel.
func makeInputDevice(phaseSetting uint, ch <-chan int) func() int {
	callCount := 0
	return func() (n int) {
		defer func() { callCount++ }()
		if callCount == 0 {
			n = int(phaseSetting)
		} else {
			n = <-ch
		}
		return
	}
}

// makeOutputDevice returns an output function to be used by the amplifier intcode program.
// This function sends its argument to the given channel.
func makeOutputDevice(ch chan<- int) func(int) {
	return func(n int) {
		ch <- n
	}
}

// makeLoopingOutputDevice is like makeOutputDevice, but the function it produces sends values to
// two given channels instead of one. This allows for signals to be passed in a loop but also received
// by an outside function that's interested in the final output of the amplifiers.
func makeLoopingOutputDevice(loop chan<- int, output chan<- int) func(int) {
	return func(n int) {
		loop <- n
		output <- n
	}
}

// runAmplifiers runs a series of amplifiers with the given phase settings and returns their output.
func runAmplifiers(initMem interpreter.Program, settings *phaseSettings) (signal int) {
	// 0 -> Amp A -> Amp B -> Amp C -> Amp D -> Amp E -> (to thrusters)
	// 5 amps, 6 channels
	chs := [ampCount + 1]chan int{}
	for i := range chs {
		chs[i] = make(chan int)
	}

	for i := 0; i < ampCount; i++ {
		go func(icpy int) {
			interpreter.New(initMem, makeInputDevice(settings[icpy], chs[icpy]), makeOutputDevice(chs[icpy+1])).Run()
		}(i)
	}

	chs[0] <- initialInput
	return <-chs[ampCount]
}

// runAmplifierLoop runs a series of amplifiers in a loop with the given phase settings and returns their final output when they halt.
func runAmplifierLoop(initMem interpreter.Program, settings *phaseSettings) (signal int) {
	// 0 -> Amp A -> Amp B -> Amp C -> Amp D -> Amp E -> (to thrusters upon Amp E halt)
	//    0        1        2        3        4        0
	// 5 amps, 5 channels
	chs := [ampCount]chan int{}
	// use a buffered channel for the first channel so that it can receive
	// one more value that won't be consumed during the final iteration of the loop
	chs[0] = make(chan int, 1)
	for i := range chs[1:] {
		chs[i+1] = make(chan int)
	}
	// final amplifier also sends to this channel so that we can capture outputs
	output := make(chan int)

	for i := 0; i < ampCount; i++ {
		go func(icpy int) {
			var outDevice func(int)
			if icpy == ampCount-1 {
				// when this interpreter halts, the whole amplifier loop is done
				defer close(output)
				outDevice = makeLoopingOutputDevice(chs[(icpy+1)%ampCount], output)
			} else {
				outDevice = makeOutputDevice(chs[(icpy+1)%ampCount])
			}
			interpreter.New(initMem, makeInputDevice(settings[icpy], chs[icpy]), outDevice).Run()
		}(i)
	}

	chs[0] <- initialInput
	var finalSignal int
	for n := range output {
		finalSignal = n
	}
	return finalSignal
}

func run(initMem interpreter.Program, phaseSettingOffset uint, runner func(interpreter.Program, *phaseSettings) int) (maxSignal int) {
	ch := make(chan int)
	wg := sync.WaitGroup{}

	for settings := range phaseSettingsGenerator(phaseSettingOffset) {
		wg.Add(1)
		go func(settings *phaseSettings) {
			defer wg.Done()
			ch <- runner(initMem, settings)
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

func part1(initMem interpreter.Program, ch chan<- int) {
	ch <- run(initMem, 0, runAmplifiers)
}

func part2(initMem interpreter.Program, ch chan<- int) {
	ch <- run(initMem, 5, runAmplifierLoop)
}

// Solve provides the day 7 puzzle solution.
func Solve(input string) (*solutions.Solution, error) {
	initMem := interpreter.ParseMem(input)
	ch1, ch2 := make(chan int), make(chan int)

	go part1(initMem, ch1)
	go part2(initMem, ch2)
	return &solutions.Solution{Part1: <-ch1, Part2: <-ch2}, nil
}
