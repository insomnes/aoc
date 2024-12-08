package solution

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Equation struct {
	Result  uint64
	Numbers []uint64
}

func (e Equation) String() string {
	var sb strings.Builder
	sb.WriteString(strconv.FormatUint(e.Result, 10))
	sb.WriteString(" = ")
	for i, n := range e.Numbers {
		if i > 0 {
			sb.WriteString(" ? ")
		}
		sb.WriteString(strconv.FormatUint(n, 10))
	}
	return sb.String()
}

func EquationFromString(s string) (Equation, error) {
	colonIndex := strings.Index(s, ":")
	if colonIndex == -1 {
		return Equation{}, fmt.Errorf("no colon found")
	}
	result, err := strconv.ParseUint(s[:colonIndex], 10, 64)
	if err != nil {
		return Equation{}, fmt.Errorf("parse result: %w", err)
	}

	raw := strings.Split(s[colonIndex+2:], " ")
	numbers := make([]uint64, len(raw))
	for i, n := range raw {
		num, err := strconv.ParseUint(n, 10, 64)
		if err != nil {
			return Equation{}, fmt.Errorf("parse number: %w", err)
		}
		numbers[i] = num
	}

	if len(numbers) < 2 {
		return Equation{}, fmt.Errorf("not enough numbers")
	}

	return Equation{Result: result, Numbers: numbers}, nil
}

type ParsedInput = []Equation

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track(time.Now(), "ParseInput")

	inp := make(ParsedInput, len(lines))
	for i, line := range lines {
		eq, err := EquationFromString(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", i, err)
		}
		inp[i] = eq
	}
	return inp, nil
}

type State struct {
	Current   uint64
	NextI     int
	Operation Operation
}

type Operation uint8

const (
	_ Operation = iota
	OpAdd
	OpMul
	OpConcat
)

func (o Operation) ApplyBackward(result, next uint64) (uint64, bool) {
	// Based on operation we can cut our backward search early
	// For example, if we know that we are looking for a multiplication
	// we can stop if the result is not divisible by the next number
	switch o {
	case OpAdd:
		return result - next, true
	case OpMul:
		return result / next, result%next == 0
	case OpConcat:
		return UnConcat(result, next)
	default:
		panic(fmt.Sprintf("unknown operation: %d", o))
	}
}

func UnConcat(a, b uint64) (uint64, bool) {
	for b > 0 {
		if a%10 != b%10 {
			return 0, false
		}
		a /= 10
		b /= 10
	}

	return a, true
}

func prepareNextStates(current uint64, nextI int, ops []Operation) []State {
	states := make([]State, len(ops))
	for i, o := range ops {
		states[i] = State{
			Current:   current,
			NextI:     nextI,
			Operation: o,
		}
	}
	return states
}

// Operation order is important for our LIFO queue, cause we want to check
// multiplication first to cut the whole branch early
var PartOneOperations = []Operation{OpAdd, OpMul}

func ProcessState(s State, next uint64, ops []Operation) ([]State, bool) {
	if s.NextI == 0 {
		if s.Current == next {
			return nil, true
		}
		return nil, false
	}
	opResult, ok := s.Operation.ApplyBackward(s.Current, next)
	if !ok {
		return nil, false
	}
	return prepareNextStates(opResult, s.NextI-1, ops), false
}

func IsSolvable(eq Equation, operations []Operation) bool {
	lastI := len(eq.Numbers) - 1
	queue := prepareNextStates(eq.Result, lastI, operations)

	for len(queue) > 0 {
		s := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		newStates, solved := ProcessState(s, eq.Numbers[s.NextI], operations)
		if solved {
			return true
		}
		queue = append(queue, newStates...)
	}

	return false
}

func PartOne(inp ParsedInput) uint64 {
	defer Track(time.Now(), "PartOne")

	var res uint64 = 0
	for _, eq := range inp {
		if IsSolvable(eq, PartOneOperations) {
			res += eq.Result
		}
	}

	return res
}

// Same logic as in PartOne, but with concatenation operation first
var PartTwoOperations = []Operation{OpAdd, OpMul, OpConcat}

func PartTwo(inp ParsedInput) uint64 {
	defer Track(time.Now(), "PartTwo")

	var res uint64 = 0
	for _, eq := range inp {
		if IsSolvable(eq, PartTwoOperations) {
			res += eq.Result
		}
	}

	return res
}
