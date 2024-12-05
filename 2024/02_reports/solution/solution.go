package solution

import (
	"fmt"
	"strconv"
	"strings"
)

type ParsedInput = [][]int

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track("ParseInput")()
	inp := make(ParsedInput, 0, len(lines))
	for _, l := range lines {
		split := strings.Split(l, " ")
		nums := make([]int, len(split))
		for i, s := range split {
			n, err := strconv.Atoi(s)
			if err != nil {
				return nil, fmt.Errorf("convert %s to int: %w", s, err)
			}
			nums[i] = n
		}
		inp = append(inp, nums)
	}
	return inp, nil
}

func isDeltaValid(delta int) bool {
	return abs(delta) >= 1 && abs(delta) <= 3
}

func areDeltasSafe(delta, prevDelta int) bool {
	return delta*prevDelta >= 0 && isDeltaValid(delta)
}

func isRowSafe(row []int) bool {
	prev := row[0]
	prevDelta := 0
	for _, n := range row[1:] {
		delta := n - prev
		if !areDeltasSafe(delta, prevDelta) {
			return false
		}

		prev, prevDelta = n, delta
	}
	return true
}

func PartOne(inp ParsedInput) int {
	defer Track("PartOne")()

	safe := 0
	for _, row := range inp {
		if isRowSafe(row) {
			safe++
		}
	}
	return safe
}

type State struct {
	fv   int
	pi   int
	ci   int
	ni   int
	skip bool
}

func (s *State) isSafe(row []int) bool {
	if s.ci >= len(row) {
		return true
	}

	delta := row[s.pi] - row[s.ci]
	if !isDeltaValid(delta) {
		return false
	}

	if (s.pi == 0 && row[s.pi] == s.fv) || (s.pi == 1 && row[s.pi] == s.fv) {
		return true
	}

	trend := s.fv - row[s.pi]

	if trend*delta < 0 {
		return false
	}

	return true
}

func (s *State) next() {
	s.pi = s.ci
	s.ci = s.ni
	s.ni++
}

func (s *State) createSkip() State {
	return State{
		fv:   s.fv,
		pi:   s.pi,
		ci:   s.ni,
		ni:   s.ni + 1,
		skip: true,
	}
}

type Stack []State

func (s *Stack) Push(state State) {
	*s = append(*s, state)
}

func (s *Stack) Pop() State {
	state := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return state
}

func IsRowSafeWithSkip(row []int) bool {
	queue := make(Stack, 0, len(row)*2)
	// LIFO
	queue.Push(State{fv: row[1], pi: 1, ci: 2, ni: 3, skip: true})
	queue.Push(State{fv: row[0], pi: 0, ci: 1, ni: 2, skip: false})

	for len(queue) > 0 {
		state := queue.Pop()

		safe := state.isSafe(row)
		// LIFO
		if !state.skip {
			queue.Push(state.createSkip())
		}

		if safe {
			if state.ci >= len(row) {
				return true
			}
			state.next()
			queue.Push(state)
		}

	}
	return false
}

func PartTwo(inp ParsedInput) int {
	defer Track("PartTwo")()

	safe := 0
	for _, row := range inp {
		if IsRowSafeWithSkip(row) {
			safe++
		}
	}
	return safe
}
