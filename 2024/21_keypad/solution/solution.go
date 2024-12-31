package solution

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func ToNumberCode(dc []rune) int {
	num := 0
	for _, b := range dc {
		if b == 'A' {
			continue
		}
		num = num*10 + int(b-'0')
	}
	return num
}

type ParsedInput = [][]rune

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track(time.Now(), "ParseInput")

	codes := make([][]rune, len(lines))
	for i, line := range lines {
		code := make([]rune, len(line))
		for j, r := range line {
			code[j] = r
		}
		codes[i] = code
	}
	return codes, nil
}

func codeToNumber(code []rune) uint {
	num := uint(0)
	for _, c := range code[0 : len(code)-1] {
		if c == 'A' {
			panic("Why is there an A in the cut code?")
		}
		num = num*10 + uint(c-'0')
	}
	return num
}

const partOneDepth = 2

func PartOne(inp ParsedInput) uint {
	defer Track(time.Now(), "PartOne")

	total := uint(0)
	codes := inp
	cr := NewCodeRouter(partOneDepth)

	for _, code := range codes {
		cost := cr.CodeCost(code)
		total += cost * codeToNumber(code)
	}

	return total
}

const partTwoDepth = 25

func PartTwo(inp ParsedInput) uint {
	defer Track(time.Now(), "PartTwo")

	var total uint = 0
	cr := NewCodeRouter(partTwoDepth)

	for _, code := range inp {
		cost := cr.CodeCost(code)
		total += cost * codeToNumber(code)
	}

	return total
}

type CodeRouter struct {
	costCalculator *CostCalculator

	depth uint
}

func NewCodeRouter(depth uint) *CodeRouter {
	cc := NewCostCalculator()
	return &CodeRouter{costCalculator: cc, depth: depth}
}

func (cr *CodeRouter) CodeCost(code []rune) uint {
	from := 'A'
	codeCost := uint(0)

	for _, to := range code {
		moveCost := cr.MoveCost(from, to)
		codeCost += moveCost
		from = to
	}

	return codeCost
}

func (cr *CodeRouter) MoveCost(from, to rune) uint {
	mpadRoutes := cr.findAllMovepadRoutes(from, to)
	minCost := uint(math.MaxUint)
	for _, route := range mpadRoutes {
		routeCost := uint(0)
		for _, press := range route {
			routeCost += cr.costCalculator.Cost(press)
			if routeCost >= minCost {
				break
			}
		}
		if routeCost < minCost {
			minCost = routeCost
		}
	}
	return minCost
}

// From 0 to 1, A to 9, 7 to 2 ets
// Returns Move pad presses one level deeper
func (cr *CodeRouter) findAllMovepadRoutes(from, to rune) [][]Press {
	allRoutes := findAllNumpadRoutes(from, to)
	allPresses := make([][]Press, 0, len(allRoutes))
	for _, route := range allRoutes {
		moveSeq := routeToMoveSequence(route)
		routePresses := toPresses(moveSeq, cr.depth)
		allPresses = append(allPresses, routePresses)
	}
	return allPresses
}

func toPresses(pseq string, depth uint) []Press {
	prev := 'A'
	presses := make([]Press, len(pseq))
	for i, r := range pseq {
		presses[i] = Press{Prev: prev, Button: r, Depth: depth}
		prev = r
	}
	return presses
}

type CostCalculator struct {
	cache map[Press]uint
}

func NewCostCalculator() *CostCalculator {
	return &CostCalculator{cache: make(map[Press]uint)}
}

func (cc *CostCalculator) Cost(press Press) uint {
	if press.Depth == 0 {
		return 1
	}
	if cost, ok := cc.cache[press]; ok {
		return cost
	}

	deeperPresses := press.GoDeeper()
	minCost := uint(math.MaxUint)
	for _, dps := range deeperPresses {
		cost := uint(0)
		for _, dp := range dps {
			cost += cc.Cost(dp)
			if cost > minCost {
				break
			}
		}
		if cost < minCost {
			minCost = cost
		}
	}
	cc.cache[press] = minCost

	return minCost
}

type Press struct {
	Prev   rune
	Button rune
	Depth  uint
}

func (p Press) String() string {
	return fmt.Sprintf("%d: %c/%c", p.Depth, p.Prev, p.Button)
}

func (p Press) GoDeeper() [][]Press {
	if p.Depth == 0 {
		panic("Cannot go deeper")
	}
	allRoutes := findMpadRoutes(p.Prev, p.Button)
	allPresses := make([][]Press, 0, len(allRoutes))

	var prev rune
	for _, route := range allRoutes {
		presses := make([]Press, len(route)+1)
		for i, mv := range route {
			if i == 0 {
				prev = 'A'
			}
			presses[i] = Press{Prev: prev, Button: mv, Depth: p.Depth - 1}
			prev = mv
		}

		press := Press{Prev: prev, Button: 'A', Depth: p.Depth - 1}
		if len(route) == 0 {
			press = Press{Prev: 'A', Button: 'A', Depth: p.Depth - 1}
		}
		presses[len(route)] = press
		allPresses = append(allPresses, presses)
	}
	return allPresses
}

