package solution

import (
	"fmt"
	"slices"
)

type ParsedInput = [][2]int

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
		var a, b int
		n, err := fmt.Sscan(l, &a, &b)
		if n != 2 || err != nil {
			return nil, fmt.Errorf("parse input line %v", l)
		}
		inp = append(inp, [2]int{a, b})
	}
	return inp, nil
}

func PartOne(inp ParsedInput) int {
	defer Track("PartOne")()

	inpLen := len(inp)
	lList := make([]int, 0, inpLen)
	rList := make([]int, 0, inpLen)

	for _, p := range inp {
		left, right := p[0], p[1]
		lList = append(lList, left)
		rList = append(rList, right)
	}

	slices.Sort(lList)
	slices.Sort(rList)

	res := 0

	for i, n := range lList {
		res += abs(n - rList[i])
	}

	return res
}

func PartTwo(inp ParsedInput) int {
	defer Track("PartTwo")()

	inpLen := len(inp)
	lList := make([]int, 0, inpLen)
	lCount := make(map[int]int, inpLen)

	rList := make([]int, 0, inpLen)
	rCount := make(map[int]int, inpLen)

	res := 0

	for _, p := range inp {
		left, right := p[0], p[1]

		lList = append(lList, left)
		lCount[left]++

		rList = append(rList, right)
		rCount[right]++

		res += left * rCount[left]

		res += right * lCount[right]
		if left == right {
			res -= left
		}

	}

	return res
}
