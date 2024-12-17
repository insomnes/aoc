package solution

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"time"
)

type Program = []uint64

type DebugInfo struct {
	RegA uint64
	RegB uint64
	RegC uint64

	Program Program
}

type ParsedInput = DebugInfo

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track(time.Now(), "ParseInput")

	if len(lines) != 5 {
		return ParsedInput{}, fmt.Errorf("has %d lines not 5", len(lines))
	}

	if lines[3] != "" {
		return ParsedInput{}, fmt.Errorf("line 4 is not empty")
	}

	regA, err := parseRegister(lines[0], "A")
	if err != nil {
		return ParsedInput{}, err
	}
	regB, err := parseRegister(lines[1], "B")
	if err != nil {
		return ParsedInput{}, err
	}
	regC, err := parseRegister(lines[2], "C")
	if err != nil {
		return ParsedInput{}, err
	}

	program, err := parseProgram(lines[4])
	if err != nil {
		return ParsedInput{}, err
	}

	return DebugInfo{RegA: regA, RegB: regB, RegC: regC, Program: program}, nil
}

func parseRegister(line string, regName string) (uint64, error) {
	templateBase := fmt.Sprintf("Register %s:", regName)
	template := templateBase + " %d"

	var value uint64
	_, err := fmt.Sscanf(line, template, &value)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func parseProgram(line string) (Program, error) {
	split := strings.Split(line, ": ")
	if len(split) != 2 {
		return nil, fmt.Errorf("invalid program line: %s", line)
	}
	rawS := []rune(split[1])
	program := make(Program, 0, len(rawS)/2)
	comas := 0

	for _, char := range rawS {
		if char == ',' {
			comas++
			continue
		}
		val := uint64(char - '0')
		if val < 0 || val > 7 {
			return nil, fmt.Errorf("invalid program line: %s", line)
		}
		program = append(program, val)
	}

	if comas != len(program)-1 {
		return nil, fmt.Errorf("invalid program line: %s", line)
	}

	return program, nil
}

type ComboOperator uint64

func (c ComboOperator) ToNum() uint64 {
	return uint64(c)
}

type Operator uint64

func (o Operator) ToNum() uint64 {
	return uint64(o)
}

func NewOperator(n uint64) Operator {
	if n > 7 {
		panic(fmt.Sprintf("invalid op: %d", n))
	}
	return Operator(n)
}

type Instruction uint64

const (
	adv Instruction = iota
	bxl
	bst
	jnz
	bxc
	out
	bdv
	cdv
)

func SelectInstruction(opcode uint64) Instruction {
	switch opcode {
	case 0:
		return adv
	case 1:
		return bxl
	case 2:
		return bst
	case 3:
		return jnz
	case 4:
		return bxc
	case 5:
		return out
	case 6:
		return bdv
	case 7:
		return cdv
	default:
		panic(fmt.Sprintf("invalid opcode: %d", opcode))
	}
}

func (i Instruction) String() string {
	switch i {
	case adv:
		return "adv"
	case bxl:
		return "bxl"
	case bst:
		return "bst"
	case jnz:
		return "jnz"
	case bxc:
		return "bxc"
	case out:
		return "out"
	case bdv:
		return "bdv"
	case cdv:
		return "cdv"
	default:
		panic(fmt.Sprintf("invalid instruction: %d", i))
	}
}

type ChronoComputer struct {
	RegA uint64
	RegB uint64
	RegC uint64

	OutNumbers []uint64
}

func (c ChronoComputer) String() string {
	return fmt.Sprintf("CC{A: %d, B: %d, C: %d}", c.RegA, c.RegB, c.RegC)
}

func (c *ChronoComputer) Run(program Program) {
	insPtr := 0
	for insPtr < len(program) {
		opPointer := insPtr + 1
		// Halt case
		if opPointer >= len(program) {
			break
		}
		opcode := program[insPtr]
		instruction := SelectInstruction(opcode)

		opLiteral := program[opPointer]
		operator := NewOperator(opLiteral)

		insPtr = c.RunInstruction(instruction, operator, insPtr)
		if insPtr < 0 {
			panic(fmt.Sprintf("invalid instruction pointer: %d", insPtr))
		}
	}
}

func (c *ChronoComputer) RunTillJump(program Program) int {
	insPtr := 0
	for insPtr < len(program) {
		opPointer := insPtr + 1
		if opPointer >= len(program) {
			return -1
		}
		opcode := program[insPtr]
		instruction := SelectInstruction(opcode)

		opLiteral := program[opPointer]
		operator := NewOperator(opLiteral)

		insPtr = c.RunInstruction(instruction, operator, insPtr)
		if insPtr < 0 {
			panic(fmt.Sprintf("invalid instruction pointer: %d", insPtr))
		}
		if instruction == jnz {
			return insPtr
		}
	}
	panic("no jump found")
}

func (c *ChronoComputer) RunInstruction(instruction Instruction, op Operator, insPtr int) int {
	newInsPtr := insPtr + 2

	switch instruction {
	case adv:
		c.adv(op)
	case bxl:
		c.bxl(op)
	case bst:
		c.bst(op)
	case jnz:
		newInsPtr = c.jumpLogic(op, insPtr)
	case bxc:
		c.bxc(op)
	case out:
		c.out(op)
	case bdv:
		c.bdv(op)
	case cdv:
		c.cdv(op)
	default:
		panic(fmt.Sprintf("unknown instruction: %d", instruction))
	}

	return newInsPtr
}

func (c *ChronoComputer) jumpLogic(op Operator, curInsPtr int) int {
	newInsPtr := c.jnz(op)
	if newInsPtr == -1 {
		return curInsPtr + 2
	}
	return newInsPtr
}

// A = A / 2 ** combo(op)
func (c *ChronoComputer) adv(op Operator) {
	c.RegA = c.divideA(op)
}

// B = A / 2 ** combo(op)
func (c *ChronoComputer) bdv(op Operator) {
	c.RegB = c.divideA(op)
}

// C = A / 2 ** combo(op)
func (c *ChronoComputer) cdv(op Operator) {
	c.RegC = c.divideA(op)
}

// B = B XOR op
func (c *ChronoComputer) bxl(op Operator) {
	c.RegB = c.RegB ^ op.ToNum()
}

// B = combo(op) % 8
func (c *ChronoComputer) bst(op Operator) {
	comboOperator := c.MakeCmbOp(op)
	res := comboOperator.ToNum() % 8
	c.RegB = res
}

// If A == 0 -> nothing (value is -1)
// If A != 0 -> jump to op (we just return the value)
func (c *ChronoComputer) jnz(op Operator) int {
	if c.RegA == 0 {
		return -1
	}
	if op.ToNum() >= math.MaxInt32 {
		panic(fmt.Sprintf("invalid jump: %d", op.ToNum()))
	}
	newInsPtr := int(op.ToNum())
	return newInsPtr
}

// B = B XOR C, op is ignored
func (c *ChronoComputer) bxc(op Operator) {
	c.RegB = c.RegB ^ c.RegC
}

// prints combo(op) % 8
func (c *ChronoComputer) out(op Operator) {
	comboOperator := c.MakeCmbOp(op)
	res := comboOperator.ToNum() % 8
	c.OutNumbers = append(c.OutNumbers, res)
}

// divideA helper
func (c *ChronoComputer) divideA(op Operator) uint64 {
	numerator := c.RegA
	comboOperator := c.MakeCmbOp(op)
	var denom uint64 = 1 << comboOperator.ToNum()
	return numerator / denom
}

func (c *ChronoComputer) StringNumbers() string {
	var sb strings.Builder
	for i, v := range c.OutNumbers {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	return sb.String()
}

const (
	NoRegister = 0
	RegisterA  = 1
	RegisterB  = 2
	RegisterC  = 3
)

var RegisterNames = []string{"NoRegister", "A", "B", "C"}

func (c *ChronoComputer) MakeCmbOp(op Operator) ComboOperator {
	switch op {
	case 0, 1, 2, 3:
		return ComboOperator(op)
	case 4:
		return ComboOperator(c.RegA)
	case 5:
		return ComboOperator(c.RegB)
	case 6:
		return ComboOperator(c.RegC)
	default:
		panic(fmt.Sprintf("invalid operand: %d", op))
	}
}

func NewChronoComputer(info DebugInfo) ChronoComputer {
	return ChronoComputer{
		RegA: info.RegA,
		RegB: info.RegB,
		RegC: info.RegC,
	}
}

func PartOne(inp ParsedInput) int {
	defer Track(time.Now(), "PartOne")
	total := 0

	computer := NewChronoComputer(inp)
	computer.Run(inp.Program)

	fmt.Println(computer.StringNumbers())

	return total
}

// We assume or know that:
//  0. A is the input register, for output we use B by default
//  1. To print 0 we need to have 8 in the register A, this is our start
//  2. Printing depends on 3 bit intervals (mod 8)
//  3. Each normal order iteration A value is divided by 8, so we
//  4. Multiplication of found value by 8 will give us the next starting value and save
//     our bit interval for printing
//  3. There is also no reason to check values from val to val * 8 after we found the
//     correct val. Cause the number is chunked in 3 bit intervals (see 2)
func ReverseBruteForce(program Program, getReg func(c *ChronoComputer) uint64) uint64 {
	var start uint64 = 1

	// Skip last '0' operator from the program
	for i := len(program) - 2; i >= 0; i-- {
		start *= 8
		num := program[i]
		found := false

		for a := start; a < start*8; a++ {
			c := ChronoComputer{RegA: a, RegB: 0, RegC: 0}
			// We should run only one iteration of our programn and check if the output
			// in register is what we expect, it is fast check and can lead to
			// false positives, which we check later
			c.RunTillJump(program)
			targetVal := getReg(&c)
			if targetVal%8 != num {
				continue
			}
			// If we found interesting value we should run the program to the end
			// and check if the full output is what we expect, program is small, so we
			// can run full calculation here
			c = ChronoComputer{RegA: a, RegB: 0, RegC: 0}
			c.Run(program)
			if !slices.Equal(c.OutNumbers, program[i:]) {
				continue
			}
			start = a
			found = true
			break

		}

		if !found {
			panic(fmt.Sprintf("not found: %v (%d) %d->%d", program[i:], num, start, start*8))
		}
	}

	return start
}

func getRegA(c *ChronoComputer) uint64 {
	return c.RegA
}

func getRegB(c *ChronoComputer) uint64 {
	return c.RegB
}

func PartTwo(inp ParsedInput) uint64 {
	defer Track(time.Now(), "PartTwo")

	return ReverseBruteForce(inp.Program, getRegB)
}
