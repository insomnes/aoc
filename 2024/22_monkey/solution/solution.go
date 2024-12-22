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

const dwSize = 4

type DeltaWindow [dwSize]int

func (dw *DeltaWindow) Shift(delta int) {
	for i := 0; i < dwSize-1; i++ {
		dw[i] = dw[i+1]
	}
	dw[dwSize-1] = delta
}

const (
	keyBits   = 5
	posOffset = 9
	bits      = keyBits * dwSize
)

// We have numbers from -9 to 9
// To not work with negative numbers we shift them by 9
// So we have numbers from 0 to 18, we would need 5 bits to represent them
// So we would have 20 bits number to represent all deltas
func (dw *DeltaWindow) Key() int {
	key := 0
	for i := 0; i < dwSize; i++ {
		key |= ((dw[i] + posOffset) << (keyBits * (dwSize - 1 - i)))
	}
	return key
}

func decodeKey(key int) DeltaWindow {
	dw := DeltaWindow{}
	bitMask := 1<<keyBits - 1
	for i := 0; i < dwSize; i++ {
		shifted := key >> (keyBits * (dwSize - 1 - i))
		dw[i] = (shifted & bitMask) - 9
	}
	return dw
}

func calculateFirstWindow(secNum int) (DeltaWindow, int) {
	dw := DeltaWindow{}
	prevPrice := secNum % 10
	for i := 0; i < dwSize; i++ {
		secNum = nextSecretNumber(secNum)
		curPrice := secNum % 10
		delta := curPrice - prevPrice
		dw[i] = delta
		prevPrice = curPrice
	}

	return dw, secNum
}

func findDeltas(secNum int, iters int, allUniqueDeltas map[int]map[int]int) {
	origSecNum := secNum

	// Init first window
	deltaWindow, secNum := calculateFirstWindow(secNum)
	deltaKey := deltaWindow.Key()
	snPriceMap, found := allUniqueDeltas[deltaKey]
	if !found {
		snPriceMap = make(map[int]int)
	}
	// Set max price for this delta
	curPrice := secNum % 10
	snPriceMap[origSecNum] = max(snPriceMap[origSecNum], curPrice)
	allUniqueDeltas[deltaKey] = snPriceMap
	prevPrice := curPrice

	for i := dwSize + 1; i <= iters; i++ {
		secNum = nextSecretNumber(secNum)
		curPrice := secNum % 10
		delta := curPrice - prevPrice
		prevPrice = curPrice
		deltaWindow.Shift(delta)
		// Leave max price for this delta
		deltaKey = deltaWindow.Key()

		snPriceMap, found := allUniqueDeltas[deltaKey]
		if !found {
			snPriceMap = make(map[int]int)
		} else if _, found := snPriceMap[origSecNum]; found {
			continue
		}
		snPriceMap[origSecNum] = max(snPriceMap[origSecNum], curPrice)
		allUniqueDeltas[deltaKey] = snPriceMap
	}
}

const partTwoIters = 2000

func PartTwo(inp ParsedInput) int {
	defer Track(time.Now(), "PartTwo")
	// delta as a "key" -> Map{originalSn -> price}
	allUniqueDeltas := make(map[int]map[int]int)
	secretNumbers := inp

	for _, sn := range secretNumbers {
		findDeltas(sn, partTwoIters, allUniqueDeltas)
	}

	maxDeltaKey, maxBananas := -1, 0

	for deltaKey, deltaSNPrices := range allUniqueDeltas {
		deltaBananas := 0
		for _, price := range deltaSNPrices {
			deltaBananas += price
		}

		if deltaBananas > maxBananas {
			maxBananas = deltaBananas
			maxDeltaKey = deltaKey
		}

	}
	fmt.Println("Delta found", decodeKey(maxDeltaKey), "Bananas:", maxBananas)
	return maxBananas
}
