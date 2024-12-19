package main

import (
	"adventofcode"
	"adventofcode/toolbox/assert"
	"adventofcode/toolbox/conversion"
	"adventofcode/toolbox/datatypes"
	"adventofcode/toolbox/fs"
	"adventofcode/toolbox/text"
	"bytes"
	"fmt"
	"math"
	"slices"
	"strings"
)

func main() {
	adventofcode.Time(func() {
		d := fs.LoadFile("2024/day_seventeen/input.txt")
		c := NewComputer(d)

		c.Run()

		fmt.Printf("Part 1: %s\n", c.StdoutString())
	})
	adventofcode.Time(func() {
		d := fs.LoadFile("2024/day_seventeen/input.txt")
		comp := NewComputer(d)

		a := solveReverse(comp)

		fmt.Printf("Part 2: %d\n", a)
	})
}

// create recursive function to iterate through the possible results
func manualFind(nextI, nextA int, program []int) (int, bool) {
	if len(program) == 0 {
		return 0, false
	}
	nextB := program[len(program)-1]
	fmt.Printf("nextI: %d, nextA: %d, nextB: %d\n", nextI, nextA, nextB)

	stopI := (nextA ^ 3) + 8
	for i := nextI; i < stopI; i++ {
		a, b, _ := manual(i)

		if a == nextA && b == nextB {
			//fmt.Printf("\t%d\n", i)
			if rI, ok := manualFind(i, i, program[:len(program)-1]); ok {
				return rI, true
			}
		}
	}
	return 0, false
}

func bruteforce(result int) (a, b, c int) {
	x := 1
	_ = 6*x + 4
	_ = b ^ c //
	_ = 2

	return 0, 0, 0
}

// ZST A=2024	B=46187030	C=0
// BST B A%8
// ZST A=2024	B=0	C=0
// BXL B B⊻5
// ZST A=2024	B=5	C=0
// CDV C A/(2^B)
// ZST A=2024	B=5	C=63
// ADV A A/(2^3)
// ZST A=253	B=5	C=63
// BXC B B⊻C
// ZST A=253	B=58	C=63
// BXL B B⊻6
// ZST A=253	B=60	C=63
//
//	OUT B
func manual(a int) (ra, rb, rc int) {
	b := a & 7
	b = (b ^ 0b101)
	c := int(float64(a) / math.Pow(2.0, float64(b)))
	//fmt.Printf("\t\t%f\n", float64(a)/math.Pow(2.0, float64(b)))
	a = a / 8
	b = b ^ c
	b = b ^ 0b110

	return a, b & 7, c
}

func solveReverse(c *computer) int {
	answer := c.program
	options := []int{}
	nextOptions := []int{0}
	solutions := []int{}

	for len(solutions) == 0 {
		options = nextOptions
		nextOptions = []int{}
		for _, option := range options {
			a := option << 3

			for i := range 8 {
				input := a + i
				current := []int{}

				var b int
				for input > 0 {
					input, b, _ = manual(input)
					current = append(current, b)
				}

				if len(current) != 0 && slices.Equal(current, answer[len(answer)-len(current):]) {
					if slices.Equal(current, answer) {
						solutions = append(solutions, a+i)
					}
					nextOptions = append(nextOptions, a+i) // was missing this from initial attempts, only took the first few results
				}
			}
		}
	}

	slices.Sort(solutions)

	if len(solutions) == 0 {
		return 0
	}

	return solutions[0]
}

func solveBrute2(c *computer) int {
	p := c.program

	var i, step int
	o := make([]int, len(p))
	for i = 1; ; i++ {
		for slices.Equal(p[:step], o[:step]) {
			if step >= 4 {
				fmt.Printf("%v\n", o)
			}
			_, o[step], _ = manual(i)
			step++
		}
		o = make([]int, len(p))
	}

	return 0
}

