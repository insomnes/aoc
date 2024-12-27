package solution

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

type GateType int8

func (gt GateType) String() string {
	switch gt {
	case AND:
		return "&&"
	case OR:
		return "||"
	case XOR:
		return "^^"
	default:
		return "Unknown"
	}
}

func ParseType(s string) (GateType, error) {
	if s == "AND" {
		return AND, nil
	}
	if s == "OR" {
		return OR, nil
	}
	if s == "XOR" {
		return XOR, nil
	}

	return -1, fmt.Errorf("bad gate type: %s", s)
}

const (
	Unknown GateType = iota
	AND
	OR
	XOR
)

const OutGateMarker = 'z'

type (
	WireName   = string
	WireVal    = uint64
	WireValues = map[WireName]WireVal
)

type Gate struct {
	Left  WireName
	Right WireName
	Out   WireName
	Type  GateType
}

var UnknownGate = Gate{"", "", "", Unknown}

func (g Gate) String() string {
	return fmt.Sprintf("%s %s %s -> %s", g.Left, g.Type, g.Right, g.Out)
}

func (g Gate) IsINGate() bool {
	return (g.Left[0] == 'x' && g.Right[0] == 'y') || (g.Left[0] == 'y' && g.Right[0] == 'x')
}

func (g Gate) Eval(left, right WireVal) WireVal {
	switch g.Type {
	case AND:
		return left & right
	case OR:
		return left | right
	case XOR:
		return left ^ right
	default:
		panic(fmt.Sprintf("Unknown gate type: %d of %s", g.Type, g))
	}
}

func (g Gate) EvalBit(left, right uint8) uint8 {
	switch g.Type {
	case AND:
		return left & right
	case OR:
		return left | right
	case XOR:
		return left ^ right
	default:
		panic(fmt.Sprintf("Unknown gate type: %d of %s", g.Type, g))
	}
}

// fcw AND hrn -> jjw
func ParseGate(s string) (Gate, error) {
	split := strings.Split(s, " ")
	if len(split) != 5 || split[3] != "->" {
		return Gate{}, fmt.Errorf("bad gate line: %s", s)
	}
	gt, err := ParseType(split[1])
	if err != nil {
		return Gate{}, err
	}

	left, right, out := split[0], split[2], split[4]
	if right < left {
		left, right = right, left
	}
	return Gate{left, right, out, gt}, nil
}

// y04: 1

type DeviceInfo struct {
	AllGates     []Gate
	InitWires    WireValues
	SetInitWires map[WireName]struct{}
	OutMap       map[WireName]Gate
	InMap        map[WireName][]Gate
	OutWires     []WireName
	AllWires     map[WireName]struct{}
}

func (di DeviceInfo) CreateDevice() Device {
	return NewDevice(di.OutWires, di.InitWires, di.OutMap)
}

func (di DeviceInfo) CreateAdderDevice() *AdderDevice {
	adder := NewAdderDevice(di.OutMap, len(di.OutWires))
	adder.SetInputsByWires(di.SetInitWires)
	return adder
}

func (di DeviceInfo) CreateAndSet(x, y uint64) Device {
	initWires := make(WireValues)
	for i := 0; i < len(di.OutWires); i++ {
		val := uint64(0)
		if x&1 == 1 {
			val = 1
		}
		initWires[wireName("x", i)] = val
		x >>= 1

		val = 0
		if y&1 == 1 {
			val = 1
		}
		initWires[wireName("y", i)] = val
		y >>= 1

	}
	return NewDevice(di.OutWires, initWires, di.OutMap)
}

