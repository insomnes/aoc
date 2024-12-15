package solution

import (
	"fmt"
	"strings"
	"time"
)

type (
	WhsGrid = Grid[rune]
	Cell    = Point[rune]
)

const (
	EmptyChar   = '.'
	RobotChar   = '@'
	BoxChar     = 'O'
	WallChar    = '#'
	LboxChar    = '['
	RboxChar    = ']'
	UnknownChar = '?'
)

type WarehouseInfo struct {
	Map   [][]rune
	Moves []GridMove
	Robot Cell
}

type ParsedInput = WarehouseInfo

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track(time.Now(), "ParseInput")
	var cols int

	warhouseMap := make([][]rune, 0)

	r := -1
	robot := Cell{-1, -1, UnknownChar}
	rows := len(lines)
	for {
		r++
		if r >= rows-1 {
			return ParsedInput{}, fmt.Errorf("No empty line found")
		}

		line := lines[r]
		if line == "" {
			break
		}

		if cols != 0 && len(line) != cols {
			return ParsedInput{}, fmt.Errorf("%d != %d at %d", len(line), cols, r)
		}
		cols = len(line)

		row := make([]rune, cols)
		for c, char := range line {
			if char == RobotChar {
				robot = Cell{r, c, RobotChar}
				char = EmptyChar
			}
			row[c] = char
		}
		warhouseMap = append(warhouseMap, row)
	}

	if robot.Value == UnknownChar {
		return ParsedInput{}, fmt.Errorf("robot not found")
	}
	if err := checkWarehouseWalls(warhouseMap); err != nil {
		return ParsedInput{}, err
	}

	moves := make([]GridMove, 0)
	for _, line := range lines[r+1:] {
		for _, char := range line {
			move, err := parseMove(char)
			if err != nil {
				return ParsedInput{}, fmt.Errorf("parsing move: %s", err)
			}
			moves = append(moves, move)
		}
	}

	return ParsedInput{
		Map:   warhouseMap,
		Moves: moves,
		Robot: robot,
	}, nil
}

// The walls' presence allows us to skip most out-of-bounds checks
func checkWarehouseWalls(warhouseMap [][]rune) error {
	rows, cols := len(warhouseMap), len(warhouseMap[0])

	for r := 0; r < rows; r++ {
		if r == 0 || r == rows-1 {
			for _, char := range warhouseMap[r] {
				if char != WallChar {
					return fmt.Errorf("wall expected at %d, %d", r, 0)
				}
			}
			continue
		}
		for c := 0; c < cols; c++ {
			if c == 0 || c == cols-1 {
				if warhouseMap[r][c] != WallChar {
					return fmt.Errorf("wall expected at %d, %d", r, c)
				}
			}
		}
	}
	return nil
}

func parseMove(char rune) (GridMove, error) {
	switch char {
	case '^':
		return MoveUp, nil
	case 'v':
		return MoveDown, nil
	case '<':
		return MoveLeft, nil
	case '>':
		return MoveRight, nil
	default:
		return MoveUp, fmt.Errorf("unknown move %c", char)
	}
}

func PartOne(inp ParsedInput) int {
	defer Track(time.Now(), "PartOne")

	robot := inp.Robot
	grid, err := NewGridFromSlices(inp.Map)
	if err != nil {
		panic(fmt.Sprintf("creating grid: %s", err))
	}

	// We can check if the move is the same, but it's quite fast as is
	for _, move := range inp.Moves {
		robot = moveRobot(grid, robot, move)
	}
	return calculateGPS(grid, BoxChar)
}

func moveRobot(grid *WhsGrid, robot Cell, move GridMove) Cell {
	next := grid.GetAnyByMove(robot.Row, robot.Col, move)
	if next.Value == EmptyChar {
		next.Value = RobotChar
		return next
	}
	if next.Value == WallChar {
		return robot
	}

	endCell := findBoxChainEnd(grid, next, move)
	if endCell.Value != EmptyChar {
		return robot
	}
	// Moving straight line of the boxes in our case can be coded as
	// "teleporting" the first box to the end of the chain. For our purposes
	// any box marker is indistinguishable from other boxes
	grid.Set(next.Row, next.Col, EmptyChar)
	grid.Set(endCell.Row, endCell.Col, BoxChar)

	next.Value = RobotChar
	return next
}

// Check next cell in the direction of the move until we hit a wall or empty cell
func findBoxChainEnd(grid *WhsGrid, startBox Cell, move GridMove) Cell {
	current, next := startBox, startBox
	for next.Value != WallChar && next.Value != EmptyChar {
		next = grid.GetAnyByMove(current.Row, current.Col, move)
		current = next
	}
	return current
}

func calculateGPS(g *WhsGrid, boxChar rune) int {
	gps := 0
	for _, cell := range g.Cells {
		if cell.Value != boxChar {
			continue
		}
		gps += calcCellGPS(cell)
	}
	return gps
}

func calcCellGPS(cell Cell) int {
	return 100*cell.Row + cell.Col
}

func PartTwo(inp ParsedInput) int {
	defer Track(time.Now(), "PartTwo")

	grid := expandWarehouse(inp.Map)
	// Expanded warehouse means that we need to double robot's starting column
	robot, moves := inp.Robot, inp.Moves
	robot.Col = robot.Col * 2

	for _, move := range moves {
		robot = moveRobotWide(grid, robot, move)
	}
	return calculateGPS(grid, LboxChar)
}

