package solution

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

type NumButton rune

type DoorCode []NumButton

func (dc DoorCode) ToNumber() int {
	num := 0
	for _, b := range dc {
		if b == 'A' {
			continue
		}
		num = num*10 + int(b-'0')
	}
	return num
}

type ParsedInput = []DoorCode

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track(time.Now(), "ParseInput")

	codes := make([]DoorCode, len(lines))
	for i, line := range lines {
		code := make(DoorCode, len(line))
		for j, r := range line {
			code[j] = NumButton(r)
		}
		codes[i] = code
	}
	return codes, nil
}

// +---+---+---+
// | 7 | 8 | 9 |
// +---+---+---+
// | 4 | 5 | 6 |
// +---+---+---+
// | 1 | 2 | 3 |
// +---+---+---+
//     | 0 | A |
//     +---+---+

// To not format
func (n NumButton) String() string {
	return string(n)
}

func (n NumButton) Neighbors() []NumButton {
	// 7 8 9
	// 4 5 6
	// 1 2 3
	//   0 A
	switch n {
	case '0':
		return []NumButton{'2', 'A'}
	case '1':
		return []NumButton{'2', '4'}
	case '2':
		return []NumButton{'1', '3', '5', '0'}
	case '3':
		return []NumButton{'2', '6', 'A'}
	case '4':
		return []NumButton{'1', '5', '7'}
	case '5':
		return []NumButton{'2', '4', '6', '8'}
	case '6':
		return []NumButton{'3', '5', '9'}
	case '7':
		return []NumButton{'4', '8'}
	case '8':
		return []NumButton{'5', '7', '9'}
	case '9':
		return []NumButton{'6', '8'}
	case 'A':
		return []NumButton{'0', '3'}
	default:
		panic("Invalid button")
	}
}

func NumPairToKeyButton(from, to NumButton) KeyButton {
	// 7 8 9
	// 4 5 6
	// 1 2 3
	//   0 A
	switch from {
	case '0':
		switch to {
		case '2':
			return '^'
		case 'A':
			return '>'
		default:
			panic("Invalid button")
		}
	case '1':
		switch to {
		case '2':
			return '>'
		case '4':
			return '^'
		default:
			panic("Invalid button")
		}
	case '2':
		switch to {
		case '1':
			return '<'
		case '3':
			return '>'
		case '5':
			return '^'
		case '0':
			return 'v'
		default:
			panic("Invalid button")
		}
	case '3':
		switch to {
		case '2':
			return '<'
		case '6':
			return '^'
		case 'A':
			return 'v'
		default:
			panic("Invalid button")
		}
	case '4':
		switch to {
		case '1':
			return 'v'
		case '5':
			return '>'
		case '7':
			return '^'
		default:
			panic("Invalid button")
		}
	case '5':
		switch to {
		case '2':
			return 'v'
		case '4':
			return '<'
		case '6':
			return '>'
		case '8':
			return '^'
		default:
			panic("Invalid button")
		}
	case '6':
		switch to {
		case '3':
			return 'v'
		case '5':
			return '<'
		case '9':
			return '^'
		default:
			panic("Invalid button")
		}
	case '7':
		switch to {
		case '4':
			return 'v'
		case '8':
			return '>'
		default:
			panic("Invalid button")
		}
	case '8':
		switch to {
		case '5':
			return 'v'
		case '7':
			return '<'
		case '9':
			return '>'
		default:
			panic("Invalid button")
		}
	case '9':
		switch to {
		case '6':
			return 'v'
		case '8':
			return '<'
		default:
			panic("Invalid button")
		}
	case 'A':
		switch to {
		case '0':
			return '<'
		case '3':
			return '^'
		default:
			panic("Invalid button")
		}
	default:
		panic("Invalid button")
	}
}

var AllNumButtons = []NumButton{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'A'}

type (
	NumButtonRoute = []NumButton
)

type Numpad struct {
	routes map[NumButton]map[NumButton][]NumButtonRoute
}

func NewNumpad() *Numpad {
	n := &Numpad{}
	n.initRoutes()
	return n
}

