package solution

import (
	"fmt"
)

type GridMove [2]int

var (
	MoveUp        GridMove = [2]int{-1, 0}
	MoveRight     GridMove = [2]int{0, 1}
	MoveDown      GridMove = [2]int{1, 0}
	MoveLeft      GridMove = [2]int{0, -1}
	MoveUpLeft    GridMove = [2]int{-1, -1}
	MoveUpRight   GridMove = [2]int{-1, 1}
	MoveDownLeft  GridMove = [2]int{1, -1}
	MoveDownRight GridMove = [2]int{1, 1}
)

func (m GridMove) String() string {
	switch m {
	case MoveUp:
		return "Up"
	case MoveRight:
		return "Right"
	case MoveDown:
		return "Down"
	case MoveLeft:
		return "Left"
	case MoveUpLeft:
		return "UpLeft"
	case MoveUpRight:
		return "UpRight"
	case MoveDownLeft:
		return "DownLeft"
	case MoveDownRight:
		return "DownRight"
	default:
		return "Unknown"
	}
}

func (m GridMove) LeftHandSearchOrderCross() [4]GridMove {
	switch m {
	case MoveRight:
		return [4]GridMove{MoveUp, MoveRight, MoveDown, MoveLeft}
	case MoveDown:
		return [4]GridMove{MoveRight, MoveDown, MoveLeft, MoveUp}
	case MoveLeft:
		return [4]GridMove{MoveDown, MoveLeft, MoveUp, MoveRight}
	case MoveUp:
		return [4]GridMove{MoveLeft, MoveUp, MoveRight, MoveDown}
	default:
		panic(fmt.Sprintf("invalid move for left hand search order cross %v", m))
	}
}

func (m GridMove) Opposite() GridMove {
	switch m {
	case MoveUp:
		return MoveDown
	case MoveRight:
		return MoveLeft
	case MoveDown:
		return MoveUp
	case MoveLeft:
		return MoveRight
	case MoveUpLeft:
		return MoveDownRight
	case MoveUpRight:
		return MoveDownLeft
	case MoveDownLeft:
		return MoveUpRight
	case MoveDownRight:
		return MoveUpLeft
	default:
		panic(fmt.Sprintf("invalid move for opposite %v", m))
	}
}

// Row, Col Clockwise from top left
var (
	// Up, Right, Down, Left
	GridCrossMoves [4]GridMove = [4]GridMove{MoveUp, MoveRight, MoveDown, MoveLeft}
	// UpLeft, UpRight, DownLeft, DownRight
	GridDiagonalMoves [4]GridMove = [4]GridMove{
		MoveUpLeft,
		MoveUpRight,
		MoveDownLeft,
		MoveDownRight,
	}
	// UpLeft, Up, UpRight, Right, DownRight, Down, DownLeft, Left
	GridAllMoves [8]GridMove = [8]GridMove{
		MoveUpLeft,
		MoveUp,
		MoveUpRight,
		MoveRight,
		MoveDownRight,
		MoveDown,
		MoveDownLeft,
		MoveLeft,
	}
)

type Point[T any] struct {
	Row, Col int
	Value    T
}

func (p Point[T]) Move(move GridMove) Point[T] {
	return Point[T]{Row: p.Row + move[0], Col: p.Col + move[1], Value: p.Value}
}

// ReflectPoint computes the reflection (mirror image) of other
// across the current point p.
// .....
// .....
// ..p..  (p is the reflection center)
// .....
// ....o  (o is the point being reflected)
//
// p.ReflectPoint(o):
// N....   (N is the reflected point, result of the method call)
// .....
// ..p..
// .....
// ....o
func (p Point[T]) ReflectPoint(other Point[T]) Point[T] {
	return Point[T]{
		// Simplified p.Row - (other.Row - p.Row)
		Row:   2*p.Row - other.Row,
		Col:   2*p.Col - other.Col,
		Value: other.Value,
	}
}

func (p Point[T]) ToArray() [2]int {
	return [2]int{p.Row, p.Col}
}

func (p Point[T]) ToCellIndex(columns int) int {
	return p.Row*columns + p.Col
}

func (p Point[T]) EqualCoordinates(other Point[T]) bool {
	return p.Row == other.Row && p.Col == other.Col
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
	columns := len(slices[0])
	for _, row := range slices {
		if len(row) != columns {
			return nil, fmt.Errorf("rows have different lengths")
		}
	}
	rows := len(slices)
	g := NewGrid[T](rows, columns)
	for r, row := range slices {
		for c, cell := range row {
			g.Set(r, c, cell)
		}
	}
	return g, nil
}

func (g *Grid[T]) Get(row, column int) Point[T] {
	return g.Cells[row*g.Columns+column]
}

func (g *Grid[T]) GetAny(row, column int) Point[T] {
	if !g.IsInside(row, column) {
		return Point[T]{Row: row, Col: column}
	}

	return g.Get(row, column)
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

// Up, Right, Down, Left
func (g *Grid[T]) GetNeighborsCross(row, column int) []Point[T] {
	neighbors := make([]Point[T], 0, 4)

	for _, move := range GridCrossMoves {
		newRow, newCol := row+move[0], column+move[1]

		if g.IsInside(newRow, newCol) {
			n := g.Get(newRow, newCol)
			neighbors = append(neighbors, n)
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

func (g *Grid[T]) GetAnyByMove(row, column int, move GridMove) Point[T] {
	newRow, newCol := row+move[0], column+move[1]
	return g.GetAny(newRow, newCol)
}

func (g *Grid[T]) GetNeighborByMove(row, column int, move GridMove) (Point[T], bool) {
	if !g.IsInside(row, column) {
		panic(
			fmt.Sprintf(
				"point %d, %d is outside of the grid %d x %d",
				row,
				column,
				g.Rows,
				g.Columns,
			),
		)
	}
	newRow, newCol := row+move[0], column+move[1]
	return g.GetWithCheck(newRow, newCol)
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

func (g *Grid[T]) MovePoint(row, col int, move GridMove) (Point[T], bool) {
	// This is bad practice, but it's ok for our aoc-solving case
	if !g.IsInside(row, col) {
		panic(
			fmt.Sprintf(
				"point %d, %d is outside of the grid %d x %d",
				row,
				col,
				g.Rows,
				g.Columns,
			),
		)
	}
	return g.GetWithCheck(row+move[0], col+move[1])
}

func (g *Grid[T]) Copy() *Grid[T] {
	newGrid := NewGrid[T](g.Rows, g.Columns)
	copy(newGrid.Cells, g.Cells)
	return newGrid
}

func (g *Grid[T]) ToSlices() [][]T {
	slices := make([][]T, g.Rows)
	for i := range slices {
		slices[i] = make([]T, g.Columns)
	}

	for i, cell := range g.Cells {
		slices[i/g.Columns][i%g.Columns] = cell.Value
	}

	return slices
}
