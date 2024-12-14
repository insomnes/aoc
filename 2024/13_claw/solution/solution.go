package solution

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"
)

const (
	MaxPushes = 100.0
	TokensA   = 3
	TokensB   = 1
)

type Answer struct {
	A, B  int
	Price int
	Valid bool
}

func NewAnswerFromFloats(a, b float64) Answer {
	if a < 0 || b < 0 {
		return Answer{}
	}

	if a > MaxPushes || b > MaxPushes {
		return Answer{}
	}
	if math.Mod(a, 1) != 0 || math.Mod(b, 1) != 0 {
		return Answer{}
	}
	ansA, ansB := int(a), int(b)
	price := ansA*TokensA + ansB*TokensB

	return Answer{A: ansA, B: ansB, Price: price, Valid: true}
}

func NewAnswer(a, b int, maxVal int) Answer {
	if a <= 0 || b <= 0 {
		return Answer{}
	}

	if maxVal > 0 && (a > MaxPushes || b > MaxPushes) {
		return Answer{}
	}

	price := a*TokensA + b*TokensB
	return Answer{A: a, B: b, Price: price, Valid: true}
}

type EquationSystem struct {
	Xa, Ya   float64
	Xb, Yb   float64
	Xpr, Ypr float64

	XaInt, YaInt   int
	XbInt, YbInt   int
	XprInt, YprInt int
}

// Equations system:
// Xa * a + Xb * b = Xpr
// Ya * a + Yb * b = Ypr
// !!!!!!!!!!!!!!!!!!!!!!!
// a = (Xpr - Xb * b) / Xa
// !!!!!!!!!!!!!!!!!!!!!!!
// Ya * ((Xpr - Xb * b) / Xa) + Yb * b = Ypr
// Ya * Xpr - Ya * Xb * b  + Yb * Xa * b = Ypr * Xa
// b * (Yb * Xa - Ya * Xb) = Ypr * Xa - Ya * Xpr
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// b = (Ypr * Xa - Ya * Xpr) / (Yb * Xa - Ya * Xb)
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func (eq EquationSystem) Solve() Answer {
	if eq.Xa == eq.Xb && eq.Ya == eq.Yb {
		panic(fmt.Sprintf("Xa=%f, Xb=%f, Ya=%f, Yb=%f", eq.Xa, eq.Xb, eq.Ya, eq.Yb))
	}

	b := (eq.Ypr*eq.Xa - eq.Ya*eq.Xpr) / (eq.Yb*eq.Xa - eq.Ya*eq.Xb)
	a := (eq.Xpr - eq.Xb*b) / eq.Xa
	if a == 0.0 || b == 0.0 {
		panic(fmt.Sprintf("zero values: a=%f, b=%f", a, b))
	}
	return NewAnswerFromFloats(a, b)
}

// Determinants approach (bonus: without floating point arithmetic)
// Xa * a + Xb * b = Xpr
// Ya * a + Yb * b = Ypr
//
// Matrix A:  Vec v:    Vec r:
// | Xa Xb |  | a |     | Xpr |
// | Ya Yb |  | b |     | Ypr |
//
// det(A) = Xa * Yb - Xb * Ya
// In our case det(A) == 0 -> unsolvable system
// det(A) != 0 -> may be solvable system, depending on if the roots are integer
// a = det(Aa) / det(A)
// b = det(Ab) / det(A)
func (eq EquationSystem) SolveWithDeterminants() Answer {
	detA := eq.DeterminantA()
	if detA == 0 {
		return Answer{}
	}

	detAa, detAb := eq.DeterminantAa(), eq.DeterminantAb()
	a, remA := divmod(detAa, detA)
	b, remB := divmod(detAb, detA)
	if remA != 0 || remB != 0 { // integer check
		return Answer{}
	}

	return NewAnswer(a, b, -1)
}

func (eq EquationSystem) DeterminantA() int {
	return determinant(eq.XaInt, eq.XbInt, eq.YaInt, eq.YbInt)
}

// Aa -- change a column in matrix A with vector r
// | Xpr Xb |
// | Ypr Yb |
func (eq EquationSystem) DeterminantAa() int {
	return determinant(eq.XprInt, eq.XbInt, eq.YprInt, eq.YbInt)
}

// Ab -- change b column in matrix A with vector r
// | Xa Xpr |
// | Ya Ypr |
func (eq EquationSystem) DeterminantAb() int {
	return determinant(eq.XaInt, eq.XprInt, eq.YaInt, eq.YprInt)
}

