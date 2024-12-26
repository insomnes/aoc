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
	WireVal    = int
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
		panic(fmt.Sprintf("Unknown gate type: %d", g.Type))
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
	AllGates  []Gate
	InitWires WireValues
	OutMap    map[WireName]Gate
	InMap     map[WireName][]Gate
	OutWires  []WireName
	AllWires  map[WireName]struct{}
}

func (di DeviceInfo) CreateDevice() Device {
	return NewDevice(di.OutWires, di.InitWires, di.OutMap)
}

type ParsedInput = DeviceInfo

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track(time.Now(), "ParseInput")
	initWires := make(WireValues, 0)
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
		val := int(valStr[0] - '0')
		initWires[wireName] = val
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
		AllGates:  gates,
		InitWires: initWires,
		OutMap:    outMap,
		InMap:     inMap,
		OutWires:  outWires,
		AllWires:  allWires,
	}, nil
}

type Device struct {
	OutWires []WireName
	Values   WireValues
	OutMap   map[WireName]Gate
	run      bool
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

func (d *Device) CollectOutput() int {
	if !d.run {
		d.Run()
	}

	result := 0
	for i := range d.OutWires {
		wn := fmt.Sprintf("z%02d", i)
		val, set := d.Values[wn]
		if !set {
			panic(fmt.Sprintf("Wire %s not set", wn))
		}
		result |= val << uint(i)

	}
	return result
}

func (d *Device) CollectInput(name string) int {
	result := 0
	for i := 0; i < len(d.OutWires)-1; i++ {
		wn := fmt.Sprintf("%s%02d", name, i)
		val, set := d.Values[wn]
		if !set {
			panic(fmt.Sprintf("Wire %s not set", wn))
		}
		result |= val << uint(i)
	}

	return result
}

func (d *Device) Trace(outBit int) {
	wn := fmt.Sprintf("z%02d", outBit)
	fmt.Println("Tracing", wn)
	gates := 0
	inWires := make([]WireName, 0)
	toVisit := []WireName{wn}
	visited := make(map[WireName]struct{})
	for len(toVisit) > 0 {
		wn = toVisit[0]
		toVisit = toVisit[1:]
		if _, been := visited[wn]; been {
			continue
		}
		visited[wn] = struct{}{}
		gates++

		gate, ok := d.OutMap[wn]
		if !ok {
			inWires = append(inWires, wn)
			continue
		}
		fmt.Println("  Gate", gate)
		toVisit = append(toVisit, gate.Left, gate.Right)
	}

	fmt.Println("In wires", len(inWires), inWires)
	if len(inWires) != 2*(outBit+1) {
		fmt.Println("WRONG")
	}
	fmt.Println("Gates", gates)
}

func NewDevice(outWires []WireName, initWires WireValues, outMap map[WireName]Gate) Device {
	values := make(WireValues)
	for n, v := range initWires {
		values[n] = v
	}

	return Device{OutWires: outWires, Values: values, OutMap: outMap}
}

func PartOne(inp ParsedInput) int {
	defer Track(time.Now(), "PartOne")

	di := inp
	device := di.CreateDevice()
	out := device.CollectOutput()

	return out
}

type Debugger struct {
	OutBits  int
	Ins      map[WireName][]Gate
	Outs     map[WireName]Gate
	AllWires map[WireName]struct{}
}

func (d *Debugger) GetIns(bit int) (Gate, Gate) {
	x, y := fmt.Sprintf("x%02d", bit), fmt.Sprintf("y%02d", bit)
	xIns, yIns := d.Ins[x], d.Ins[y]
	if len(xIns) != 2 || len(yIns) != 2 {
		panic("Wrong number inputs")
	}
	if xIns[0].Type == AND {
		return xIns[0], xIns[1]
	}
	return xIns[1], xIns[0]
}

func (d *Debugger) FindInputGateByType(g Gate, t GateType) Gate {
	l, r := g.Left, g.Right
	lOut := d.Outs[l]
	if lOut.Type == t {
		return lOut
	}
	rOut := d.Outs[r]
	if rOut.Type == t {
		return rOut
	}
	return UnknownGate
}

func (d *Debugger) AssertInvariants(bitLen int) {
	d.assertXYZ(bitLen)
	d.assertAllGatesAreConnected()
}

func (d *Debugger) assertAllGatesAreConnected() {
	for wire := range d.AllWires {
		outGate, wireIsOut := d.Outs[wire]
		if wire[0] == 'x' || wire[0] == 'y' {
			if wireIsOut {
				panic(fmt.Sprintf("Wire %s is OUT for %s", wire, outGate))
			}
		} else {
			if !wireIsOut {
				panic(fmt.Sprintf("Wire %s is NOT OUT", wire))
			}
		}

		ins, wireIsIn := d.Ins[wire]
		if wire[0] == 'z' {
			if wireIsIn {
				panic(fmt.Sprintf("Wire %s is IN for %v", wire, ins))
			}
		} else {
			if !wireIsIn {
				panic(fmt.Sprintf("Wire %s is NOT IN for anything", wire))
			}
		}
	}
}

func (d *Debugger) assertXYZ(bitLen int) {
	lastBit := bitLen - 1
	for bit := 0; bit <= lastBit; bit++ {
		x, y, z := xyz(bit)

		// All outs have gates
		zOut, ok := d.Outs[z]
		if !ok {
			panic(fmt.Sprintf("No output %s: %v", z, zOut))
		}

		// All input lead to two gates except the last one
		xIns, yIns := d.Ins[x], d.Ins[y]
		if bit != lastBit && len(xIns) != 2 {
			panic(fmt.Sprintf("Wrong gates for %s: %v", x, xIns))
		}
		if bit != lastBit && len(yIns) != 2 {
			panic(fmt.Sprintf("Wrong gates for %s: %v", y, yIns))
		}
		if bit != lastBit && !slices.Equal(xIns, yIns) {
			panic(fmt.Sprintf("Different gates for %s: %v, %s: %v", x, xIns, y, yIns))
		}

		if bit != lastBit {
			xAND, xXOR := d.GetIns(bit)
			if xAND.Type != AND {
				panic(fmt.Sprintf("Wrong gate for %s: %v", x, xAND))
			}
			if xXOR.Type != XOR {
				panic(fmt.Sprintf("Wrong gate for %s: %v", x, xXOR))
			}
			yAND, yXOR := d.GetIns(bit)
			if yAND.Type != AND {
				panic(fmt.Sprintf("Wrong gate for %s: %v", y, yAND))
			}
			if yXOR.Type != XOR {
				panic(fmt.Sprintf("Wrong gate for %s: %v", y, yXOR))
			}
		}

		if bit == lastBit && (len(xIns) != 0 || len(yIns) != 0) {
			panic(fmt.Sprintf("Last bit leads to inputs: %s: %v, %s: %v", x, xIns, y, yIns))
		}

	}
}

func xyz(bit int) (string, string, string) {
	return identif("x", bit), identif("y", bit), identif("z", bit)
}

func identif(name string, bit int) string {
	return fmt.Sprintf("%s%02d", name, bit)
}

func PartTwo(inp ParsedInput) int {
	defer Track(time.Now(), "PartTwo")

	di := inp
	device := di.CreateDevice()

	bits := len(device.OutWires)
	dbg := Debugger{
		OutBits:  bits,
		Ins:      di.InMap,
		Outs:     di.OutMap,
		AllWires: di.AllWires,
	}

	// out := device.CollectOutput()
	// x, y := device.CollectInput("x"), device.CollectInput("y")
	// fmt.Println("Device", x, "+", y, "=", out)
	// fmt.Println("Real  ", x, "+", y, "=", x+y)
	// fmt.Printf("X    %b\n", x)
	// fmt.Printf("Y    %b\n", y)
	// fmt.Printf("D   %b\n", out)
	// fmt.Printf("R   %b\n", x+y)
	// fmt.Printf("^   %046b\n", out^(x+y))

	dbg.AssertInvariants(bits)

	return 0
}

func getPairs(slice []string) [][2]string {
	var pairs [][2]string
	n := len(slice)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			pairs = append(pairs, [2]string{slice[i], slice[j]})
		}
	}
	return pairs
}

func getPairCombinations(pairs [][2]string, k int) [][][2]string {
	var result [][][2]string
	comb := make([][2]string, k)
	var generate func(start, depth int)
	generate = func(start, depth int) {
		if depth == k {
			combCopy := make([][2]string, k)
			copy(combCopy, comb)
			result = append(result, combCopy)
			return
		}
		for i := start; i < len(pairs); i++ {
			comb[depth] = pairs[i]
			generate(i+1, depth+1)
		}
	}
	generate(0, 0)
	return result
}
