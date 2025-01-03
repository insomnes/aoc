package solution

import (
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
	total := 0
	for _, num := range inp {
		result := runOperationsNtimes(num, partOneIters)
		total += result
	}

	return total
}

const partTwoIters = 2000

// We want to count only prices for each delta and only once for secret number
// For delta identification we use 20 bits number (5 bits for each delta part)
// Each delta is the product of 5 prices (price window). Some price windows lead
// to the same delta. We precalculate Price Window -> Delta key mapping.
// For each price window we "count" the price which leads to the delta, and delta
// will not calculate twice for the same original secret number.
func PartTwo(inp ParsedInput) int {
	defer Track(time.Now(), "PartTwo")
	maxBananas := 0
	secretNumbers := inp

	// Mapping of delta key to delta counts pre-calculated here
	allCounts := NewAllCounts(len(secretNumbers))
	for _, sn := range secretNumbers {
		findDeltasCountForSN(sn, partTwoIters, allCounts)
		allCounts.Iter++
	}

	maxBananas = allCounts.Max

	return maxBananas
}

func findDeltasCountForSN(secNum int, iters int, allCounts *AllCounts) {
	origSecNum := secNum
	priceWindow, secNum := prepareFirstPriceWindow(secNum)
	allCounts.Count(priceWindow, origSecNum)

	for i := deltaSize + 1; i <= iters; i++ {
		secNum = nextSecretNumber(secNum)
		curPrice := secNum % 10
		priceWindow = shiftPriceWindow(priceWindow, curPrice)

		// Here we count the price for each delta only once
		// per original secret number
		allCounts.Count(priceWindow, origSecNum)
	}
}

type DeltaCounts struct {
	LastOrig int
	Sum      int
	MaxPrice int
	Key      int
}

func NewDeltaCounts(key int) *DeltaCounts {
	return &DeltaCounts{
		LastOrig: -1,
		Sum:      0,
		MaxPrice: 0,
		Key:      key,
	}
}

func (dc *DeltaCounts) Count(price int, origSecNum int) int {
	// Do not count same secret number twice
	if origSecNum == dc.LastOrig {
		return -1
	}

	dc.Sum += price
	dc.LastOrig = origSecNum

	return dc.Sum
}

func (dc *DeltaCounts) AddPrice(price int) {
	if price > dc.MaxPrice {
		dc.MaxPrice = price
	}
}

type AllCounts struct {
	// Price window as 5 numbers to delta key
	DeltaCounts [100_000]*DeltaCounts
	Max         int
	Iter        int
	total       int
}

func NewAllCounts(total int) *AllCounts {
	mapping := [100_000]*DeltaCounts{}
	prepared := make(map[int]*DeltaCounts, 55_000)
	for i := 0; i < 100_000; i++ {
		key := calcDeltaKey(i)
		dc, found := prepared[key]
		if !found {
			dc = NewDeltaCounts(key)
			prepared[key] = dc
		}
		dc.AddPrice(i % 10)
		mapping[i] = dc
	}
	return &AllCounts{
		DeltaCounts: mapping,
		Max:         0,
		Iter:        0,
		total:       total,
	}
}

func (ac *AllCounts) Count(priceWindow, origSecNum int) {
	dc := ac.DeltaCounts[priceWindow]
	if dc == nil {
		return
	}
	price := priceWindow % 10
	deltaGainUpperBound := (ac.total - ac.Iter) * dc.MaxPrice
	if dc.Sum+deltaGainUpperBound < ac.Max {
		ac.DeltaCounts[priceWindow] = nil
		return
	}

	sum := dc.Count(price, origSecNum)
	if sum == -1 {
		return
	}

	if sum > ac.Max {
		ac.Max = sum
		return
	}
}

const (
	deltaSize      = 4
	keyBits        = 5
	positiveOffset = 9
	bits           = keyBits * deltaSize
	initDelim      = 10000
)

// Unique key for delta from price window
// We have price window as 5 numbers: 96428
// Delta is 4 numbers: (6-9, 4-6, 2-4, 8-2) = (-3, -2, -2, 6)
// So we can have numbers from -9 to 9 in delta
// To not work with negative numbers we shift them by 9
// So we have numbers from 0 to 18, we would need 5 bits to represent them
// So we would have 20 bits number to represent all deltas
func calcDeltaKey(priceWindow int) int {
	delim := initDelim
	key := 0
	for i := 0; i < deltaSize; i++ {
		prev := priceWindow / delim
		priceWindow %= delim
		delim /= 10
		cur := priceWindow / delim
		delta := cur - prev
		key |= ((delta + positiveOffset) << (keyBits * (deltaSize - 1 - i)))
	}
	return key
}

// For debug
type DeltaWindow [deltaSize]int

func decodeKey(key int) DeltaWindow {
	dw := DeltaWindow{}
	bitMask := 1<<keyBits - 1
	for i := 0; i < deltaSize; i++ {
		shifted := key >> (keyBits * (deltaSize - 1 - i))
		dw[i] = (shifted & bitMask) - 9
	}
	return dw
}

func shiftPriceWindow(pw, next int) int {
	return (pw%10000)*10 + next
}

func prepareFirstPriceWindow(secNum int) (int, int) {
	pw := secNum % 10
	for i := 1; i < 5; i++ {
		secNum = nextSecretNumber(secNum)
		pw = pw*10 + secNum%10
	}
	return pw, secNum
}