func (n *Numpad) initRoutes() {
	if n.routes != nil {
		return
	}
	n.routes = make(map[NumButton]map[NumButton][]NumButtonRoute)
	for _, start := range AllNumButtons {
		routes := make(map[NumButton][]NumButtonRoute)
		for _, end := range AllNumButtons {
			if start == end {
				continue
			}
			routes[end] = findAllShortRoutes(start, end)
		}
		n.routes[start] = routes
	}
}

func (n *Numpad) GetRoutes(start, end NumButton) []NumButtonRoute {
	routes, ok := n.routes[start][end]
	if !ok {
		panic(fmt.Sprintf("Invalid route %s -> %s", start, end))
	}
	return routes
}

func (n *Numpad) CodeToNumPositionSequences(code DoorCode) []NumPositionSequence {
	sequences := make([]NumPositionSequence, 0)
	prev := code[0]
	initialRoutes := n.GetRoutes('A', prev)
	for _, route := range initialRoutes {
		seq := NPSFromRoute('A', route)
		sequences = append(sequences, seq)
	}
	for _, button := range code[1:] {
		curRoutes := n.GetRoutes(prev, button)
		newSequences := make([]NumPositionSequence, 0, len(sequences)*len(curRoutes))
		for _, seq := range sequences {
			for _, route := range curRoutes {
				newSeq := seq.Clone()
				newSeq.AddRoute(prev, route)
				newSequences = append(newSequences, newSeq)
			}
		}
		sequences = newSequences
		prev = button
	}
	return sequences
}

// inefficient but who cares
func findAllShortRoutes(start, end NumButton) []NumButtonRoute {
	// from end to start
	visited := make(map[NumButton]bool)
	toVisit := []NumButton{end}
	routeMap := make(map[NumButton][]NumButtonRoute)
	routeMap[end] = []NumButtonRoute{[]NumButton{}}
	for len(toVisit) > 0 {
		current := toVisit[0]
		toVisit = toVisit[1:]
		if current == start {
			continue
		}
		if visited[current] {
			continue
		}
		visited[current] = true

		curRoutes := routeMap[current]
		for _, neighbor := range current.Neighbors() {
			if visited[neighbor] {
				continue
			}
			toVisit = append(toVisit, neighbor)

			neighborRoutes := routeMap[neighbor]

			if len(neighborRoutes) == 0 {
				for _, route := range curRoutes {
					neighborRoutes = append(neighborRoutes, append(route, current))
				}
				routeMap[neighbor] = neighborRoutes
				continue
			}
			nRoute := neighborRoutes[0]
			if len(nRoute) < len(curRoutes[0])+1 {
				continue
			}
			if len(nRoute) == len(curRoutes[0])+1 {
				for _, route := range curRoutes {
					neighborRoutes = append(neighborRoutes, append(route, current))
				}
				routeMap[neighbor] = neighborRoutes
				continue
			}
			neighborRoutes = make([]NumButtonRoute, 0, len(curRoutes))
			for i, route := range curRoutes {
				neighborRoutes[i] = append(route, current)
			}
			routeMap[neighbor] = neighborRoutes

		}
	}
	routes := routeMap[start]
	for i, route := range routes {
		slices.Reverse(route)
		routes[i] = route
	}

	return routes
}

type KeyButton rune

func (k KeyButton) String() string {
	return string(k)
}

//	 ^ A
// < v >

// It should be
func KeyPairToKeyButton(from, to KeyButton) KeyButton {
	switch from {
	case '^':
		switch to {
		case 'A':
			return '>'
		case 'v':
			return 'v'
		}
	case 'A':
		switch to {
		case '^':
			return '<'
		case '>':
			return 'v'
		}
	case '>':
		switch to {
		case 'A':
			return '^'
		case 'v':
			return '<'
		}
	case 'v':
		switch to {
		case '^':
			return '^'
		case '>':
			return '>'
		case '<':
			return '<'
		}
	case '<':
		if to == 'v' {
			return '>'
		}
	}
	panic("Invalid key pair")
}

//     +---+---+
//     | ^ | A |
// +---+---+---+
// | < | v | > |
// +---+---+---+

