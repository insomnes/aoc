package solution

import (
	"fmt"
	"time"
)

type HikingMap struct {
	Map        [][]int8
	ZeroPoints []Point[int8]
}

type ParsedInput = HikingMap

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track(time.Now(), "ParseInput")
	zeroPoints := make([]Point[int8], 0)
	elevations := make([][]int8, len(lines))
	var cols int

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
		parsedRow := make([]int8, cols)

		for c, char := range row {
			n := int(char - '0')
			if n < 0 || n > 9 {
				return HikingMap{}, fmt.Errorf("invalid elevation %c at %d,%d", char, r, c)
			}
			parsedRow[c] = int8(n)
			if n == 0 {
				zeroPoints = append(zeroPoints, Point[int8]{r, c, int8(n)})
			}
		}
		elevations[r] = parsedRow
	}

	return HikingMap{elevations, zeroPoints}, nil
}

func PartOne(inp ParsedInput) int {
	defer Track(time.Now(), "PartOne")
	total := 0
	grid, err := NewGridFromSlices(inp.Map)
	if err != nil {
		panic(fmt.Sprintf("creating grid %s", err))
	}
	cache := NewGridWithDefault[map[[2]int]struct{}](grid.Rows, grid.Columns, nil)

	for _, zp := range inp.ZeroPoints {
		total += len(findNines(grid, zp, cache))
	}

	return total
}

func findNines(
	grid *Grid[int8],
	point Point[int8],
	cache *Grid[map[[2]int]struct{}],
) map[[2]int]struct{} {
	if point.Value == 9 {
		return map[[2]int]struct{}{{point.Row, point.Col}: {}}
	}
	if val := cache.Get(point.Row, point.Col).Value; val != nil {
		return val
	}
	nines := make(map[[2]int]struct{})
	for _, neighbor := range grid.GetNeighborsCross(point.Row, point.Col) {
		if neighbor.Value != point.Value+1 {
			continue
		}

		for p := range findNines(grid, neighbor, cache) {
			nines[p] = struct{}{}
		}
	}
	cache.Set(point.Row, point.Col, nines)
	return nines
}

func PartTwo(inp ParsedInput) int {
	defer Track(time.Now(), "PartTwo")
	total := 0
	grid, err := NewGridFromSlices(inp.Map)
	if err != nil {
		panic(fmt.Sprintf("creating grid %s", err))
	}

	cache := NewGridWithDefault(grid.Rows, grid.Columns, -1)

	for _, zp := range inp.ZeroPoints {
		total += findDistinctWays(grid, zp, cache)
	}

	return total
}

func findDistinctWays(grid *Grid[int8], point Point[int8], cache *Grid[int]) int {
	if point.Value == 9 {
		return 1
	}
	cached := cache.Get(point.Row, point.Col).Value
	if cached != -1 {
		return cached
	}

	ways := 0
	for _, neighbor := range grid.GetNeighborsCross(point.Row, point.Col) {
		if neighbor.Value != point.Value+1 {
			continue
		}
		ways += findDistinctWays(grid, neighbor, cache)
	}
	cache.Set(point.Row, point.Col, ways)
	return ways
}