func solveBrute(c *computer) int {
	program := text.Texts(conversion.To(c.program, func(i int) text.Text {
		return text.Text(fmt.Sprintf("%d", i))
	})).Join(",").String()

	fmt.Printf("Looking for: %s\n", program)
	// 7065690000
	// 16848745000
	// 19984957000
	// 4960738152201
	// 5029457628937
	// 7724811783945
	// 7864406576905
	// 8001845530377 (34359738368) (33554432)
	// 136902133485321 (matches)
	var lastMatch int
	for a := 1; a <= 12303; a++ {
		if a%1000 == 0 {
			fmt.Printf("A: %d\r", a)
		}
		snapshot := c.Image()
		snapshot.debug = false
		snapshot.registers.SetX("A", a)
		snapshot.Run(func(c *computer) bool {
			output := c.StdoutString()
			if output == program {
				return true
			} else if strings.HasPrefix(program, output) {
				// track the numbers that are matching
				if len(output) >= 4 {
					fmt.Printf("Found match with %d\t%d\t %s\n", a, a-lastMatch, output)
					lastMatch = a
				}
				return false
			}
			return true
		})

		if snapshot.StdoutString() == program {
			return a
		}
	}
	return -1
}

// computer reads an instruction set
//
// it will start with 3 registers A, B, C
//
// Combo operands 0 through 3 represent literal values 0 through 3.
// Combo operand 4 represents the value of register A.
// Combo operand 5 represents the value of register B.
// Combo operand 6 represents the value of register C.
// Combo operand 7 is reserved and will not appear in valid programs.
type computer struct {
	registers *datatypes.Inventory[string]

	instructionPointer int
	program            []int

	Stdout bytes.Buffer
	Rawout []int

	debug bool
}

func (c *computer) Image() *computer {
	return &computer{
		registers: c.registers.Copy(),

		instructionPointer: 0,
		program:            c.program,

		Stdout: bytes.Buffer{},
		Rawout: []int{},
	}
}

func (c *computer) Run(halt ...func(c *computer) (stop bool)) {
	// Allows intercepting of computer
	h := func(c *computer) bool {
		return false
	}
	if len(halt) > 0 {
		h = halt[0]
	}

	for c.instructionPointer < len(c.program) {
		instruction, operand := c.program[c.instructionPointer], c.program[c.instructionPointer+1]
		//fmt.Printf("ZST A=%d\tB=%d\tC=%d\n", c.registers.Count("A"), c.registers.Count("B"), c.registers.Count("C"))

		switch instruction {
		case 0: // Opcode 0 - division
			if !c.adv(operand) {
				continue
			}
		case 1: // Opcode 1 - bitwise XOR
			if !c.bxl(operand) {
				continue
			}
		case 2: // Opcode 2 - bst (unspecified behavior in your code comments)
			if !c.bst(operand) {
				continue
			}
		case 3: // Opcode 3 - jnz (jump if not zero)
			if !c.jnz(operand) {
				continue
			}
		case 4: // Opcode 4 - bxc (unspecified behavior)
			if !c.bxc(operand) {
				continue
			}
		case 5: // Opcode 5 - out (unspecified behavior)
			if !c.out(operand) {
				continue
			}
			if h(c) {
				return
			}
		case 6: // Opcode 6 - bdv (unspecified behavior)
			if !c.bdv(operand) {
				continue
			}
		case 7: // Opcode 7 - cdv (unspecified behavior)
			if !c.cdv(operand) {
				continue
			}
		default:
			panic("segfault: invalid instruction")
		}

		c.instructionPointer += 2
	}
}

func (c *computer) SetPointer(pointer int) {
	c.instructionPointer = pointer
}

// return bool indicates whether point should jump 2
func (c *computer) Opcode(opcode int, operand int) bool {
	return true
}

// Opcode 0 - division
//
//	numerator in A, denominator 2^operand
//	result is truncated to an integer
func (c *computer) adv(operand int) bool {
	value := c.Operand(operand)
	if c.debug {
		fmt.Printf("ADV A A/(2^%s)\n", c.DebugOperand(operand))
	}
	numerator := c.registers.Count("A")
	denominator := math.Pow(2, float64(value))
	result := int(math.Trunc(float64(numerator) / denominator))

	c.registers.SetX("A", result)
	return true
}

// Bitwise XOR of rB and operand
func (c *computer) bxl(operand int) bool {
	if c.debug {
		fmt.Printf("BXL B B⊻%d\n", operand)
	}
	c.registers.SetX("B", c.registers.Count("B")^operand)
	return true
}

