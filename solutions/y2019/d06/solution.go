package d06

import (
	"regexp"
	"strings"

	"github.com/jzimbel/adventofcode-go/solutions"
)

// intermediate data structure to make building the tree easier
type orbitMap map[string][]string

// standard tree structure with awareness of its parent and depth from the root
type tree struct {
	label    string
	depth    uint
	parent   *tree
	children []*tree
}

// allows for fast access to any node in a tree by its label
type flatTree map[string]*tree

const (
	centerOfMass string = "COM"
	you          string = "YOU"
	santa        string = "SAN"
)

var pattern = regexp.MustCompile(`(.+)\)(.+)`)

// makeFlatTree creates a flatTree that points to nodes in an underlying tree.
func makeFlatTree(om orbitMap) flatTree {
	ft := make(flatTree, len(om))

	var buildTree func(string, uint, *tree) *tree
	buildTree = func(label string, depth uint, parent *tree) *tree {
		t := tree{
			label:    label,
			children: make([]*tree, 0, len(om[label])),
			depth:    depth,
			parent:   parent,
		}
		for i := range om[label] {
			t.children = append(t.children, buildTree(om[label][i], depth+1, &t))
		}
		ft[label] = &t
		return &t
	}

	buildTree(centerOfMass, 0, nil)
	return ft
}

func (ft flatTree) sumDepths() (count uint) {
	for _, t := range ft {
		count += t.depth
	}
	return
}

func (t *tree) getAncestors() []*tree {
	ancestors := make([]*tree, t.depth)
	current := t
	for i := uint(0); i < t.depth; i++ {
		current = current.parent
		ancestors[i] = current
	}
	return ancestors
}

func getCommonAncestor(t1 *tree, t2 *tree) (common *tree) {
	ancestors1 := t1.getAncestors()
	ancestors2 := t2.getAncestors()

	// put the second list of ancestors into a set for faster existence checking between the two lists
	compareSet := make(map[*tree]struct{}, len(ancestors2))
	for _, ancestor := range ancestors2 {
		compareSet[ancestor] = struct{}{}
	}
	for _, ancestor := range ancestors1 {
		if _, ok := compareSet[ancestor]; ok {
			common = ancestor
			break
		}
	}
	return
}

func parseInput(input string) (om orbitMap) {
	lines := strings.Split(input, "\n")
	om = make(orbitMap, len(lines))
	for _, line := range lines {
		matches := pattern.FindStringSubmatch(line)
		base, satellite := matches[1], matches[2]
		if _, ok := om[base]; ok {
			om[base] = append(om[base], satellite)
		} else {
			om[base] = []string{satellite}
		}
	}
	return
}

func part1(ft flatTree) uint {
	return ft.sumDepths()
}

func part2(ft flatTree) uint {
	common := getCommonAncestor(ft[you].parent, ft[santa].parent)
	return (ft[you].parent.depth - common.depth) + (ft[santa].parent.depth - common.depth)
}

// Solve provides the day 6 puzzle solution.
func Solve(input string) (*solutions.Solution, error) {
	ft := makeFlatTree(parseInput(input))
	return &solutions.Solution{Part1: part1(ft), Part2: part2(ft)}, nil
}