type ParsedInput = DeviceInfo

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track(time.Now(), "ParseInput")
	initWires := make(WireValues, 0)
	setInitWires := make(map[WireName]struct{})
	allWires := make(map[WireName]struct{})
	i := 0
	for {
		line := lines[i]
		i++
		if line == "" {
			break
		}
		if i == len(lines) {
			panic("Wrong input")
		}

		// x04: 1
		split := strings.Split(line, ": ")
		if len(split) != 2 || len(split[1]) != 1 {
			return ParsedInput{}, fmt.Errorf("bad init wire line: %s", line)
		}
		wireName, valStr := split[0], split[1]
		val := uint64(valStr[0] - '0')
		initWires[wireName] = val
		if val == 1 {
			setInitWires[wireName] = struct{}{}
		}

		allWires[wireName] = struct{}{}
	}

	gates := make([]Gate, 0, len(lines[i:]))
	outWires := make([]WireName, 0)
	inMap := make(map[WireName][]Gate)
	outMap := make(map[WireName]Gate)

	for i < len(lines) {
		line := lines[i]
		i++
		gate, err := ParseGate(line)
		if err != nil {
			return ParsedInput{}, err
		}
		gates = append(gates, gate)
		if gate.Out[0] == OutGateMarker {
			outWires = append(outWires, gate.Out)
		}
		inMap[gate.Left] = append(inMap[gate.Left], gate)
		inMap[gate.Right] = append(inMap[gate.Right], gate)
		if _, ok := outMap[gate.Out]; ok {
			panic(fmt.Sprintf("Duplicate output %s", gate.Out))
		}
		outMap[gate.Out] = gate
		allWires[gate.Left] = struct{}{}
		allWires[gate.Right] = struct{}{}
		allWires[gate.Out] = struct{}{}
	}

	return DeviceInfo{
		AllGates:     gates,
		InitWires:    initWires,
		SetInitWires: setInitWires,
		OutMap:       outMap,
		InMap:        inMap,
		OutWires:     outWires,
		AllWires:     allWires,
	}, nil
}

type Device struct {
	OutWires []WireName
	Values   WireValues
	OutMap   map[WireName]Gate
	run      bool
}

func NewDevice(
	outWires []WireName,
	initWires WireValues,
	outMap map[WireName]Gate,
) Device {
	values := make(WireValues)
	for n, v := range initWires {
		values[n] = v
	}

	return Device{OutWires: outWires, Values: values, OutMap: outMap}
}

func (d *Device) Run() {
	if d.run {
		return
	}

	for i := 0; i < len(d.OutWires); i++ {
		wn := fmt.Sprintf("z%02d", i)
		d.EnsureValue(wn)
	}

	d.run = true
}

func (d *Device) RunGate(gate Gate) {
	if _, set := d.Values[gate.Out]; set {
		fmt.Println("Already set", gate)
		return
	}

	left, right := d.GetGateInput(gate)
	result := gate.Eval(left, right)
	d.Values[gate.Out] = result
}

func (d *Device) GetGateInput(gate Gate) (WireVal, WireVal) {
	d.EnsureValue(gate.Left)
	d.EnsureValue(gate.Right)
	return d.Values[gate.Left], d.Values[gate.Right]
}

func (d *Device) EnsureValue(wn WireName) {
	if _, set := d.Values[wn]; set {
		return
	}
	gate, ok := d.OutMap[wn]
	if !ok {
		panic(fmt.Sprintf("No gate for %s", wn))
	}
	d.RunGate(gate)
	if _, set := d.Values[gate.Out]; !set {
		panic(fmt.Sprintf("Gate %s not set after running", gate.Out))
	}
}

func (d *Device) CollectOutput() uint64 {
	if !d.run {
		d.Run()
	}

	result := uint64(0)
	for i := range d.OutWires {
		wn := fmt.Sprintf("z%02d", i)
		val, set := d.Values[wn]
		if !set {
			panic(fmt.Sprintf("Wire %s not set", wn))
		}
		result |= val << uint64(i)

	}
	return result
}

func PartOne(inp ParsedInput) uint64 {
	defer Track(time.Now(), "PartOne")

	devInfo := inp
	device := devInfo.CreateDevice()
	out := device.CollectOutput()

	return out
}

func PartTwo(inp ParsedInput) int {
	defer Track(time.Now(), "PartTwo")

	devInfo := inp
	adder := devInfo.CreateAdderDevice()
	wrongRes := adder.CollectOutput()
	badBits := SearchBadBits(adder)
	gates := make([]string, 0, 2*len(badBits))
	for _, bit := range badBits {
		swap := FindSwapForBit(bit, adder)
		g1, g2 := swapAdderGates(swap, bit, adder, &devInfo)
		gates = append(gates, g1, g2)
	}
	slices.Sort(gates)
	x, y := adder.Inputs()
	expected := x + y
	withSwap := adder.CollectOutput()
	fmt.Println()
	fmt.Printf(" Bad bits:  %v\n", badBits)
	fmt.Printf("        X:  %d\n", x)
	fmt.Printf("        Y:  %d\n", y)
	fmt.Printf("   Broken:  %d\n", wrongRes)
	fmt.Printf(" Expected:  %d\n", expected)
	fmt.Printf("With swap:  %d\n", withSwap)
	fmt.Printf("----------\n")
	fmt.Printf("    Gates:  %s\n\n", strings.Join(gates, ","))
	if withSwap != expected {
		panic("Not equal")
	}

	return 0
}