// bst - modulus of operand % 8 and save to B
func (c *computer) bst(operand int) bool {
	if c.debug {
		fmt.Printf("BST B %s%%8\n", c.DebugOperand(operand))
	}
	c.registers.SetX("B", c.Operand(operand)%8)
	return true
}

// Nothing when A is 0
// If not 0, jump to operand
func (c *computer) jnz(operand int) bool {
	if c.registers.Count("A") == 0 {
		//fmt.Printf("JNZ A=%d\tB=%d\tC=%d\n", c.registers.Count("A"), c.registers.Count("B"), c.registers.Count("C"))
		return true
	}

	if c.debug {
		fmt.Printf("JNZ P %d\tA=%d\tB=%d\tC=%d\n", operand, c.registers.Count("A"), c.registers.Count("B"), c.registers.Count("C"))
	}

	c.SetPointer(operand)
	return false
}

// bxc - Bitwise XOR B ^ C into B
func (c *computer) bxc(operand int) bool {
	if c.debug {
		fmt.Printf("BXC B B⊻C\n")
	}
	c.registers.SetX("B", c.registers.Count("B")^c.registers.Count("C"))
	return true
}

// out - calc operand % 8 and outputs to buffer
func (c *computer) out(operand int) bool {
	if c.debug {
		fmt.Printf("\tOUT %s\n", c.DebugOperand(operand))
	}
	value := c.Operand(operand) % 8
	c.Rawout = append(c.Rawout, value)
	c.Stdout.WriteString(fmt.Sprintf("%d", value))
	return true
}

// bdv - same and adv for B register
func (c *computer) bdv(operand int) bool {
	if c.debug {
		fmt.Printf("BDV B A/(2^%s)\n", c.DebugOperand(operand))
	}
	numerator := c.registers.Count("A")
	denominator := math.Pow(2, float64(c.Operand(operand)))
	result := int(math.Trunc(float64(numerator) / denominator))

	c.registers.SetX("B", result)
	return true
}

// cdv - same as adv,bdv for C
func (c *computer) cdv(operand int) bool {
	if c.debug {
		fmt.Printf("CDV C A/(2^%s)\n", c.DebugOperand(operand))
	}
	numerator := c.registers.Count("A")
	denominator := math.Pow(2, float64(c.Operand(operand)))
	result := int(math.Trunc(float64(numerator) / denominator))

	c.registers.SetX("C", result)
	return true
}

func (c *computer) Operand(operand int) int {
	assert.LessThan(operand, 7)

	if operand >= 0 && operand <= 3 {
		return operand
	}

	switch operand {
	case 4:
		return c.registers.Count("A")
	case 5:
		return c.registers.Count("B")
	case 6:
		return c.registers.Count("C")
	default:
		panic("segfault: invalid operand")
	}
}

func (c *computer) DebugOperand(operand int) string {
	assert.LessThan(operand, 7)

	if operand >= 0 && operand <= 3 {
		return fmt.Sprintf("%d", operand)
	}

	switch operand {
	case 4:
		return "A"
	case 5:
		return "B"
	case 6:
		return "C"
	default:
		panic("segfault: invalid operand")
	}
}

func (c *computer) StdoutString() string {
	var t text.Texts

	for _, b := range c.Stdout.Bytes() {
		t = append(t, text.Text(fmt.Sprintf("%c", b)))
	}

	return t.Join(",").String()
}

func NewComputer(data []byte) *computer {
	registers := datatypes.NewInventory("Register")

	t := text.Text(data)
	config := t.Split("\n\n", 1)

	r, p := config[0], config[1]
	for _, line := range r.Lines() {
		c := line.Split(": ", 1)
		register, valueText := c[0], c[1]
		value, err := conversion.ToInt(valueText)
		assert.NoError(err)

		registers.SetX(string(register[len(register)-1:]), value)
	}
	c := conversion.To((p.Split(": ", 1)[1]).Split(",", -1), func(t text.Text) int {
		v, err := conversion.ToInt(t)
		assert.NoError(err)
		return v
	})

	return &computer{
		registers:          registers,
		instructionPointer: 0,
		program:            c,

		Stdout: bytes.Buffer{},
	}

}
