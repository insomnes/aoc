package solution

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ParsedInput = []uint64

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track(time.Now(), "ParseInput")

	var numbers []uint64

	if len(lines) != 1 {
		return numbers, fmt.Errorf("expected 1 line, got %d", len(lines))
	}

	split := strings.Split(lines[0], " ")
	numbers = make([]uint64, len(split))

	for i, s := range split {
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return numbers, fmt.Errorf("parsing %s: %w", s, err)
		}
		numbers[i] = n
	}

	return numbers, nil
}

const (
	partOneBlinks = 25
	partTwoBlinks = 75
)

var blinkCache = make(map[uint64][]uint64, 4000)

func blink(n uint64) []uint64 {
	if n == 0 {
		return []uint64{1}
	}
	if val, ok := blinkCache[n]; ok {
		return val
	}
	s := strconv.FormatUint(n, 10)
	// Slow but fine for this problem
	if len(s)%2 != 0 {
		result := n * 2024
		blinkCache[n] = []uint64{result}
		return []uint64{result}
	}
	a, err := strconv.ParseUint(s[:len(s)/2], 10, 64)
	if err != nil {
		panic(err)
	}
	b, err := strconv.ParseUint(s[len(s)/2:], 10, 64)
	if err != nil {
		panic(err)
	}

	result := []uint64{a, b}
	blinkCache[n] = result

	return result
}

const distinctCap = 4000

func SimulateBlinks(initNumbers []uint64, totalBlinks int) uint64 {
	counts := make(map[uint64]uint64, 4000)
	blinkCounts := make(map[uint64]uint64, 4000)
	for _, n := range initNumbers {
		counts[n]++
	}

	for range totalBlinks {
		for num, count := range counts {
			blinkResult := blink(num)
			for _, br := range blinkResult {
				blinkCounts[br] += count
			}
		}

		// To avoid re-allocating memory
		counts, blinkCounts = blinkCounts, counts
		clear(blinkCounts)
	}

	var total uint64 = 0
	for _, count := range counts {
		total += count
	}

	return total
}

func PartOne(inp ParsedInput) uint64 {
	defer Track(time.Now(), "PartOne")
	var total uint64 = 0

	total = SimulateBlinks(inp, partOneBlinks)

	return total
}

func PartTwo(inp ParsedInput) uint64 {
	defer Track(time.Now(), "PartTwo")
	var total uint64 = 0

	total = SimulateBlinks(inp, partTwoBlinks)

	return total
}
