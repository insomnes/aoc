package solution

import (
	"fmt"
	"regexp"
)

type ParsedInput = []string

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track("ParseInput")()

	return lines, nil
}

func PartOne(inp ParsedInput) int {
	defer Track("PartOne")()
	pattern := `mul\(\d{1,3},\d{1,3}\)`
	re := regexp.MustCompile(pattern)

	sum := 0

	for _, row := range inp {
		matches := re.FindAllString(row, -1)
		for _, match := range matches {
			var a, b int
			_, err := fmt.Sscanf(match, "mul(%d,%d)", &a, &b)
			if err != nil {
				panic(err)
			}
			sum += a * b
		}
	}
	return sum
}

func PartTwo(inp ParsedInput) int {
	defer Track("PartTwo")()
	pattern := `(do\(\)|don't\(\)|mul\(\d{1,3},\d{1,3}\))`
	re := regexp.MustCompile(pattern)

	sum := 0
	active := true
	var a, b int

	for _, row := range inp {
		matches := re.FindAllString(row, -1)
		if len(matches) == 0 {
			panic("no matches")
		}

		for _, match := range matches {
			switch match {
			case "do()":
				active = true
			case "don't()":
				active = false
			default:
				if !active {
					continue
				}
				if _, err := fmt.Sscanf(match, "mul(%d,%d)", &a, &b); err != nil {
					panic(fmt.Sprintf("parsing %s: %s", match, err))
				}
				sum += a * b
			}
		}

	}
	return sum
}