// It should be
var kbRoutes = map[KeyButton]map[KeyButton][][]KeyButton{
	'A': {
		'^': {{'^'}},
		'>': {{'>'}},
		'v': {{'^', 'v'}, {'>', 'v'}},
		'<': {{'^', 'v', '<'}, {'>', 'v', '<'}},
		'A': {{}},
	},
	'v': {
		'^': {{'^'}},
		'>': {{'>'}},
		'<': {{'<'}},
		'A': {{'^', 'A'}, {'>', 'A'}},
		'v': {{}},
	},
	'^': {
		'A': {{'A'}},
		'v': {{'v'}},
		'<': {{'v', '<'}},
		'>': {{'v', '>'}, {'A', '>'}},
		'^': {{}},
	},
	'>': {
		'A': {{'A'}},
		'v': {{'v'}},
		'<': {{'v', '<'}},
		'^': {{'v', '^'}, {'A', '^'}},
		'>': {{}},
	},
	'<': {
		'v': {{'v'}},
		'>': {{'v', '>'}},
		'^': {{'v', '^'}},
		'A': {{'v', '^', 'A'}, {'v', '>', 'A'}},
		'<': {{}},
	},
}

type NumPositionMove struct {
	From  NumButton
	To    NumButton
	Route NumButtonRoute
}

func (npm NumPositionMove) String() string {
	return fmt.Sprintf("%s-%s %s", npm.From, npm.To, npm.Route)
}

func NewNumPositionMove(from NumButton, route NumButtonRoute) NumPositionMove {
	to := route[len(route)-1]
	return NumPositionMove{
		From:  from,
		To:    to,
		Route: route,
	}
}

func (npm *NumPositionMove) Clone() NumPositionMove {
	return NumPositionMove{
		From:  npm.From,
		To:    npm.To,
		Route: slices.Clone(npm.Route),
	}
}

type NumPositionSequence struct {
	Buttons []NumButton
	Moves   []NumPositionMove
}

func (nps NumPositionSequence) String() string {
	var sb strings.Builder
	sb.WriteRune('<')
	for i, b := range nps.Buttons {
		sb.WriteRune(rune(b))
		if i < len(nps.Buttons)-1 {
			sb.WriteRune(' ')
		}
	}
	sb.WriteRune('>')
	for _, m := range nps.Moves {
		sb.WriteString(fmt.Sprintf(" -> %s", m))
	}
	return sb.String()
}

func NPSFromRoute(from NumButton, route NumButtonRoute) NumPositionSequence {
	if len(route) == 0 {
		panic(fmt.Sprintf("Invalid route %s -> %s", from, route))
	}
	move := NewNumPositionMove(from, route)
	return NPSFromMove(move)
}

func NPSFromMove(move NumPositionMove) NumPositionSequence {
	return NumPositionSequence{
		Buttons: []NumButton{move.From, move.To},
		Moves:   []NumPositionMove{move},
	}
}

func (nps *NumPositionSequence) AddRoute(from NumButton, route NumButtonRoute) {
	move := NewNumPositionMove(from, route)
	nps.AddMove(move)
}

func (nps *NumPositionSequence) AddMove(move NumPositionMove) {
	last := nps.Buttons[len(nps.Buttons)-1]
	if last != move.From {
		panic("Invalid move")
	}
	nps.Buttons = append(nps.Buttons, move.To)
	nps.Moves = append(nps.Moves, move)
}

func (nps *NumPositionSequence) Clone() NumPositionSequence {
	clone := NumPositionSequence{
		Buttons: make([]NumButton, len(nps.Buttons)),
		Moves:   make([]NumPositionMove, len(nps.Moves)),
	}
	copy(clone.Buttons, nps.Buttons)
	for i, move := range nps.Moves {
		clone.Moves[i] = move.Clone()
	}
	return clone
}

type KeyPadPressSequence []KeyButton

func NewKeyPadPressSequence(nps NumPositionSequence) KeyPadPressSequence {
	seq := make(KeyPadPressSequence, 0)
	for _, move := range nps.Moves {
		prev := move.From
		for _, numButton := range move.Route {
			seq = append(seq, NumPairToKeyButton(prev, numButton))
			prev = numButton
		}
		seq = append(seq, 'A')
	}
	return seq
}

func NewKeyPadPressSequenceFromKeyPos(kpps KeyPadPositionSequence) KeyPadPressSequence {
	seq := make(KeyPadPressSequence, 0)
	for _, move := range kpps.Moves {
		prev := move.From
		for _, keyButton := range move.Route {
			seq = append(seq, KeyPairToKeyButton(prev, keyButton))
			prev = keyButton
		}
		seq = append(seq, 'A')
	}
	return seq
}

