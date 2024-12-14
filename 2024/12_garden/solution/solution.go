package solution

import (
	"fmt"
	"time"
)

type GardenMap struct {
	Map [][]rune
}

type ParsedInput = GardenMap

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track(time.Now(), "ParseInput")
	var cols int

	parsed := make([][]rune, len(lines))

	for r, row := range lines {
		if cols != 0 && len(row) != cols {
			return ParsedInput{}, fmt.Errorf(
				"inconsistent row length %d != %d at row %d",
				len(row),
				cols,
				r,
			)
		}
		cols = len(row)
		parsed[r] = []rune(row)

	}

	return GardenMap{Map: parsed}, nil
}

func PartOne(inp ParsedInput) int {
	defer Track(time.Now(), "PartOne")
	total := 0
	grid, err := NewGridFromSlices(inp.Map)
	if err != nil {
		panic(fmt.Sprintf("creating grid %s", err))
	}

	visited := NewGridWithDefault(grid.Rows, grid.Columns, false)

	for pi, point := range grid.Cells {
		if visited.Cells[pi].Value {
			continue
		}

		area, perimeter := findFullPlot(grid, visited, point)
		total += area * perimeter
	}

	return total
}

const (
	cellSides = 4
)

// area, perimeter
func findFullPlot(grid *Grid[rune], globalVisited *Grid[bool], start Point[rune]) (int, int) {
	queue := make([]Point[rune], 0, 16)
	queue = append(queue, start)

	area, perimeter := 0, 0

	for len(queue) > 0 {
		point := queue[0]
		queue = queue[1:]
		if globalVisited.Get(point.Row, point.Col).Value {
			continue
		}

		globalVisited.Set(point.Row, point.Col, true)
		area++

		nCode, polyNeighbors := encodePointNeighbors(grid, point)
		cellPerimiter := calcPerimeter(nCode)
		perimeter += cellPerimiter

		for _, neighbor := range polyNeighbors {
			queue = append(queue, neighbor)
		}
	}

	return area, perimeter
}

func PartTwo(inp ParsedInput) int {
	defer Track(time.Now(), "PartTwo")
	total := 0
	grid, err := NewGridFromSlices(inp.Map)
	if err != nil {
		panic(fmt.Sprintf("creating grid %s", err))
	}

	visited := NewGridWithDefault(grid.Rows, grid.Columns, false)

	for pi, point := range grid.Cells {
		if visited.Cells[pi].Value {
			continue
		}
		// We can calculate sides by finding angles count
		area, sides := findFullPlotWithSides(grid, visited, point)
		total += area * sides

	}

	return total
}

// area, sides
func findFullPlotWithSides(
	grid *Grid[rune],
	globalVisited *Grid[bool],
	start Point[rune],
) (int, int) {
	queue := make([]Point[rune], 0, 16)
	queue = append(queue, start)

	area := 0
	angles := 0

	for len(queue) > 0 {
		point := queue[0]
		queue = queue[1:]

		if globalVisited.Get(point.Row, point.Col).Value {
			continue
		}
		globalVisited.Set(point.Row, point.Col, true)
		area++

		nCode, polyNeighbors := encodePointNeighbors(grid, point)
		// how many sides our polygon has can be calculated by the number of angles it has
		pAngles := findAngles(nCode)
		angles += pAngles

		for _, neighbor := range polyNeighbors {
			queue = append(queue, neighbor)
		}
	}

	return area, angles
}

// Each cell has 1 bit around
// 0 1 2
// 3 X 4
// 5 6 7
// So we can use bit mask to find the borders and angles by setting
// the bits for each neighbor: 0 for the same name, 1 for different
func prepareByte(powers []int) int {
	mask := 0
	for _, power := range powers {
		if power < 0 || power > 7 {
			panic(fmt.Sprintf("invalid power %d", power))
		}
		mask |= 1 << power
	}
	return mask
}

var (
	// 2, 16, 64, 8
	upBorderMask    int = prepareByte([]int{1})
	rightBorderMask int = prepareByte([]int{4})
	downBorderMask  int = prepareByte([]int{6})
	leftBorderMask  int = prepareByte([]int{3})

	// 11, 22, 208, 104
	upLeftMask    int = prepareByte([]int{0, 1, 3})
	upRightMask   int = prepareByte([]int{1, 2, 4})
	downRightMask int = prepareByte([]int{4, 6, 7})
	downLeftMask  int = prepareByte([]int{3, 5, 6})

	// 1 0 .
	// 0 X .
	upLeftInner int = prepareByte([]int{0})
	// 1 1 .
	// 1 X .
	upLeftOuter int = prepareByte([]int{0, 1, 3})
	// . 1 .
	// 1 X .
	upLeftInside int = prepareByte([]int{1, 3})

	// . 0 1
	// . X 0
	upRightInner int = prepareByte([]int{2})
	// . 1 1
	// . X 1
	upRightOuter int = prepareByte([]int{1, 2, 4})
	// . 1 .
	// . X 1
	upRightInside int = prepareByte([]int{1, 4})

	// . X 0
	// . 0 1
	downRightInner int = prepareByte([]int{7})
	// . X 1
	// . 1 1
	downRightOuter int = prepareByte([]int{4, 6, 7})
	// . X 1
	// . 1 .
	downRightInside int = prepareByte([]int{4, 6})

	// 0 X .
	// 1 0 .
	downLeftInner int = prepareByte([]int{5})
	// 1 X .
	// 1 1 .
	downLeftOuter int = prepareByte([]int{3, 5, 6})
	// 1 X .
	// . 1 .
	downLeftInside int = prepareByte([]int{3, 6})
)

var borderMasks = []int{upBorderMask, rightBorderMask, downBorderMask, leftBorderMask}

func calcPerimeter(nCode int) int {
	perimeter := 0
	for _, mask := range borderMasks {
		if nCode&mask == mask {
			perimeter++
		}
	}
	return perimeter
}

func findAngles(nCode int) int {
	angles := 0
	upLeft := nCode & upLeftMask
	upRight := nCode & upRightMask
	downRight := nCode & downRightMask
	downLeft := nCode & downLeftMask

	if upLeft == upLeftOuter || upLeft == upLeftInner || upLeft == upLeftInside {
		angles++
	}

	if upRight == upRightOuter || upRight == upRightInner || upRight == upRightInside {
		angles++
	}

	if downRight == downRightOuter || downRight == downRightInner || downRight == downRightInside {
		angles++
	}

	if downLeft == downLeftOuter || downLeft == downLeftInner || downLeft == downLeftInside {
		angles++
	}

	return angles
}

func encodePointNeighbors(grid *Grid[rune], point Point[rune]) (int, []Point[rune]) {
	// 0 1 2
	// 3 X 4
	// 5 6 7
	row, col := point.Row, point.Col
	moves := []GridMove{
		MoveUpLeft, MoveUp, MoveUpRight,
		MoveLeft, MoveRight,
		MoveDownLeft, MoveDown, MoveDownRight,
	}

	name := point.Value

	codedNeighbors := 0
	crossNeighbors := make([]Point[rune], 0, 4)

	for i, move := range moves {
		nPoint := grid.GetAnyByMove(row, col, move)
		if nPoint.Value != name {
			codedNeighbors |= 1 << i
			continue
		}
		if move == MoveUp || move == MoveRight || move == MoveDown || move == MoveLeft {
			crossNeighbors = append(crossNeighbors, nPoint)
		}
	}

	return codedNeighbors, crossNeighbors
}