// Only columns are expanded, rows are the same
func expandWarehouse(whmap [][]rune) *WhsGrid {
	origRows, origCols := len(whmap), len(whmap[0])
	rows, cols := origRows, origCols*2

	grid := NewGrid[rune](rows, cols)

	for r := 0; r < origRows; r++ {
		for c := 0; c < origCols; c++ {
			char := whmap[r][c]
			expCol := c * 2
			if char != BoxChar {
				grid.Set(r, expCol, char)
				grid.Set(r, expCol+1, char)
				continue
			}
			grid.Set(r, expCol, LboxChar)
			grid.Set(r, expCol+1, RboxChar)
		}
	}
	return grid
}

func moveRobotWide(grid *WhsGrid, robot Cell, move GridMove) Cell {
	next := grid.GetAnyByMove(robot.Row, robot.Col, move)
	if next.Value == WallChar {
		return robot
	}
	if next.Value == EmptyChar || tryMoveBox(grid, next, move) {
		robot.Row, robot.Col = next.Row, next.Col
	}
	return robot
}

// Horizontal move is very similar to the PartOne's moveRobot, but we need to
// account for the new box size and recursively move all boxes in the chain
func tryMoveBox(grid *WhsGrid, box Cell, move GridMove) bool {
	if move == MoveLeft || move == MoveRight {
		return moveBoxHorizontally(grid, box, move)
	}
	return movePartVertically(grid, box, move)
}

// Horizontal move is very similar to the PartOne's moveRobot, but we need to
// account for the new box size and recursively move all boxes in the chain
func moveBoxHorizontally(grid *WhsGrid, box Cell, move GridMove) bool {
	// In the end of the box chain we can decide if we can move all previous boxes
	// based on the next cell value
	if box.Value == EmptyChar || box.Value == WallChar {
		return box.Value == EmptyChar
	}
	side := toOtherBoxSide(box)
	next := grid.GetAnyByMove(box.Row, side.Col, move)
	if !moveBoxHorizontally(grid, next, move) {
		return false
	}

	grid.Set(next.Row, next.Col, side.Value)
	grid.Set(box.Row, side.Col, box.Value)
	grid.Set(box.Row, box.Col, EmptyChar)
	return true
}

func movePartVertically(grid *WhsGrid, box Cell, move GridMove) bool {
	stackedParts := findStackedParts(grid, box, move)
	if len(stackedParts) == 0 {
		return false
	}

	// It is important to move boxes from the far end of the stack
	for i := len(stackedParts) - 1; i >= 0; i-- {
		stackedPart := stackedParts[i]
		movePart(grid, stackedPart, move)
	}
	return true
}

// This is kind of a BFS, but with the specific neighbor finding
func findStackedParts(grid *WhsGrid, part Cell, move GridMove) []Cell {
	visited := make(map[int]struct{})

	sidePart := toOtherBoxSide(part)
	parts := make([]Cell, 0)
	queue := []Cell{part, sidePart}

	for len(queue) > 0 {
		// It is important to stack the boxes in the BFS order, layer by layer
		cur := queue[0]
		queue = queue[1:]
		key := calcCellGPS(cur)
		if _, ok := visited[key]; ok {
			continue
		}
		visited[key] = struct{}{}
		parts = append(parts, cur)

		following := findFollowingParts(grid, cur, move)
		if len(following) == 1 {
			return []Cell{}
		}

		for _, neighbor := range following {
			if _, ok := visited[calcCellGPS(neighbor)]; ok {
				continue
			}
			queue = append(queue, neighbor)
		}
	}
	return parts
}

func findFollowingParts(grid *WhsGrid, box Cell, move GridMove) []Cell {
	following := grid.GetAnyByMove(box.Row, box.Col, move)
	if following.Value == EmptyChar {
		return []Cell{}
	}
	if following.Value == WallChar {
		return []Cell{following}
	}
	return []Cell{following, toOtherBoxSide(following)}
}

func movePart(grid *WhsGrid, part Cell, move GridMove) {
	if part.Value != LboxChar && part.Value != RboxChar {
		panic(fmt.Sprintf("moving box %c at %d, %d", part.Value, part.Row, part.Col))
	}
	next := grid.GetAnyByMove(part.Row, part.Col, move)
	if next.Value != EmptyChar {
		panic(fmt.Sprintf("move to '%c' at %d, %d", next.Value, next.Row, next.Col))
	}

	grid.Set(next.Row, next.Col, part.Value)
	grid.Set(part.Row, part.Col, EmptyChar)
}

func toOtherBoxSide(cell Cell) Cell {
	sideChar, sideCol := RboxChar, cell.Col+1
	if cell.Value == RboxChar {
		sideChar, sideCol = LboxChar, cell.Col-1
	}

	return Cell{cell.Row, sideCol, sideChar}
}

func StringWarehouse(grid *WhsGrid, robot Cell) string {
	var sb strings.Builder

	for i, cell := range grid.Cells {
		if cell.Row == robot.Row && cell.Col == robot.Col {
			sb.WriteRune(RobotChar)
		} else {
			sb.WriteRune(cell.Value)
		}
		if (i+1)%grid.Columns == 0 && i != len(grid.Cells)-1 {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}
