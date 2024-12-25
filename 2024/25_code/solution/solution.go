package solution

import (
	"fmt"
	"time"
)

const (
	fullLine  = "#####"
	emptyLine = "....."
	pins      = 5
	pinH      = 5
	// Raw input height
	rawInpHeight = 7
	full         = '#'
	empty        = '.'
)

type PinInfo [pins]int

type Lock struct {
	Idx  int
	Pins PinInfo
}

func parse(lines []string) (PinInfo, error) {
	var pinHeight PinInfo
	for _, line := range lines {
		if len(line) != pins {
			return pinHeight, fmt.Errorf("invalid line %s", line)
		}
		for i, c := range line {
			if c == full {
				pinHeight[i]++
			}
		}
	}
	return pinHeight, nil
}

func NewHeightMap() HeightMap {
	hm := HeightMap{}
	for p := 0; p < pins; p++ {
		for h := 0; h < pinH+1; h++ {
			hm[p][h] = NewSet[int](0)
		}
	}
	return hm
}

type HeightMap [pins][pinH + 1]Set[int]

type TaskInfo struct {
	MatchingKeys HeightMap
	Locks        []Lock
}

type ParsedInput = TaskInfo

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track(time.Now(), "ParseInput")
	locks := make([]Lock, 0)
	keyCount := 0
	matchingKeys := NewHeightMap()

	i := 0
	for i < len(lines) {
		line := lines[i]
		var toParse []string
		isLock := true
		toParse = lines[i+1 : i+rawInpHeight]
		if line == emptyLine { // Key
			toParse = toParse[:len(toParse)-1]
			isLock = false
		}
		i += rawInpHeight + 1 // Skip the next empty line

		pinInfo, err := parse(toParse)
		if err != nil {
			return ParsedInput{}, err
		}
		if isLock {
			locks = append(locks, Lock{Idx: len(locks), Pins: pinInfo})
			continue
		}

		for p, kh := range pinInfo {
			for h := kh; h < pinH; h++ {
				matchingKeys[p][h].Add(keyCount)
			}
		}
		keyCount++
	}
	return TaskInfo{MatchingKeys: matchingKeys, Locks: locks}, nil
}

func PartOne(inp ParsedInput) int {
	defer Track(time.Now(), "PartOne")
	total := 0

	ti := inp
	locks, matchingKeys := ti.Locks, ti.MatchingKeys

	for _, lock := range locks {
		fits := CountFits(lock, matchingKeys)
		total += fits
	}

	return total
}

func CountFits(lock Lock, matchingKeys HeightMap) int {
	keyPool := NewSet[int](0)

	for pin, height := range lock.Pins {
		if height == 0 {
			continue
		}
		pinFits := matchingKeys[pin][pinH-height]
		if keyPool.Empty() {
			keyPool = pinFits
			continue
		}
		keyPool = keyPool.Intersection(pinFits)
		if keyPool.Empty() {
			return 0
		}

	}

	return keyPool.Size()
}

func PartTwo(inp ParsedInput) int {
	defer Track(time.Now(), "PartTwo")

	return 0
}
