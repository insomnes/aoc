package solution

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

const (
	testRows  = 7
	testCols  = 11
	inputRows = 103
	inputCols = 101
)

type Robot struct {
	Row, Col   int
	VelR, VelC int
}

func (r Robot) Simulate(modRow int, modCol int, ticks int) Robot {
	if ticks <= 0 {
		panic(fmt.Sprintf("Invalid number of ticks: %d", ticks))
	}

	rowDelta := calculateDelta(r.VelR, ticks, modRow)
	newRow := modAdd(r.Row, rowDelta, modRow)

	colDelta := calculateDelta(r.VelC, ticks, modCol)
	newCol := modAdd(r.Col, colDelta, modCol)

	return Robot{Row: newRow, Col: newCol, VelR: r.VelR, VelC: r.VelC}
}

func (r Robot) String() string {
	return fmt.Sprintf(
		"Robot(%d, %d -> %d,%d)",
		r.Row,
		r.Col,
		r.VelR,
		r.VelC,
	)
}

func (r Robot) ToKey() int {
	return toKey(r.Row, r.Col)
}

func toKey(row, col int) int {
	return row*100000 + col
}

func calculateDelta(velocity int, ticks int, mod int) int {
	delta := velocity * ticks % mod
	if delta < 0 {
		delta += mod
	}
	return delta
}

func modAdd(a int, b int, mod int) int {
	res := (a + b) % mod
	if res < 0 {
		res += mod
	}
	return res
}

func ParseRobot(line string) (Robot, error) {
	var row, col, vRow, vCol int
	// p=x,y v=vx,vy
	inputFmt := "p=%d,%d v=%d,%d"
	// x is the column, y is the row
	_, err := fmt.Sscanf(line, inputFmt, &col, &row, &vCol, &vRow)
	if err != nil {
		return Robot{}, err
	}

	return Robot{Row: row, Col: col, VelR: vRow, VelC: vCol}, nil
}

type ParsedInput = []Robot

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track(time.Now(), "ParseInput")

	parsed := make([]Robot, 0, len(lines))
	for _, row := range lines {
		robot, err := ParseRobot(row)
		if err != nil {
			return ParsedInput{}, err
		}
		parsed = append(parsed, robot)
	}

	return parsed, nil
}

const partOneSteps = 100

func PartOne(inp ParsedInput) int {
	defer Track(time.Now(), "PartOne")
	rows, cols := inputRows, inputCols

	robots := slices.Clone(inp)
	simulated := make([]Robot, len(robots))
	for i, robot := range robots {
		simulated[i] = robot.Simulate(rows, cols, partOneSteps)
	}

	sf := safetyFactor(rows, cols, simulated)

	return sf
}

func getQuadrant(midRow, midCol int, r Robot) int {
	if r.Row == midRow || r.Col == midCol {
		panic(fmt.Sprintf("middle row or column: %s", r))
	}
	if r.Row < midRow {
		if r.Col < midCol {
			return 0
		}
		return 1
	}
	if r.Col < midCol {
		return 2
	}
	return 3
}

func safetyFactor(rows, cols int, robots []Robot) int {
	// Odd number by task rules, so to find the middle we can just divide by 2
	midRow, midCol := rows/2, cols/2
	quadrants := make([]int, 4)
	for _, r := range robots {
		if r.Row == midRow || r.Col == midCol {
			continue
		}
		q := getQuadrant(midRow, midCol, r)
		quadrants[q] += 1
	}
	sf := 1
	for _, q := range quadrants {
		sf *= q
	}

	return sf
}

const (
	trunkSeqSize = 10
	printTree    = false
)

func PartTwo(inp ParsedInput) int {
	defer Track(time.Now(), "PartTwo")
	rows, cols := inputRows, inputCols

	robots := slices.Clone(inp)
	for sec := 1; sec <= 100000; sec++ {
		robotsByCol := make([][]int, cols)
		for i, robot := range robots {
			newRobot := robot.Simulate(rows, cols, 1)
			robotsByCol[newRobot.Col] = append(robotsByCol[newRobot.Col], newRobot.Row)
			robots[i] = newRobot
		}

		for c, column := range robotsByCol {
			slices.Sort(column)
			if !findSequence(column, trunkSeqSize) {
				continue
			}
			if printTree {
				fmt.Println(StringForTree(robots, rows, cols))
				fmt.Printf("\n%d sec, found %d in column %d:\n%v\n", sec, trunkSeqSize, c, column)
			}
			return sec
		}

	}

	return 0
}

func findSequence(toCheck []int, size int) bool {
	if len(toCheck) < size {
		return false
	}
	foundSize := 1
	previous := toCheck[0]

	for i := 1; i < len(toCheck); i++ {
		if i+size >= len(toCheck) {
			return false
		}
		num := toCheck[i]
		if num == previous {
			continue
		}
		if num == previous+1 {
			foundSize++
		} else {
			foundSize = 1
		}
		previous = num
		if foundSize == size {
			return true
		}
	}

	return false
}

func StringForTree(robots []Robot, rows, cols int) string {
	var sb strings.Builder

	counts := make(map[int]int8)
	for _, robot := range robots {
		key := robot.ToKey()
		counts[key] += 1
	}

	for r := range rows {
		for c := range cols {
			if val, ok := counts[toKey(r, c)]; ok && val > 0 {
				if val > 9 {
					panic(fmt.Sprintf("Too many robots at %d,%d: %d", r, c, val))
				}
				sb.WriteRune('#')
				continue
			}
			sb.WriteRune('.')
		}
		sb.WriteRune('\n')
	}

	return sb.String()
}

func StringRobots(robots []Robot, rows, cols int) string {
	var sb strings.Builder

	counts := make(map[int]int8)
	for _, robot := range robots {
		key := robot.ToKey()
		counts[key] += 1
	}

	for r := range rows {
		for c := range cols {
			if val, ok := counts[toKey(r, c)]; ok && val > 0 {
				if val > 9 {
					panic(fmt.Sprintf("Too many robots at %d,%d: %d", r, c, val))
				}
				sb.WriteString(fmt.Sprintf("%d", val))
				continue
			}
			sb.WriteRune('.')
		}
		sb.WriteRune('\n')
	}

	return sb.String()
}