func findMpadRoutes(from, to rune) []MovepadRoute {
	switch from {
	case 'A':
		return aRoutes[to]
	case '^':
		return upRoutes[to]
	case 'v':
		return downRoutes[to]
	case '>':
		return rightRoutes[to]
	case '<':
		return leftRoutes[to]
	}
	return nil
}

type MovepadRoute []rune

func (mr MovepadRoute) String() string {
	return string(mr)
}

// # ^ A
// < v >
// skipfmt
var aRoutes = map[rune][]MovepadRoute{
	'A': {{}},
	'^': {{'<'}},
	'v': {{'v', '<'}, {'<', 'v'}},
	'>': {{'v'}},
	'<': {{'v', '<', '<'}, {'<', 'v', '<'}},
}

var upRoutes = map[rune][]MovepadRoute{
	'A': {{'>'}},
	'^': {{}},
	'v': {{'v'}},
	'>': {{'v', '>'}, {'>', 'v'}},
	'<': {{'v', '<'}},
}

var downRoutes = map[rune][]MovepadRoute{
	'A': {{'>', '^'}, {'^', '>'}},
	'^': {{'^'}},
	'v': {{}},
	'>': {{'>'}},
	'<': {{'<'}},
}

var rightRoutes = map[rune][]MovepadRoute{
	'A': {{'^'}},
	'^': {{'^', '<'}, {'<', '^'}},
	'v': {{'<'}},
	'>': {{}},
	'<': {{'<', '<'}},
}

var leftRoutes = map[rune][]MovepadRoute{
	'A': {{'>', '>', '^'}, {'>', '^', '>'}},
	'^': {{'>', '^'}},
	'v': {{'>'}},
	'>': {{'>', '>'}},
	'<': {{}},
}

func isNumpadNeighbor(fromButton, toButton rune) (rune, bool) {
	neighs := getNumpadNeighbors(fromButton)
	for i, n := range neighs {
		if n == toButton {
			return neighborMoves[i], true
		}
	}
	return rune(0), false
}

func getNumpadNeighbors(button rune) [4]rune {
	return numpadNeighbors[numButtToKey(button)]
}

func numButtToKey(b rune) int {
	if b == 'A' {
		return 10
	}
	return int(b - '0')
}

// Missing neighbor
const noNeigh = rune(0)

var neighborMoves = [4]rune{'<', '^', '>', 'v'}

// 7 8 9
// 4 5 6
// 1 2 3
// # 0 A
// Neighbors go as: left (<), up (^), right (>), down (v)
var numpadNeighbors = [11][4]rune{
	{noNeigh, '2', 'A', noNeigh}, // 0
	{noNeigh, '4', '2', noNeigh}, // 1
	{'1', '5', '3', '0'},         // 2
	{'2', '6', noNeigh, 'A'},     // 3
	{noNeigh, '7', '5', '1'},     // 4
	{'4', '8', '6', '2'},         // 5
	{'5', '9', noNeigh, '3'},     // 6
	{noNeigh, noNeigh, '8', '4'}, // 7
	{'7', noNeigh, '9', '5'},     // 8
	{'8', noNeigh, noNeigh, '6'}, // 9
	{'0', '3', noNeigh, noNeigh}, // A
}

func routeToMoveSequence(route []rune) string {
	sb := strings.Builder{}
	prev := route[0]
	for _, r := range route[1:] {
		move, ok := isNumpadNeighbor(prev, r)
		if !ok {
			panic("Invalid route")
		}
		sb.WriteRune(move)
		prev = r
	}
	sb.WriteRune('A')
	return sb.String()
}

func findAllNumpadRoutes(fromButton, toButton rune) [][]rune {
	toVisit := [][]rune{{fromButton}}
	routes := make([][]rune, 0, 8)
	minRouteLen := 999
	for len(toVisit) > 0 {
		curPath := toVisit[len(toVisit)-1]
		curButton := curPath[len(curPath)-1]
		toVisit = toVisit[:len(toVisit)-1]

		if curButton != toButton {
			if len(curPath) >= minRouteLen {
				continue
			}
			for _, n := range getNumpadNeighbors(curButton) {
				if n == noNeigh || pathContains(curPath, n) {
					continue
				}
				toVisit = append(toVisit, addToPath(curPath, n))
			}
			continue
		}

		if len(curPath) == minRouteLen {
			routes = append(routes, curPath)
			continue
		}

		minRouteLen = len(curPath)
		routes = [][]rune{curPath}
		continue

	}
	return routes
}

func addToPath(path []rune, n rune) []rune {
	newPath := make([]rune, len(path)+1)
	copy(newPath, path)
	newPath[len(newPath)-1] = n
	return newPath
}

func pathContains(path []rune, n rune) bool {
	for _, p := range path {
		if p == n {
			return true
		}
	}
	return false
}