type KeyButtonRoute = []KeyButton

type KeyPad struct {
	routes map[KeyButton]map[KeyButton][]KeyButtonRoute
}

func (kp *KeyPad) GetRoutes(start, end KeyButton) []KeyButtonRoute {
	routes, ok := kp.routes[start][end]
	if !ok {
		panic(fmt.Sprintf("Invalid route %s -> %s", start, end))
	}
	return routes
}

func (kp *KeyPad) PressSeqToPositionSeqs(pressSeq KeyPadPressSequence) []KeyPadPositionSequence {
	sequences := make([]KeyPadPositionSequence, 0)
	prev := pressSeq[0]
	initialRoutes := kp.GetRoutes('A', prev)
	for _, route := range initialRoutes {
		seq := KPPSFromRoute('A', route)
		sequences = append(sequences, seq)
	}
	for _, key := range pressSeq[1:] {
		curRoutes := kp.GetRoutes(prev, key)
		newSequences := make([]KeyPadPositionSequence, 0, len(sequences)*len(curRoutes))
		for _, seq := range sequences {
			for _, route := range curRoutes {
				newSeq := seq.Clone()
				newSeq.AddRoute(prev, route)
				newSequences = append(newSequences, newSeq)
			}
		}
		sequences = newSequences
		prev = key
	}
	return sequences
}

func NewKeyPad() *KeyPad {
	kp := &KeyPad{}
	m := make(map[KeyButton]map[KeyButton][]KeyButtonRoute)
	for k, v := range kbRoutes {
		m[k] = v
	}
	kp.routes = m
	return kp
}

func createKeyPadPrssSeqsFromNumPos(numPosSeqs []NumPositionSequence) []KeyPadPressSequence {
	seqs := make([]KeyPadPressSequence, len(numPosSeqs))
	for i, nps := range numPosSeqs {
		seqs[i] = NewKeyPadPressSequence(nps)
	}
	return seqs
}

func createKeyPadPrssSeqsFromKeyPos(keyPosSeqs []KeyPadPositionSequence) []KeyPadPressSequence {
	seqs := make([]KeyPadPressSequence, len(keyPosSeqs))
	for i, kps := range keyPosSeqs {
		seqs[i] = NewKeyPadPressSequenceFromKeyPos(kps)
	}
	return seqs
}

type KeyPadPositionMove struct {
	From  KeyButton
	To    KeyButton
	Route KeyButtonRoute
}

func (kppm KeyPadPositionMove) String() string {
	return fmt.Sprintf("%s - %s %s", kppm.From, kppm.To, kppm.Route)
}

func NewKeyPadPositionMove(from KeyButton, route KeyButtonRoute) KeyPadPositionMove {
	var to KeyButton
	if len(route) == 0 {
		to = from
	} else {
		to = route[len(route)-1]
	}
	return KeyPadPositionMove{
		From:  from,
		To:    to,
		Route: route,
	}
}

func (kppm *KeyPadPositionMove) Clone() KeyPadPositionMove {
	return KeyPadPositionMove{
		From:  kppm.From,
		To:    kppm.To,
		Route: slices.Clone(kppm.Route),
	}
}

type KeyPadPositionSequence struct {
	Buttons []KeyButton
	Moves   []KeyPadPositionMove
}

func KPPSFromRoute(from KeyButton, route KeyButtonRoute) KeyPadPositionSequence {
	move := NewKeyPadPositionMove(from, route)
	return KPPSFromMove(move)
}

func KPPSFromMove(move KeyPadPositionMove) KeyPadPositionSequence {
	return KeyPadPositionSequence{
		Buttons: []KeyButton{move.From, move.To},
		Moves:   []KeyPadPositionMove{move},
	}
}

func (kpps KeyPadPositionSequence) String() string {
	var sb strings.Builder
	sb.WriteRune('(')
	for i, b := range kpps.Buttons {
		sb.WriteRune(rune(b))
		if i < len(kpps.Buttons)-1 {
			sb.WriteRune(' ')
		}
	}
	sb.WriteRune(')')
	for _, m := range kpps.Moves {
		sb.WriteString(fmt.Sprintf(" -> %s", m))
	}
	return sb.String()
}