func swapAdderGates(swap gswp, bit int, adder *AdderDevice, devInfo *DeviceInfo) (string, string) {
	switch swap {
	case AxB:
		return swapGatesAxB(adder, bit, devInfo)
	case BxS:
		return swapGatesBxS(adder, bit, devInfo)
	case DxS:
		return swapGatesDxS(adder, bit, devInfo)
	case SxC:
		return swapGatesSxC(adder, bit, devInfo)
	default:
		panic("Unknown swap")
	}
}

type AdderDevice struct {
	X, Y             []uint8
	xNms, yNms, zNms []string
	nmsToIdx         map[string]int

	gatesByOut map[WireName]Gate
	outValues  map[WireName]uint8

	bitLen int
}

func NewAdderDevice(
	gatesByOut map[WireName]Gate,
	bitLen int,
) *AdderDevice {
	x, y := make([]uint8, bitLen), make([]uint8, bitLen)
	xNms, yNms := createToNameTranslation("x", bitLen), createToNameTranslation("y", bitLen)
	zNms := createToNameTranslation("z", bitLen)
	nmsToIdx := createFromNameTranslation([]string{"x", "y", "z"}, bitLen)

	outValues := make(map[WireName]uint8)

	return &AdderDevice{
		X: x, Y: y,
		xNms: xNms, yNms: yNms, zNms: zNms,
		nmsToIdx:   nmsToIdx,
		gatesByOut: gatesByOut,
		outValues:  outValues,
		bitLen:     bitLen,
	}
}

func (ad *AdderDevice) RunWithInputs(x, y uint64) uint64 {
	ad.Reset()
	for i := 0; i < ad.bitLen; i++ {
		if x&1 == 1 {
			ad.SetXbit(i)
		}
		x >>= 1
		if y&1 == 1 {
			ad.SetYbit(i)
		}
		y >>= 1
	}
	return ad.CollectOutput()
}

func (ad *AdderDevice) SetBitAndRun(xBit, yBit int) uint64 {
	ad.Reset()
	ad.SetXbit(xBit)
	ad.SetYbit(yBit)
	return ad.CollectOutput()
}

func (ad *AdderDevice) CollectOutput() uint64 {
	result := uint64(0)
	for i := range ad.bitLen {
		zName := ad.zNms[i]
		val := uint64(ad.GetVal(zName))
		result |= val << uint64(i)

	}
	return result
}

func (ad *AdderDevice) Inputs() (uint64, uint64) {
	x, y := uint64(0), uint64(0)
	for i := 0; i < ad.bitLen; i++ {
		if ad.X[i] == 1 {
			x |= 1 << i
		}
		if ad.Y[i] == 1 {
			y |= 1 << i
		}
	}
	return x, y
}

func (ad *AdderDevice) CalculateValue(name string) uint8 {
	gate, ok := ad.gatesByOut[name]
	if !ok {
		panic(fmt.Sprintf("No out gate '%s'", name))
	}
	left, right := gate.Left, gate.Right
	lv, rv := ad.GetVal(left), ad.GetVal(right)
	res := gate.EvalBit(lv, rv)
	ad.outValues[name] = res
	return res
}

func (ad *AdderDevice) GetVal(name string) uint8 {
	switch name[0] {
	case 'x':
		return ad.X[ad.ToIdx(name)]
	case 'y':
		return ad.Y[ad.ToIdx(name)]
	default:
		if val, ok := ad.outValues[name]; ok {
			return val
		}
		return ad.CalculateValue(name)
	}
}

func (ad *AdderDevice) SetInputsByWires(setWires map[WireName]struct{}) {
	ad.Reset()
	for wire := range setWires {
		idx := ad.ToIdx(wire)
		if wire[0] == 'x' {
			ad.SetXbit(idx)
			continue
		}
		ad.SetYbit(idx)
	}
}

func (ad *AdderDevice) SetXbit(bit int) {
	ad.X[bit] = 1
}

func (ad *AdderDevice) SetYbit(bit int) {
	ad.Y[bit] = 1
}

func (ad *AdderDevice) Reset() {
	ad.X = make([]uint8, ad.bitLen)
	ad.Y = make([]uint8, ad.bitLen)
	// All out values gates except x, y
	ad.outValues = make(map[WireName]uint8, len(ad.gatesByOut)-(2*ad.bitLen))
}

