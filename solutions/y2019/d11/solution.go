package d11

import (
	"strings"

	"github.com/jzimbel/adventofcode-go/solutions"
	"github.com/jzimbel/adventofcode-go/solutions/y2019/interpreter"
)

const (
	outputTypeColor = iota
	outputTypeTurn
)

const (
	black = iota
	white
)

// mod performs a modulo that always floors the remainder toward 0 and uses the sign of the divisor.
// mod(6,5) = 1
// mod(-6,5) = 4
// mod(6,-5) = -4
// mod(-6,-5) = -1
func mod(a, b int) (n int) {
	n = a % b
	if (a < 0 && b > 0) || (a > 0 && b < 0) {
		n += b
	}
	return
}

type point struct {
	x int
	y int
}
type direction point

type grid map[point]int

type row []rune
type image []row

var (
	dirs = [...]*direction{
		&direction{0, -1},
		&direction{1, 0},
		&direction{0, 1},
		&direction{-1, 0},
	}
	colorMap = map[int]rune{
		black: ' ',
		white: '\u2588',
	}
)

func (p *point) move(d *direction) {
	p.x += d.x
	p.y += d.y
}

func (g grid) getBounds() (xOffset, yOffset, xMax, yMax int) {
	for p := range g {
		if p.x < xOffset {
			xOffset = p.x
		}
		if p.x > xMax {
			xMax = p.x
		}
		if p.y < yOffset {
			yOffset = p.y
		}
		if p.y > yMax {
			yMax = p.y
		}
	}
	xMax, yMax = xMax-xOffset, yMax-yOffset
	return
}

func (g grid) String() string {
	xOffset, yOffset, xMax, yMax := g.getBounds()
	im := make(image, yMax+1)
	for y := range im {
		im[y] = make(row, xMax+1)
		for x := range im[y] {
			im[y][x] = colorMap[black]
		}
	}

	for p, c := range g {
		im[p.y-yOffset][p.x-xOffset] = colorMap[c]
	}

	builder := make([]string, len(im))
	for y := range im {
		builder[y] = string(im[y])
	}
	return strings.Join(builder, "\n")
}

func run(initMem interpreter.Program, startColor int) (paintedCount int, g grid) {
	var outputType, dirIndex int
	g = make(grid)
	if startColor == white {
		g[point{}] = startColor
	}
	cursor := point{}

	input := func() (color int) {
		var ok bool
		if color, ok = g[cursor]; !ok {
			color = black
		}
		return
	}

	output := func(n int) {
		switch outputType {
		case outputTypeColor:
			if _, ok := g[cursor]; !ok {
				paintedCount++
			}
			g[cursor] = n
		case outputTypeTurn:
			dirIndex = mod(dirIndex+(n*2-1), len(dirs))
			cursor.move(dirs[dirIndex])
		}
		outputType = (outputType + 1) % 2
	}

	interpreter.New(initMem, input, output).Run()
	return
}

func part1(initMem interpreter.Program) (paintedCount int) {
	paintedCount, _ = run(initMem, black)
	return
}

func part2(initMem interpreter.Program) (g grid) {
	_, g = run(initMem, white)
	return
}

// Solve provides the day 11 puzzle solution.
func Solve(input string) (*solutions.Solution, error) {
	initMem := interpreter.ParseMem(input)
	return &solutions.Solution{Part1: part1(initMem), Part2: part2(initMem)}, nil
}