// | a  b |
// | c  d |
func determinant(a, b, c, d int) int {
	return a*d - b*c
}

func divmod(a, b int) (int, int) {
	return a / b, a % b
}

// Button A: X+<Xa:%d>, Y+<Ya:%d>
// Button B: X+<Xb:%d>, Y+<Yb:%d>
// Prize: X=<Xpr:%d>, Y=<Ypr:%d>
func ParseEquationSystem(raw []string) (EquationSystem, error) {
	if len(raw) != 3 {
		return EquationSystem{}, fmt.Errorf("Expected 3 lines, got %d", len(raw))
	}

	a, err := ParseButton(raw[0])
	if err != nil {
		return EquationSystem{}, err
	}

	b, err := ParseButton(raw[1])
	if err != nil {
		return EquationSystem{}, err
	}

	pr, err := parsePrize(raw[2])
	if err != nil {
		return EquationSystem{}, err
	}

	var xpr, ypr int
	_, err = fmt.Sscanf(raw[2], "Prize: X=%d, Y=%d", &xpr, &ypr)
	if err != nil {
		return EquationSystem{}, err
	}

	if !checkPositive(a) || !checkPositive(b) || !checkPositive(pr) {
		return EquationSystem{}, fmt.Errorf("All values must be non-zero positive numbers")
	}

	return EquationSystem{
		Xa: float64(a[0]), Ya: float64(a[1]),
		Xb: float64(b[0]), Yb: float64(b[1]),
		Xpr: float64(xpr), Ypr: float64(ypr),
		XaInt: a[0], YaInt: a[1],
		XbInt: b[0], YbInt: b[1],
		XprInt: xpr, YprInt: ypr,
	}, nil
}

func checkPositive(nn [2]int) bool {
	return nn[0] > 0 && nn[1] > 0
}

const buttonPattern = `Button [A|B]: X\+(\d+), Y\+(\d+)`

var buttonRe = regexp.MustCompile(buttonPattern)

func ParseButton(line string) ([2]int, error) {
	res, err := parseEquationLine(line, buttonRe)
	if err != nil {
		return res, fmt.Errorf("parsing button: %v", err)
	}
	return res, nil
}

const prizePattern = `Prize: X=(\d+), Y=(\d+)`

var prizeRe = regexp.MustCompile(prizePattern)

func parsePrize(line string) ([2]int, error) {
	res, err := parseEquationLine(line, prizeRe)
	if err != nil {
		return res, fmt.Errorf("parsing prize: %v", err)
	}
	return res, nil
}

func parseEquationLine(line string, re *regexp.Regexp) ([2]int, error) {
	matches := re.FindStringSubmatch(line)

	var res [2]int

	// Check if the required groups are captured
	if len(matches) != 3 {
		return res, fmt.Errorf("matches for <%s>: %v", line, matches)
	}

	// Extract and convert X and Y values
	x, err := strconv.Atoi(matches[1])
	if err != nil {
		return res, fmt.Errorf("converting X to int: %v", err)
	}
	res[0] = x

	y, err := strconv.Atoi(matches[2])
	if err != nil {
		return res, fmt.Errorf("converting Y to int: %v", err)
	}
	res[1] = y

	return res, nil
}

type ParsedInput = []EquationSystem

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track(time.Now(), "ParseInput")

	equations := make([]EquationSystem, 0)
	rawEqSys := make([]string, 0, 3)

	for i, line := range lines {
		if line != "" {
			rawEqSys = append(rawEqSys, line)
			if i != len(lines)-1 {
				continue
			}
		}
		eqSys, err := ParseEquationSystem(rawEqSys)
		if err != nil {
			return ParsedInput{}, err
		}
		equations = append(equations, eqSys)
		rawEqSys = make([]string, 0, 3)
	}

	return equations, nil
}

func PartOne(inp ParsedInput) int {
	defer Track(time.Now(), "PartOne")
	total := 0

	for _, eq := range inp {
		answer := eq.Solve()
		if answer.Valid {
			total += answer.Price
		}
	}

	return total
}

const clawPosExtra = 10000000000000

func PartTwo(inp ParsedInput) int {
	defer Track(time.Now(), "PartTwo")
	total := 0

	for _, eq := range inp {
		eq.XprInt, eq.YprInt = eq.XprInt+clawPosExtra, eq.YprInt+clawPosExtra
		answer := eq.SolveWithDeterminants()
		if answer.Valid {
			total += answer.Price
		}
	}

	return total
}