func (ad *AdderDevice) ToIdx(name string) int {
	idx, ok := ad.nmsToIdx[name]
	if !ok {
		panic(fmt.Sprintf("Not indexed name %s", name))
	}
	return idx
}

func SearchBadBits(adder *AdderDevice) []int {
	badBits := make([]int, 0, 4)
	for xBit := 0; xBit < adder.bitLen-1; xBit++ {
		x := uint64(1 << xBit)
		if adder.SetBitAndRun(xBit, 0) != x+1 {
			badBits = append(badBits, xBit)
		}
	}
	return badBits
}

// Full adder
// X Y C  -> S  Carry
// 0 0 0  -> 0  0
// 0 0 1  -> 1  0
// 0 1 0  -> 1  0
// 0 1 1  -> 0  1
// 1 0 0  -> 1  0
// 1 0 1  -> 0  1
// 1 1 0  -> 0  1
// 1 1 1  -> 1  1
func FindSwapForBit(bit int, adder *AdderDevice) gswp {
	zWireName, nextZWireName := adder.zNms[bit], adder.zNms[bit+1]

	swaps := allSwaps
	for i := 0; i < 8; i++ {
		xBitVal, yBitVal, carryBitVal := (i>>2)&1, (i>>1)&1, i&1
		cinVal := carryBitVal << (bit - 1)
		xVal := uint64((xBitVal << bit) + cinVal)
		yVal := uint64((yBitVal << bit) + cinVal)
		adder.RunWithInputs(xVal, yVal)

		zBitOut, carryBitOut := adder.GetVal(zWireName), adder.GetVal(nextZWireName)

		fa := NewFullAdder(uint64(xBitVal), uint64(yBitVal), uint64(carryBitVal))
		swaps = FilterSwaps(fa, swaps, uint64(zBitOut), uint64(carryBitOut))
		if len(swaps) == 0 {
			panic("No swaps found")
		}

	}
	if len(swaps) != 1 {
		panic(fmt.Sprintf("Expected 1 swap, but: %v", swaps))
	}
	return swaps[0]
}

func swapGatesAxB(adder *AdderDevice, bit int, di *DeviceInfo) (string, string) {
	xInGates := di.InMap[wireName("x", bit)]
	aGate, bGate := xInGates[0], xInGates[1]
	if aGate.Type != XOR {
		aGate, bGate = bGate, aGate
	}
	newMap := swapOuts(aGate, bGate, adder.gatesByOut)
	adder.gatesByOut = newMap

	return aGate.Out, bGate.Out
}

func swapGatesBxS(adder *AdderDevice, bit int, di *DeviceInfo) (string, string) {
	// Swapped
	bGate := adder.gatesByOut[wireName("z", bit)]

	xInGates := di.InMap[wireName("x", bit)]
	aGate := xInGates[0]
	if aGate.Type != XOR {
		aGate = xInGates[1]
	}
	aInGates := di.InMap[aGate.Out]

	sumGate := aInGates[0]
	if sumGate.Type != XOR {
		sumGate = aInGates[1]
	}
	if len(aInGates) != 2 || sumGate.Type != XOR {
		panic("Expected other gates")
	}

	newMap := swapOuts(bGate, sumGate, adder.gatesByOut)
	adder.gatesByOut = newMap
	return bGate.Out, sumGate.Out
}

func swapGatesDxS(adder *AdderDevice, bit int, di *DeviceInfo) (string, string) {
	xInGates := di.InMap[wireName("x", bit)]
	aGate := xInGates[0]
	if aGate.Type != XOR {
		aGate = xInGates[1]
	}
	aInGates := di.InMap[aGate.Out]
	if len(aInGates) != 2 {
		panic("Expected 2 gates")
	}
	one, two := aInGates[0], aInGates[1]
	newMap := swapOuts(one, two, adder.gatesByOut)
	adder.gatesByOut = newMap

	return one.Out, two.Out
}

func swapGatesSxC(adder *AdderDevice, bit int, di *DeviceInfo) (string, string) {
	xInGates := di.InMap[wireName("x", bit)]
	aGate := xInGates[0]
	if aGate.Type != XOR {
		aGate = xInGates[1]
	}
	aInGates := di.InMap[aGate.Out]
	if len(aInGates) != 2 {
		panic("Expected 2 gates")
	}
	sumGate := aInGates[0]
	if sumGate.Type != XOR {
		sumGate = aInGates[1]
	}
	// Swapped
	carryGate := adder.gatesByOut[wireName("z", bit)]

	return carryGate.Out, sumGate.Out
}

