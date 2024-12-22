package solution

import (
	"fmt"
	"strconv"
	"time"
)

type ParsedInput = []int

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track(time.Now(), "ParseInput")
	numbers := make([]int, 0, len(lines))
	for _, line := range lines {
		num, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, int(num))
	}
	return numbers, nil
}

// 1 << 24
// => 16777216

// 100000000 & ((1 << 24) - 1)
// 16113920
const (
	pruneNumber = 16777216
	pruneMask   = 0xFFFFFF
)

func mixAndPrune(a, b int) int {
	return (a ^ b) & pruneMask
}

func nextSecretNumber(sn int) int {
	// x1 = ( x0 XOR (x0 * 64) ) % 16777216
	// x1 = ( x0 XOR (x0 << 6) ) & 0xFFFFFF
	sn = mixAndPrune(sn, sn<<6)
	// x2 = ( x1 XOR (x1 DIV 32) ) % 16777216
	// x2 = ( x1 XOR (x1 >> 5) ) & 0xFFFFFF
	sn = mixAndPrune(sn, sn>>5)
	// x3 = ( x2 XOR (x2 * 2048) ) % 16777216
	// x3 = ( x2 XOR (x2 << 11) ) & 0xFFFFFF
	sn = mixAndPrune(sn, sn<<11)
	return sn
}

func runOperationsNtimes(sn int, n int) int {
	for i := 1; i <= n; i++ {
		sn = nextSecretNumber(sn)
	}
	return sn
}

const partOneIters = 2000

func PartOne(inp ParsedInput) int {
	defer Track(time.Now(), "PartOne")
	total := int(0)
	for _, num := range inp[0:1] {
		result := runOperationsNtimes(num, partOneIters)
		total += result
	}

	return total
}

type DeltaWindow struct {
	first  int
	second int
	third  int
	fourth int
}

func (dw *DeltaWindow) Shift(delta int) {
	dw.first = dw.second
	dw.second = dw.third
	dw.third = dw.fourth
	dw.fourth = delta
}

func (dw *DeltaWindow) String() string {
	return fmt.Sprintf("(%d %d %d %d)", dw.first, dw.second, dw.third, dw.fourth)
}

func findDeltas(sn int, iters int, allUniqueDeltas map[string]map[int]int) {
	originalSn := sn
	// Mappint delta string to price
	prevPrice := sn % 10

	var first, second, third, fourth int

	sn = nextSecretNumber(sn)
	curPrice := sn % 10
	delta := curPrice - prevPrice
	prevPrice = curPrice
	first = delta

	sn = nextSecretNumber(sn)
	curPrice = sn % 10
	delta = curPrice - prevPrice
	prevPrice = curPrice
	second = delta

	sn = nextSecretNumber(sn)
	curPrice = sn % 10
	delta = curPrice - prevPrice
	prevPrice = curPrice
	third = delta

	sn = nextSecretNumber(sn)
	curPrice = sn % 10
	delta = curPrice - prevPrice
	prevPrice = curPrice
	fourth = delta

	dw := DeltaWindow{first, second, third, fourth}
	deltaKey := dw.String()

	snPriceMap, found := allUniqueDeltas[deltaKey]
	if !found {
		snPriceMap = make(map[int]int)
	}
	snPriceMap[originalSn] = max(snPriceMap[originalSn], curPrice)
	allUniqueDeltas[deltaKey] = snPriceMap

	for i := 5; i <= iters; i++ {
		sn = nextSecretNumber(sn)
		curPrice := sn % 10
		delta := curPrice - prevPrice
		prevPrice = curPrice
		dw.Shift(delta)
		// Leave max price for this delta
		deltaKey = dw.String()

		snPriceMap, found := allUniqueDeltas[deltaKey]
		if !found {
			snPriceMap = make(map[int]int)
		} else {
			// Only first occurrence of delta counts
			if _, found := snPriceMap[originalSn]; found {
				continue
			}
		}
		snPriceMap[originalSn] = max(snPriceMap[originalSn], curPrice)
		allUniqueDeltas[deltaKey] = snPriceMap
	}
}

const partTwoIters = 2000

func PartTwo(inp ParsedInput) int {
	defer Track(time.Now(), "PartTwo")
	// delta -> Map{originalSn -> price}
	allUniqueDeltas := make(map[string]map[int]int)
	secretNumbers := inp

	for _, sn := range secretNumbers {
		findDeltas(sn, partTwoIters, allUniqueDeltas)
	}

	maxDelta := ""
	maxBananas := 0

	for delta, deltaSNPrices := range allUniqueDeltas {
		deltaBananas := 0
		for _, price := range deltaSNPrices {
			deltaBananas += price
		}

		if deltaBananas > maxBananas {
			maxBananas = deltaBananas
			maxDelta = delta
		}

	}
	fmt.Println("MAX DELTA", maxDelta, "MAX BANANAS", maxBananas)
	return maxBananas
}
