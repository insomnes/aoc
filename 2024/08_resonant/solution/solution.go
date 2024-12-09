package solution

import (
	"fmt"
	"time"
)

type (
	Antena = Point[rune]
)

type AntenaMap struct {
	Antenas map[rune][]Antena
	Rows    int
	Cols    int
}

type ParsedInput = AntenaMap

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track(time.Now(), "ParseInput")
	antenas := make(map[rune][]Point[rune])
	rows := len(lines)
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
		for c, char := range row {
			if char == EmptyCell {
				continue
			}
			antenas[char] = append(antenas[char], Antena{Row: r, Col: c, Value: char})
		}
	}

	return AntenaMap{Antenas: antenas, Rows: rows, Cols: cols}, nil
}

const (
	EmptyCell rune = '.'
)

func shouldCountPoint(p Point[rune], grid *Grid[bool]) bool {
	point, inside := grid.GetWithCheck(p.Row, p.Col)
	present := point.Value
	if inside && !present {
		grid.Set(p.Row, p.Col, true)
		return true
	}
	return false
}

func calculateAntiNodes(a, b Antena, grid *Grid[bool]) int {
	// Reflect the antenas by each other, to get the antinodes.
	// "Reflection" is: a.Row - (b.Row - a.Row), a.Col - (b.Col - a.Col)
	antiA, antiB := a.ReflectPoint(b), b.ReflectPoint(a)

	antiCount := 0
	if shouldCountPoint(antiA, grid) {
		antiCount++
	}
	if shouldCountPoint(antiB, grid) {
		antiCount++
	}
	return antiCount
}

func PartOne(inp ParsedInput) int {
	defer Track(time.Now(), "PartOne")
	antenas := inp.Antenas
	grid := NewGridWithDefault(inp.Rows, inp.Cols, false)
	total := 0

	for _, points := range antenas {
		for i := 0; i < len(points)-1; i++ {
			for j := i + 1; j < len(points); j++ {
				a, b := points[i], points[j]
				total += calculateAntiNodes(a, b, grid)
			}
		}
	}

	return total
}

func countAntiNodes(antiNode Point[rune], refPoint Point[rune], grid *Grid[bool]) int {
	count := 0
	// We draw a "ray" from the antinode to the edge of the grid and proceed to "step"
	// with the reflection, each step inside the grid is a valid antinode.
	for grid.IsInside(antiNode.Row, antiNode.Col) {
		if shouldCountPoint(antiNode, grid) {
			count++
		}
		antiNode, refPoint = antiNode.ReflectPoint(refPoint), antiNode

	}
	return count
}

func calculateResonantAntinodes(a, b Antena, grid *Grid[bool]) int {
	// In this part point of the other antena counts as antinode too, so we set
	// respective starting antinode to each antena
	antiCount := 0
	antiCount += countAntiNodes(a, b, grid)
	antiCount += countAntiNodes(b, a, grid)
	return antiCount
}

func PartTwo(inp ParsedInput) int {
	defer Track(time.Now(), "PartTwo")
	antenas := inp.Antenas
	grid := NewGridWithDefault(inp.Rows, inp.Cols, false)
	total := 0

	for _, points := range antenas {
		for i := 0; i < len(points)-1; i++ {
			for j := i + 1; j < len(points); j++ {
				a, b := points[i], points[j]
				total += calculateResonantAntinodes(a, b, grid)
			}
		}
	}

	return total
}