func (kpps *KeyPadPositionSequence) AddRoute(from KeyButton, route KeyButtonRoute) {
	move := NewKeyPadPositionMove(from, route)
	kpps.AddMove(move)
}

func (kpps *KeyPadPositionSequence) AddMove(move KeyPadPositionMove) {
	last := kpps.Buttons[len(kpps.Buttons)-1]
	if last != move.From {
		panic("Invalid move")
	}
	kpps.Buttons = append(kpps.Buttons, move.To)
	kpps.Moves = append(kpps.Moves, move)
}

func (kpps *KeyPadPositionSequence) Clone() KeyPadPositionSequence {
	clone := KeyPadPositionSequence{
		Buttons: make([]KeyButton, len(kpps.Buttons)),
		Moves:   make([]KeyPadPositionMove, len(kpps.Moves)),
	}
	copy(clone.Buttons, kpps.Buttons)
	for i, move := range kpps.Moves {
		clone.Moves[i] = move.Clone()
	}
	return clone
}

func createKeyPadPositionSeqs(
	keyPadPressSeqs []KeyPadPressSequence,
	kp KeyPad,
) []KeyPadPositionSequence {
	seqs := make([]KeyPadPositionSequence, 0)
	for _, pressSeq := range keyPadPressSeqs {
		posSeqs := kp.PressSeqToPositionSeqs(pressSeq)
		for _, posSeq := range posSeqs {
			seqs = append(seqs, posSeq)
		}
	}
	return seqs
}

func leaveOnlyMinLen(seq []KeyPadPressSequence) []KeyPadPressSequence {
	minLen := 99999999999
	for _, s := range seq {
		if len(s) < minLen {
			minLen = len(s)
		}
	}
	newSeq := make([]KeyPadPressSequence, 0, len(seq))
	for _, s := range seq {
		if len(s) == minLen {
			newSeq = append(newSeq, s)
		}
	}
	return newSeq
}

// <vA<AA>>^AvAA<^A>A<v<A>>^AvA^A<vA>^A<v<A>^A>AAvA^A<v<A>A>^AAAvA<^A>A
// v<<A>>^A<A>AvA<^AA>A<vAAA>^A
// <A^A>^^AvvvA
// 029A
func PartOne(inp ParsedInput) int {
	defer Track(time.Now(), "PartOne")
	total := 0
	for _, code := range inp {
		minLen := codeMinSequenceLen(code)
		fmt.Println(code, code.ToNumber(), minLen, minLen*code.ToNumber())
		total += minLen * code.ToNumber()

	}
	return total
}

func codeMinSequenceLen(code DoorCode) int {
	np := NewNumpad()

	// Input press sequence is code

	// Robot 1 arm movement
	numPosSequences := np.CodeToNumPositionSequences(code)

	// Robot 1 console press sequence
	keyPadOne := NewKeyPad()
	keypadOnePressSeqs := createKeyPadPrssSeqsFromNumPos(numPosSequences)
	keypadOnePressSeqs = leaveOnlyMinLen(keypadOnePressSeqs)

	// Robot 2 arm movement
	keypadOnePositionSeqs := createKeyPadPositionSeqs(keypadOnePressSeqs, *keyPadOne)

	// Robot 2 console press sequence
	keyPadTwo := NewKeyPad()
	keypadTwoPressSeqs := createKeyPadPrssSeqsFromKeyPos(keypadOnePositionSeqs)

	keypadTwoPressSeqs = leaveOnlyMinLen(keypadTwoPressSeqs)

	// Robot 3 arm movement
	keypadTwoPositionSeqs := createKeyPadPositionSeqs(keypadTwoPressSeqs, *keyPadTwo)

	// Robot 3 console press sequence
	_ = NewKeyPad()
	keypadThreePressSeqs := createKeyPadPrssSeqsFromKeyPos(keypadTwoPositionSeqs)
	// keypadThreePressSeqs = leaveOnlyMinLen(keypadThreePressSeqs)

	minLen := 99999999999
	for _, s := range keypadThreePressSeqs {
		if len(s) < minLen {
			minLen = len(s)
		}
	}

	return minLen
}

func PartTwo(inp ParsedInput) uint64 {
	defer Track(time.Now(), "PartTwo")

	var total uint64 = 0

	return total
}