func FilterSwaps(fa FullAdder, swaps []gswp, sum, carry uint64) []gswp {
	filtered := make([]gswp, 0, len(swaps))
	for _, swap := range swaps {
		s, c := fa.RunSwap(swap)
		if s == sum && c == carry {
			filtered = append(filtered, swap)
		} else {
		}
	}
	return filtered
}

// FA
// x - XOR - A
// y - |

// x - AND - B
// y - |

// A - XOR - Sum
// C - |
// ==
// x - XOR - XOR - Sum
// y - |     |
// C -------|

// A - AND - D
// C - |
// ==
// x - XOR - AND - D
// y - |     |
// C -------|
//
// D - OR - Carry
// B - |
// ==
// x - XOR - AND - OR - Carry
// y - |    |      |
// C -------|     |
// x ---------- AND
// y -----------|
type FullAdder struct {
	A, B, C, D uint64
	Sum, Carry uint64
}

type gswp int

func (g gswp) String() string {
	switch g {
	case AxB:
		return "A x B"
	case BxS:
		return "B x S"
	case DxS:
		return "D x S"
	case SxC:
		return "S x C"
	default:
		return "Unknown"
	}
}

const (
	_ gswp = iota
	AxB
	BxS
	DxS
	SxC
)

var allSwaps = []gswp{AxB, BxS, DxS, SxC}

// Sum = A ^ C
// D = A & C
// Carry = D | B = (A & C) | B
func NewFullAdder(x, y, c uint64) FullAdder {
	a := x ^ y
	b := x & y
	return FullAdder{
		A:     a,
		B:     b,
		C:     c,
		D:     a & c,
		Sum:   a ^ c,
		Carry: (a & c) | b,
	}
}

func (fa FullAdder) Result() (uint64, uint64) {
	return fa.Sum, fa.Carry
}

func (fa FullAdder) RunSwap(swap gswp) (uint64, uint64) {
	switch swap {
	case AxB:
		return fa.SwapAxB()
	case BxS:
		return fa.SwapBxSum()
	case DxS:
		return fa.SwapDxSum()
	case SxC:
		return fa.SwapSxCarry()
	default:
		panic(fmt.Sprintf("Unknown swap: %d", swap))
	}
}

// Sum = B ^ C
// D = B & C
// Carry = D | A = (B & C) | A
func (fa FullAdder) SwapAxB() (uint64, uint64) {
	return fa.B ^ fa.C, (fa.B & fa.C) | fa.A
}

// Sum = B
// Carry = D | Sum
func (fa FullAdder) SwapBxSum() (uint64, uint64) {
	return fa.B, fa.D | fa.Sum
}

// Sum = D
// Carry = Sum | B
func (fa FullAdder) SwapDxSum() (uint64, uint64) {
	return fa.D, fa.Sum | fa.B
}

// Sum = Carry
// Carry = Sum
func (fa FullAdder) SwapSxCarry() (uint64, uint64) {
	return fa.Carry, fa.Sum
}

func wireName(c string, i int) string {
	return fmt.Sprintf("%s%02d", c, i)
}

func createToNameTranslation(c string, n int) []string {
	nameTranslator := make([]string, n)
	for i := 0; i < n; i++ {
		nameTranslator[i] = wireName(c, i)
	}
	return nameTranslator
}

func createFromNameTranslation(chars []string, n int) map[string]int {
	nameTranslator := make(map[string]int)
	for i := 0; i < n; i++ {
		for _, c := range chars {
			nameTranslator[wireName(c, i)] = i
		}
	}
	return nameTranslator
}

func swapOutNames(gOne, gTwo Gate) (Gate, Gate) {
	gOne.Out, gTwo.Out = gTwo.Out, gOne.Out
	return gOne, gTwo
}

func swapOuts(gateOne, gateTwo Gate, outMap map[WireName]Gate) map[WireName]Gate {
	newMap := make(map[WireName]Gate, len(outMap))
	newOne, newTwo := swapOutNames(gateOne, gateTwo)
	for name, gate := range outMap {
		if gate == gateOne {
			newMap[newOne.Out] = newOne
		} else if gate == gateTwo {
			newMap[newTwo.Out] = newTwo
		} else {
			newMap[name] = gate
		}
	}
	return newMap
}
