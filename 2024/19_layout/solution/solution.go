package solution

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type TowelInfo struct {
	Designs  []string
	Patterns []string
	ByFirst  map[byte][]string
}

type ParsedInput = TowelInfo

func sortByLen(values []string) {
	sort.Slice(values, func(i, j int) bool {
		return len(values[i]) < len(values[j])
	})
}

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track(time.Now(), "ParseInput")

	patterns := strings.Split(lines[0], ", ")
	sortByLen(patterns)

	byFirst := make(map[byte][]string)

	for _, p := range patterns {
		byFirst[p[0]] = append(byFirst[p[0]], p)
	}

	if lines[1] != "" {
		return TowelInfo{}, fmt.Errorf("Invalid input line: %s", lines[2])
	}

	md := 0
	designs := lines[2:]
	for _, d := range designs {
		if len(d) > md {
			md = len(d)
		}
	}

	return TowelInfo{
		Designs:  designs,
		Patterns: patterns,
		ByFirst:  byFirst,
	}, nil
}

func getPatterns(line string) []string {
	patterns := strings.Split(line, ", ")
	return patterns
}

func matchAll(design string, byFirst map[byte][]string, cache map[string]bool) bool {
	if len(design) == 0 {
		return true
	}
	if v, ok := cache[design]; ok {
		return v
	}
	patterns := byFirst[design[0]]
	if len(patterns) == 0 {
		return false
	}

	for _, pattern := range patterns {
		if len(pattern) > len(design) {
			break
		}
		if design[:len(pattern)] != pattern {
			continue
		}
		if matchAll(design[len(pattern):], byFirst, cache) {
			cache[design] = true
			return true
		}
	}
	cache[design] = false
	return false
}

func PartOne(inp ParsedInput) int {
	defer Track(time.Now(), "PartOne")
	total := 0
	designs := inp.Designs
	byFirst := inp.ByFirst
	cache := make(map[string]bool, 4000)
	for _, design := range designs {
		if matchAll(design, byFirst, cache) {
			total++
		}
	}

	return total
}

func matchAllCount(
	design string,
	byFirst map[byte][]string,
	cache map[string]uint64,
) uint64 {
	if len(design) == 0 {
		return 1
	}

	if v, ok := cache[design]; ok {
		return v
	}

	patterns := byFirst[design[0]]
	if len(patterns) == 0 {
		return 0
	}

	var total uint64 = 0
	for _, pattern := range patterns {
		if len(pattern) > len(design) || design[:len(pattern)] != pattern {
			continue
		}
		total += matchAllCount(design[len(pattern):], byFirst, cache)
	}
	cache[design] = total
	return total
}

func PartTwo(inp ParsedInput) uint64 {
	defer Track(time.Now(), "PartTwo")

	var total uint64 = 0
	designs, byFirst := inp.Designs, inp.ByFirst
	cache := make(map[string]uint64, 11000)
	for _, design := range designs {
		cnt := matchAllCount(design, byFirst, cache)
		total += cnt
	}

	return total
}
