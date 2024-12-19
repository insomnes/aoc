package solution

import (
	"fmt"
	"math"
	"time"
)

const (
	StartChar   = 'S'
	EndChar     = 'E'
	WallChar    = '#'
	EmptyChar   = '.'
	UnknownChar = '?'
)

const (
	StepCost = 1
	TurnCost = 1000
)

var startMove = MoveRight

type Maze struct {
	grid  *MazeGrid
	start Cell
	end   Cell
}

type (
	Cell     = Point[rune]
	PathCell = Point[int]
	MazeGrid = Grid[rune]
)

const UnknownScore = math.MaxInt32

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type ParsedInput = Maze

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track(time.Now(), "ParseInput")
	rows := len(lines)
	columns := len(lines[0])
	grid := NewGrid[rune](rows, columns)

	start := Point[rune]{Row: -1, Col: -1, Value: UnknownChar}
	end := Point[rune]{Row: -1, Col: -1, Value: UnknownChar}

	for r, row := range lines {
		for c, cell := range row {
			if c >= columns {
				return ParsedInput{}, fmt.Errorf("inconsistent width %d at %d", c, r)
			}
			if cell == StartChar {
				start = Point[rune]{Row: r, Col: c, Value: StartChar}
				cell = EmptyChar
			}
			if cell == EndChar {
				end = Point[rune]{Row: r, Col: c, Value: EndChar}
				cell = EmptyChar
			}
			grid.Set(r, c, cell)
		}
	}
	if start.Row == -1 || end.Row == -1 {
		return ParsedInput{}, fmt.Errorf("start or end not found")
	}

	return Maze{grid: grid, start: start, end: end}, nil
}

type PathFinder struct {
	Row, Col int
	Move     GridMove
	Score    int
	From     *PathFinder
}

func (pf PathFinder) ToKey() RouteKey {
	return RouteKey{pf.Row, pf.Col, pf.Move[0], pf.Move[1]}
}

func (pf PathFinder) MakePrevious() []RouteKey {
	if pf.From == nil {
		return nil
	}
	return []RouteKey{pf.From.ToKey()}
}

func (pf PathFinder) MakePathInfo() PathInfo {
	return PathInfo{Score: pf.Score, Previous: pf.MakePrevious()}
}

func pfless(a, b PathFinder) bool {
	return a.Score < b.Score
}

type RouteKey [4]int

func (rk RouteKey) ToPosition() [2]int {
	return [2]int{rk[0], rk[1]}
}

func shortestPath(start, end Cell, grid *MazeGrid) int {
	scores := NewGridWithDefault[int](grid.Rows, grid.Columns, UnknownScore)
	startFinder := PathFinder{Row: start.Row, Col: start.Col, Move: startMove, Score: 0}

	queue := NewHeap(pfless, 100)
	queue.Push(startFinder)

	for queue.Len() > 0 {
		pf := queue.Pop()
		if pf.Row == end.Row && pf.Col == end.Col {
			return pf.Score
		}

		knownScore := scores.Get(pf.Row, pf.Col).Value
		if pf.Score >= knownScore {
			continue
		}

		scores.Set(pf.Row, pf.Col, pf.Score)
		for _, next := range getNext(pf, grid) {
			queue.Push(next)
		}
	}

	return -1
}

func getNext(pf PathFinder, grid *MazeGrid) []PathFinder {
	oppositeMove := pf.Move.Opposite()
	next := make([]PathFinder, 0, 4)
	for _, move := range GridCrossMoves {
		if move == oppositeMove {
			continue
		}

		n := grid.GetAnyByMove(pf.Row, pf.Col, move)
		if n.Value == WallChar {
			continue
		}

		newScore := pf.Score + StepCost
		if move != pf.Move {
			newScore += TurnCost
		}
		nextPf := PathFinder{
			Row:   n.Row,
			Col:   n.Col,
			Move:  move,
			Score: newScore,
			From:  &pf,
		}
		next = append(next, nextPf)
	}
	return next
}

func PartOne(inp ParsedInput) int {
	defer Track(time.Now(), "PartOne")
	start, end, grid := inp.start, inp.end, inp.grid

	return shortestPath(start, end, grid)
}

type PathInfo struct {
	Score    int
	Previous []RouteKey
}

func shortestWithPathSaving(start, end Cell, grid *MazeGrid) (map[RouteKey]PathInfo, RouteKey) {
	unknowInfo := PathInfo{Score: UnknownScore, Previous: nil}
	allPathsInfo := make(map[RouteKey]PathInfo)
	startFinder := PathFinder{Row: start.Row, Col: start.Col, Move: startMove, Score: 0}

	queue := NewHeap(pfless, 100)
	queue.Push(startFinder)

	var endKey RouteKey

	for queue.Len() > 0 {
		pf := queue.Pop()
		pfKey := pf.ToKey()

		if pf.Row == end.Row && pf.Col == end.Col {
			pathInfo := pf.MakePathInfo()
			allPathsInfo[pfKey] = pathInfo
			endKey = pfKey
			break
		}

		knownInfo, ok := allPathsInfo[pfKey]
		if !ok {
			knownInfo = unknowInfo
		}

		if pf.Score == knownInfo.Score {
			knownInfo.Previous = append(knownInfo.Previous, pf.From.ToKey())
			allPathsInfo[pfKey] = knownInfo
			continue
		}

		if pf.Score > knownInfo.Score {
			continue
		}

		allPathsInfo[pfKey] = pf.MakePathInfo()
		for _, next := range getNext(pf, grid) {
			queue.Push(next)
		}
	}

	return allPathsInfo, endKey
}

func countUniquePoints(allInfo map[RouteKey]PathInfo, endKey RouteKey) int {
	unique := make(map[[2]int]struct{})
	toVisit := []RouteKey{endKey}

	for len(toVisit) > 0 {
		key := toVisit[len(toVisit)-1]
		toVisit = toVisit[:len(toVisit)-1]

		cur, ok := allInfo[key]
		if !ok {
			panic(fmt.Sprintf("key %v not found", key))
		}

		unique[key.ToPosition()] = struct{}{}

		for _, prevKey := range cur.Previous {
			toVisit = append(toVisit, prevKey)
		}
	}

	return len(unique)
}

func PartTwo(inp ParsedInput) int {
	defer Track(time.Now(), "PartTwo")

	start, end, grid := inp.start, inp.end, inp.grid
	info, endKey := shortestWithPathSaving(start, end, grid)
	return countUniquePoints(info, endKey)
}
