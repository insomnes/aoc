package solution

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	ByteChar  = '#'
	EmptyChar = '.'
)

type FallingBytes []Point[rune]

type ParsedInput = FallingBytes

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track(time.Now(), "ParseInput")

	fallBytes := make(FallingBytes, len(lines))
	for i, line := range lines {
		p, err := parsePoint(line)
		if err != nil {
			return FallingBytes{}, err
		}
		fallBytes[i] = p
	}

	return fallBytes, nil
}

// X,Y (Col,Row)
func parsePoint(s string) (Point[rune], error) {
	var p Point[rune]
	split := strings.Split(s, ",")
	if len(split) != 2 {
		return p, fmt.Errorf("invalid point format: %s", s)
	}
	colR, rowR := split[0], split[1]
	col, err := strconv.Atoi(colR)
	if err != nil {
		return p, err
	}
	p.Col = col

	row, err := strconv.Atoi(rowR)
	if err != nil {
		return p, err
	}
	p.Row = row

	p.Value = ByteChar
	return p, nil
}

const (
	TestGridSize  = 7
	TestTimeFrame = 12
	GridSize      = 71
	TimeFrame     = 1024
)

func prepareStartGrid(fallingBytes FallingBytes, size int, timeFrame int) *Grid[rune] {
	grid := NewGridWithDefault(size, size, EmptyChar)
	for _, p := range fallingBytes[:timeFrame] {
		grid.Set(p.Row, p.Col, p.Value)
	}
	return grid
}

type Step struct {
	Row, Col int
	StpCnt   int
	Metric   int
}

func findNextSteps(g *Grid[rune], step Step, visited *Grid[Step]) []Step {
	steps := make([]Step, 0, 4)
	for _, neighbor := range g.GetNeighborsCross(step.Row, step.Col) {
		if neighbor.Value == ByteChar {
			continue
		}
		if visited.Get(neighbor.Row, neighbor.Col).Value.StpCnt != -1 {
			continue
		}
		row, col, sCount := neighbor.Row, neighbor.Col, step.StpCnt+1
		metric := manhattanToEnd(row, col, g)
		newStep := Step{Row: row, Col: col, StpCnt: sCount, Metric: metric}
		steps = append(steps, newStep)

	}
	return steps
}

func lessForDijkstra(a, b Step) bool {
	return a.StpCnt < b.StpCnt
}

func lessForAStar(a, b Step) bool {
	return a.StpCnt+a.Metric < b.StpCnt+b.Metric
}

func lessForGreedyAStar(a, b Step) bool {
	return a.Metric < b.Metric
}

func FindShortestPath(g *Grid[rune], less func(a, b Step) bool) int {
	unknownStep := Step{Row: -1, Col: -1, StpCnt: -1, Metric: -1}
	visited := NewGridWithDefault(g.Rows, g.Columns, unknownStep)
	endRow, endCol := g.Rows-1, g.Columns-1

	startStep := Step{
		Row:    0,
		Col:    0,
		StpCnt: 0,
		Metric: manhattanToEnd(0, 0, g),
	}
	queue := NewHeap[Step](less, 200)
	queue.Push(startStep)

	for queue.Len() > 0 {
		step := queue.Pop()
		if step.Row == endRow && step.Col == endCol {
			visited.Set(step.Row, step.Col, step)
			break
		}

		knownStep := visited.Get(step.Row, step.Col).Value
		if knownStep.StpCnt != -1 {
			continue
		}

		visited.Set(step.Row, step.Col, step)
		nextSteps := findNextSteps(g, step, visited)

		for _, nextStep := range nextSteps {
			queue.Push(nextStep)
		}
	}

	return visited.Get(endRow, endCol).Value.StpCnt
}

func PartOne(inp ParsedInput) int {
	defer Track(time.Now(), "PartOne")
	size, tf := GridSize, TimeFrame

	grid := prepareStartGrid(inp, size, tf)

	return FindShortestPath(grid, lessForAStar)
}

func PartTwo(inp ParsedInput) int {
	defer Track(time.Now(), "PartTwo")

	size, tf := GridSize, TimeFrame
	grid := prepareStartGrid(inp, size, tf)
	bp := bisectBlockerSearch(grid, inp, tf)

	return bp.Col*1000000 + bp.Row
}

// We can skip bytes until the time frame, cause we know they are not blocking anything
// from part one. Also grid already contains them, so we don't need to set them again.
func iterativeBlockerSearch(g *Grid[rune], bytes FallingBytes, tf int) Point[rune] {
	for i, p := range bytes[tf+1:] {
		g.Set(p.Row, p.Col, p.Value)
		if isBlocked(g) {
			fmt.Printf("Blocked at %d: %d,%d\n", tf+i, p.Col, p.Row)
			return p
		}
	}
	panic("No blocker found")
}

// This is a binary search implementation to find the blocker:
// https://en.wikipedia.org/wiki/Binary_search
// The idea is to block cells by clusters to find the blocker by log(n) steps.
// We also don't need to check starting timeframe bytes
func bisectBlockerSearch(g *Grid[rune], bytes FallingBytes, tf int) Point[rune] {
	start, end := tf+1, len(bytes)
	for start < end {
		if end-start == 1 {
			fmt.Printf("Blocked at %d: %d,%d\n", start, bytes[start].Col, bytes[start].Row)
			return bytes[start]
		}
		mid := (start + end) / 2
		blockCells(g, bytes[start:mid])
		// This means we can block further aka we need to move mid to right
		if !isBlocked(g) {
			start = mid
			continue
		}
		// Otherwise we move mid to left and end to previous mid (bisect)
		// and free the cells from new mid to new end
		end = mid
		futureMid := (start + end) / 2
		freeCells(g, bytes[futureMid:end])
	}

	panic("No blocker found")
}

func blockCells(g *Grid[rune], toBlock FallingBytes) {
	for _, p := range toBlock {
		g.Set(p.Row, p.Col, p.Value)
	}
}

func freeCells(g *Grid[rune], toFree FallingBytes) {
	for _, p := range toFree {
		g.Set(p.Row, p.Col, EmptyChar)
	}
}

func isBlocked(g *Grid[rune]) bool {
	return FindShortestPath(g, lessForGreedyAStar) == -1
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func manhattanToEnd(row, col int, g *Grid[rune]) int {
	endRow, endCol := g.Rows-1, g.Columns-1
	return abs(row-endRow) + abs(col-endCol)
}
