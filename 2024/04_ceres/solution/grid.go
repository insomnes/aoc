package solution

import (
	"fmt"
)

// Row, Col Clockwise from top left
var (
	GridCrossMoves    [4][2]int = [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	GridDiagonalMoves [4][2]int = [4][2]int{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
	GridAllMoves      [8][2]int = [8][2]int{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, 1},
		{1, 1},
		{1, 0},
		{1, -1},
		{0, -1},
	}
)

type Point[T any] struct {
	Row, Col int
	Value    T
}

func NewPointByIndex[T any](i, columns int, value T) Point[T] {
	return Point[T]{Row: i / columns, Col: i % columns, Value: value}
}

type Grid[T any] struct {
	Cells   []Point[T]
	Rows    int
	Columns int
}

func NewGrid[T any](rows, columns int) *Grid[T] {
	return &Grid[T]{Cells: make([]Point[T], rows*columns), Rows: rows, Columns: columns}
}

func NewGridWithDefault[T any](rows, columns int, defaultValue T) *Grid[T] {
	cells := make([]Point[T], rows*columns)
	for i := range cells {
		cell := NewPointByIndex[T](i, columns, defaultValue)
		cells[i] = cell
	}
	return &Grid[T]{Cells: cells, Rows: rows, Columns: columns}
}

func NewGridFromSlice[T any](slice []T, columns int) (*Grid[T], error) {
	rows := len(slice) / columns
	if len(slice)%columns != 0 {
		return nil, fmt.Errorf("slice length is not a multiple of columns")
	}

	cells := make([]Point[T], rows*columns)
	for i := range slice {
		cells[i] = NewPointByIndex[T](i, columns, slice[i])
	}

	return &Grid[T]{Cells: cells, Rows: rows, Columns: columns}, nil
}

func NewGridFromSlices[T any](slices [][]T) (*Grid[T], error) {
	width := len(slices[0])
	for _, row := range slices {
		if len(row) != width {
			return nil, fmt.Errorf("rows have different lengths")
		}
	}
	height := len(slices)
	g := NewGrid[T](width, height)
	for y, row := range slices {
		for x, cell := range row {
			g.Set(x, y, cell)
		}
	}
	return g, nil
}

func (g *Grid[T]) Get(row, column int) Point[T] {
	return g.Cells[row*g.Columns+column]
}

func (g *Grid[T]) GetWithCheck(row, column int) (Point[T], bool) {
	if !g.IsInside(row, column) {
		return Point[T]{}, false
	}
	return g.Get(row, column), true
}

func (g *Grid[T]) Set(row, column int, value T) {
	point := Point[T]{Row: row, Col: column, Value: value}
	g.Cells[row*g.Columns+column] = point
}

func (g *Grid[T]) SetWithCheck(point Point[T]) error {
	if !g.IsInside(point.Row, point.Col) {
		return fmt.Errorf("point %d, %d is outside of the grid", point.Row, point.Col)
	}
	g.Cells[point.Row*g.Columns+point.Col] = point
	return nil
}

func (g *Grid[T]) IsInside(row, column int) bool {
	return row >= 0 && row < g.Rows && column >= 0 && column < g.Columns
}

func (g *Grid[T]) GetNeighborsCross(row, column int) []Point[T] {
	neighbors := make([]Point[T], 0, 4)

	for _, move := range GridCrossMoves {
		newRow, newCol := row+move[0], column+move[1]
		if g.IsInside(newRow, newCol) {
			neighbors = append(neighbors, g.Get(newRow, newCol))
		}
	}

	return neighbors
}

func (g *Grid[T]) GetNeighborsDiagonal(row, column int) []Point[T] {
	neighbors := make([]Point[T], 0, 4)

	for _, move := range GridDiagonalMoves {
		newRow, newCol := row+move[0], column+move[1]
		if g.IsInside(newRow, newCol) {
			neighbors = append(neighbors, g.Get(newRow, newCol))
		}
	}

	return neighbors
}

func (g *Grid[T]) GetNeighbors(row, column int) []Point[T] {
	neighbors := make([]Point[T], 0, 8)

	for _, move := range GridAllMoves {
		newRow, newCol := row+move[0], column+move[1]
		if g.IsInside(newRow, newCol) {
			neighbors = append(neighbors, g.Get(newRow, newCol))
		}
	}

	return neighbors
}

func (g *Grid[T]) GetLeft(row, column int) (Point[T], bool) {
	if !g.IsInside(row, column-1) {
		return Point[T]{}, false
	}
	return g.Get(row, column-1), true
}

func (g *Grid[T]) GetRight(row, column int) (Point[T], bool) {
	if !g.IsInside(row, column+1) {
		return Point[T]{}, false
	}
	return g.Get(row, column+1), true
}

func (g *Grid[T]) GetUp(row, column int) (Point[T], bool) {
	if !g.IsInside(row-1, column) {
		return Point[T]{}, false
	}
	return g.Get(row-1, column), true
}

func (g *Grid[T]) GetDown(row, column int) (Point[T], bool) {
	if !g.IsInside(row+1, column) {
		return Point[T]{}, false
	}
	return g.Get(row+1, column), true
}

func (g *Grid[T]) GetUpLeft(row, column int) (Point[T], bool) {
	if !g.IsInside(row-1, column-1) {
		return Point[T]{}, false
	}
	return g.Get(row-1, column-1), true
}

func (g *Grid[T]) GetUpRight(row, column int) (Point[T], bool) {
	if !g.IsInside(row-1, column+1) {
		return Point[T]{}, false
	}
	return g.Get(row-1, column+1), true
}

func (g *Grid[T]) GetDownLeft(row, column int) (Point[T], bool) {
	if !g.IsInside(row+1, column-1) {
		return Point[T]{}, false
	}
	return g.Get(row+1, column-1), true
}

func (g *Grid[T]) GetDownRight(row, column int) (Point[T], bool) {
	if !g.IsInside(row+1, column+1) {
		return Point[T]{}, false
	}
	return g.Get(row+1, column+1), true
}

func (g *Grid[T]) Copy() *Grid[T] {
	newGrid := NewGrid[T](g.Rows, g.Columns)
	copy(newGrid.Cells, g.Cells)
	return newGrid
}
