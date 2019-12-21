package d08

import (
	"strings"

	"github.com/jzimbel/adventofcode-go/solutions"
)

const (
	width     = 25
	height    = 6
	layerSize = width * height
	// difference between byte values for e.g. '0' and 0, '1' and 1, etc.
	asciiDigitDiff = 48
)

const (
	black uint8 = iota
	white
	transparent
)

var colorMap = map[uint8]rune{
	black:       '\u25A0',
	white:       '\u25A1',
	transparent: ' ',
}

type row [width]uint8
type layer [height]row
type image []layer

func (r *row) String() string {
	builder := make([]rune, len(r))
	for x := range r {
		builder[x] = colorMap[r[x]]
	}
	return string(builder)
}

func (l *layer) String() string {
	rows := make([]string, len(l))
	for y := range l {
		rows[y] = l[y].String()
	}
	return "\n" + strings.Join(rows, "\n")
}

// mergeDown merges l onto l2, replacing any transparent pixels in l with
// non-transparent ones in the same location of l2.
// Only l is mutated by this function.
func (l *layer) mergeDown(l2 *layer) {
	for y := range l {
		for x := range l[y] {
			if l[y][x] == transparent && l2[y][x] != transparent {
				l[y][x] = l2[y][x]
			}
		}
	}
}

func part1(im image) uint {
	counts := make([]map[uint8]uint, len(im))
	for i := range counts {
		counts[i] = map[uint8]uint{0: 0, 1: 0, 2: 0}
	}
	for li := range im {
		for y := range im[li] {
			for _, d := range im[li][y] {
				counts[li][d]++
			}
		}
	}
	minZeroCount := counts[0]
	for i := range counts[1:] {
		if counts[i][0] < minZeroCount[0] {
			minZeroCount = counts[i]
		}
	}
	return minZeroCount[1] * minZeroCount[2]
}

func part2(im image) (combined *layer) {
	combined = &layer{}
	for y := range combined {
		for x := range combined[y] {
			combined[y][x] = transparent
		}
	}

	for li := range im {
		combined.mergeDown(&im[li])
	}
	return
}

// Solve provides the day 7 puzzle solution.
func Solve(input string) (*solutions.Solution, error) {
	numLayers := len(input) / layerSize
	im := make(image, numLayers)
	for i, b := range []byte(input) {
		// indexes in order: layer, y, x.
		im[i/layerSize][(i%layerSize)/width][i%width] = b - asciiDigitDiff
	}
	return &solutions.Solution{Part1: part1(im), Part2: part2(im)}, nil
}
