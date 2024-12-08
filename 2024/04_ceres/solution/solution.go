package solution

import (
	"fmt"
)

type ParsedInput = Grid[rune]

const (
	X = 'X'
	M = 'M'
	A = 'A'
	S = 'S'
)

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track("ParseInput")()
	rows, columns := len(lines), len(lines[0])
	g := NewGrid[rune](rows, columns)

	for row, line := range lines {
		if len(line) != columns {
			return Grid[rune]{}, fmt.Errorf(
				"Bad %dth line length: %d != %d",
				row,
				len(line),
				columns,
			)
		}
		for col, r := range line {
			g.Set(row, col, r)
		}
	}

	return *g, nil
}

func PartOne(inp ParsedInput) int {
	defer Track("PartOne")()

	xmasCount := 0

	for row := 0; row < inp.Rows; row++ {
		for col := 0; col < inp.Columns; col++ {
			point := inp.Get(row, col)
			if point.Value != X {
				continue
			}
			paths := findXMASPathsCount(&inp, point)
			xmasCount += paths
		}
	}

	return xmasCount
}

func findXMASPathsCount(g *Grid[rune], point Point[rune]) int {
	paths := 0

	for _, move := range GridAllMoves {
		row, col := point.Row, point.Col
	CharLoop:
		for _, char := range []rune{M, A, S} {
			row, col = row+move[0], col+move[1]
			nextPoint, ok := g.GetWithCheck(row, col)
			if !ok || nextPoint.Value != char {
				break CharLoop
			}
			if char == S {
				paths++
			}
		}
	}

	return paths
}

func PartTwo(inp ParsedInput) int {
	defer Track("PartTwo")()
	masCount := 0

	// We can skip the first row and column for A search, caue it's 3x3 square:
	// M . S
	// . A .
	// M . S
	for row := 1; row < inp.Rows; row++ {
		for col := 1; col < inp.Columns; col++ {
			point := inp.Get(row, col)
			if point.Value != A {
				continue
			}
			if findMAS(&inp, point) {
				masCount++
			}
		}
	}

	return masCount
}

func findMAS(g *Grid[rune], point Point[rune]) bool {
	return checkUpLeft(g, point) && checkUpRight(g, point)
}

func checkUpLeft(g *Grid[rune], point Point[rune]) bool {
	upLeft, ok := g.GetUpLeft(point.Row, point.Col)
	expectedDownRight := getMASmissing(upLeft.Value)

	if !ok || expectedDownRight == notMASchar {
		return false
	}
	downRight, ok := g.GetDownRight(point.Row, point.Col)
	if !ok || downRight.Value != expectedDownRight {
		return false
	}
	return true
}

func checkUpRight(g *Grid[rune], point Point[rune]) bool {
	upRight, ok := g.GetUpRight(point.Row, point.Col)
	expectedDownLeft := getMASmissing(upRight.Value)
	if !ok || expectedDownLeft == notMASchar {
		return false
	}

	downLeft, ok := g.GetDownLeft(point.Row, point.Col)
	if !ok || downLeft.Value != expectedDownLeft {
		return false
	}

	return true
}

const notMASchar = '-'

func getMASmissing(c rune) rune {
	switch c {
	case M:
		return S
	case S:
		return M
	default:
		return notMASchar
	}
}
