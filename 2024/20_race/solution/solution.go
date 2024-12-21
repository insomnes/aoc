package solution

import (
	"fmt"
	"slices"
	"time"
)

const (
	StartChar   = 'S'
	EndChar     = 'E'
	WallChar    = '#'
	EmptyChar   = '.'
	UnknownChar = '?'
)

type Maze struct {
	PathCells  []PathCell
	SparseGrid *PathSparseGrid
}

type PathDistance struct {
	DRow int
	DCol int
	Dist int
}

func calcPastDistance(r, c, or, oc int) PathDistance {
	return PathDistance{DRow: r - or, DCol: c - oc, Dist: abs(r-or) + abs(c-oc)}
}

type (
	PathCell       = Point[int]
	PathGrid       = Grid[int]
	PathSparseGrid = SparseGrid[int]
)

const UnknownTime = -1

type ParsedInput = Maze

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track(time.Now(), "ParseInput")
	rows := len(lines)
	columns := len(lines[0])
	pathGrid := NewGrid[int](rows, columns)

	start := PathCell{Row: -1, Col: -1, Value: -1}
	end := PathCell{Row: -1, Col: -1, Value: -1}

	for r, row := range lines {
		for c, cell := range row {
			val := 0
			if c >= columns {
				return ParsedInput{}, fmt.Errorf("inconsistent width %d at %d", c, r)
			}
			if cell == StartChar {
				start = PathCell{Row: r, Col: c, Value: 0}
			}
			if cell == EndChar {
				end = PathCell{Row: r, Col: c, Value: 0}
			}
			if cell == WallChar {
				val = UnknownTime
			}
			pathGrid.Set(r, c, val)
		}
	}
	if start.Row == -1 || end.Row == -1 {
		return ParsedInput{}, fmt.Errorf("start or end not found")
	}
	pathCells, sparseGrid := findPath(pathGrid, start, end)

	return Maze{PathCells: pathCells, SparseGrid: sparseGrid}, nil
}

func iswall(cell int) bool {
	return cell == UnknownTime
}

func findPath(grid *PathGrid, start, end PathCell) ([]PathCell, *PathSparseGrid) {
	cur := grid.GetNeighborsCrossNoWalls(start.Row, start.Col, iswall)[0]
	move := GridMove{cur.Row - start.Row, cur.Col - start.Col}
	pathCells := []PathCell{
		{Row: start.Row, Col: start.Col, Value: 0},
	}
	spGrid := NewSparseGrid[int](grid.Rows, grid.Columns)

	length := 0

MainLoop:
	for {
		length++

		grid.Set(cur.Row, cur.Col, length)
		if cur.Row == end.Row && cur.Col == end.Col {
			pc := PathCell{Row: cur.Row, Col: cur.Col, Value: length}
			pathCells = append(pathCells, pc)
			spGrid.AddPoint(pc)
			break
		}

		next, ok := grid.GetByMoveCheck(cur.Row, cur.Col, move)
		if ok && !iswall(next.Value) {
			pc := PathCell{Row: cur.Row, Col: cur.Col, Value: length}
			pathCells = append(pathCells, pc)
			spGrid.AddPoint(pc)
			cur = next
			continue
		}

		oppositeMove := move.Opposite()
		for _, nextMove := range GridCrossMoves {
			if nextMove == move || nextMove == oppositeMove {
				continue
			}
			next, ok := grid.GetByMoveCheck(cur.Row, cur.Col, nextMove)
			if ok && !iswall(next.Value) {
				move = nextMove
				length--
				continue MainLoop
			}
		}
		panic("no moves")
	}
	return pathCells, spGrid
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func manhattanDistance(a, b PathCell) int {
	return abs(a.Row-b.Row) + abs(a.Col-b.Col)
}

func findCheats(
	pathCells []PathCell,
	sparseGrid *PathSparseGrid,
	threshold int,
	cheatDuration int,
) int {
	tsCheats := 0
	j := 0
	for j < threshold {
		toRemove := pathCells[j]
		sparseGrid.RemovePoint(toRemove)
		j++
	}
	for _, cell := range pathCells {
		if j < len(pathCells) {
			toRemove := pathCells[j]
			sparseGrid.RemovePoint(toRemove)
			j++
		}
		ts := findCellCheats(cell, sparseGrid, threshold, cheatDuration)
		tsCheats += ts
	}
	return tsCheats
}

func findCellCheats(
	cell PathCell,
	sparseGrid *PathSparseGrid,
	threshold int,
	cheatDuration int,
) int {
	tsCheats := 0
	upperRow := max(cell.Row-cheatDuration, 0)
	lowerRow := min(cell.Row+cheatDuration, sparseGrid.NumRows-1)

	var leftCell, rightCell PathCell

	for r := upperRow; r <= lowerRow; r++ {
		row := sparseGrid.Rows[r]
		if len(row) == 0 {
			continue
		}

		rowDelta := abs(cell.Row - r)
		maxColDelta := cheatDuration - rowDelta
		lCol := cell.Col - maxColDelta
		if lCol < 0 {
			lCol = 0
		}
		rCol := min(cell.Col+maxColDelta, sparseGrid.NumColumns-1)

		leftCell, rightCell = PathCell{Row: r, Col: lCol}, PathCell{Row: r, Col: rCol}
		leftI, _ := slices.BinarySearchFunc(row, leftCell, LessPoint)
		rightI, rightFound := slices.BinarySearchFunc(row, rightCell, LessPoint)
		if rightFound {
			rightI++
		}

		for _, jumpCell := range row[leftI:rightI] {
			if jumpCell.Value <= cell.Value || jumpCell.Value-cell.Value-2 < threshold {
				continue
			}
			jumpDist := manhattanDistance(cell, jumpCell)
			// Couldn't be jumpDist > cheatDuration
			savedTime := jumpCell.Value - cell.Value - jumpDist
			if savedTime >= threshold {
				tsCheats++
			}
		}

	}

	return tsCheats
}

const (
	cheatDurationOne = 2
	thresholdOne     = 100
)

func PartOne(inp ParsedInput) int {
	defer Track(time.Now(), "PartOne")
	path, sparseGrid := inp.PathCells, inp.SparseGrid
	sparseGrid = sparseGrid.Copy()

	return findCheats(path, sparseGrid, thresholdOne, cheatDurationOne)
}

const (
	cheatDurationTwo = 20
	thresholdTwo     = 100
)

func PartTwo(inp ParsedInput) int {
	defer Track(time.Now(), "PartTwo")
	path, sparseGrid := inp.PathCells, inp.SparseGrid

	return findCheats(path, sparseGrid, thresholdTwo, cheatDurationTwo)
}
