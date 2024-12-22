package main

import (
	"adventofcode"
	"adventofcode/toolbox/arrays"
	"adventofcode/toolbox/conversion"
	"adventofcode/toolbox/fs"
	"adventofcode/toolbox/text"
	"bytes"
	"fmt"
	"slices"
)

func main() {
	adventofcode.Time(func() {
		designs, patterns := loadTowels(fs.LoadFile("2024/day_nineteen/input.txt"))

		var total int
		var value int
		for i, design := range designs {
			fmt.Printf("%d/%d \r", i+1, len(designs))
			if n := valid(design, patterns); n > 0 {
				value += n
				total++
			}
		}
		fmt.Println()
		fmt.Printf("Part 1: %d\n", total)
		fmt.Printf("Part 2: %d\n", value)
	})
}

func loadTowels(data []byte) ([]text.Text, []text.Text) {
	d := bytes.Split(data, []byte("\n\n"))
	patternsB, designB := d[0], d[1]
	return text.Text(designB).TrimSpace().Lines(), text.Text(patternsB).TrimSpace().Split(",", -1).Trim()
}

func valid(design text.Text, patterns text.Texts) int {
	d := design.Bytes()
	p := arrays.Filter(patterns.Bytes(), func(b []byte) bool {
		return bytes.Contains(d, b)
	})

	// sort the pattern based on its position in the design
	slices.SortFunc(p, func(a, b []byte) int {
		return len(b) - len(a)
	})
	if ok := check(d, p); ok > 0 {
		return ok
	}

	return 0
}

type match struct {
	pattern       []byte
	index         int
	length, shiny int
	valid         bool
}

func (m match) String() string {
	return fmt.Sprintf("match(idx=%d, len=%d)", m.index, m.length)
}

func (m match) Bits(length int) int {
	n := 1 << (length - m.index - 1)
	j := n
	for range m.length - 1 {
		j = j >> 1
		n += j
	}
	return n
}

func check(design []byte, patterns [][]byte) int {
	if len(design) == 0 {
		return 1
	}
	var matches []*match
	for _, pattern := range patterns {
		var pos int
		var m int
		for m != -1 && pos < len(design) {
			m = bytes.Index(design[pos:], pattern)
			if m != -1 {
				matches = append(matches, &match{
					pattern: pattern,
					index:   pos + m,
					length:  len(pattern),
					valid:   true,
				})
				pos += m + 1
			}
		}
	}
	// create a hash map for the index and then test each solution individually grabbing the next index solutions
	lookup := make(map[int][]*match)
	for _, match := range matches {
		lookup[match.index] = append(lookup[match.index], match)
	}

	_, ok := lookup[0]

	if !ok {
		// no matching start option
		return 0
	}

	designBit := (1 << len(design)) - 1
	bits := conversion.To(matches, func(m *match) int {
		return m.Bits(len(design))
	})
	// check that the number can be created
	patternBit := 0
	for _, b := range bits {
		patternBit = patternBit | b
	}
	// simple check if there is actually enough matches present to try
	if patternBit != designBit {
		return 0
	}

	// might need to do a cleanse on achievable positions at end
	validPositions := []int{len(design)}
	for i := len(design) - 1; i >= 0; i-- {
		indexMatches, ok := lookup[i]
		if !ok {
			continue
		}
		for _, m := range indexMatches {
			if slices.Contains(validPositions, m.index+m.length) {
				m.shiny += arrays.Count(validPositions, m.index+m.length)
				validPositions = append(validPositions, m.index)
			} else {
				m.valid = false
			}
		}
	}

	// opposite of the filter end, filter start to ensure that all matches are reachable
	validSteps := []int{0}
	for i := 0; i < len(validPositions); i++ {
		indexMatches, ok := lookup[i]
		if !ok {
			continue
		}
		for _, m := range indexMatches {
			if !slices.Contains(validPositions, m.index) {
				validSteps = append(validSteps, m.index+m.length)
			}
		}
	}
	// disable matches without a valid step
	for _, m := range matches {
		if !slices.Contains(validPositions, m.index) {
			fmt.Printf("disabling %s @ index %d", m.pattern, m.index)
			m.valid = false
		}
	}

	lookup = make(map[int][]*match)
	for _, m := range matches {
		if m.valid {
			lookup[m.index] = append(lookup[m.index], m)
		}
	}

	//for _, m := range startingMatches {
	//	if ok = findBit(len(design), m, lookup); ok {
	//		return true
	//	}
	//}
	// create a graph

	startingMatches, ok := lookup[0]
	if !ok {
		// no matching start option
		return 0
	}

	bitCache := make(map[int]int)
	var totalDesigns int
	for _, m := range startingMatches {
		// TODO: need to implement a reverse cache on the numbers
		totalDesigns += findBit(len(design), m, lookup, bitCache)
	}
	return totalDesigns

	totalDesigns = len(arrays.Filter(startingMatches, func(m *match) bool {
		return m.valid
	}))
	for i, indexMatches := range lookup {
		if len(indexMatches) == 1 || i == 0 || i == len(startingMatches)-1 {
			for _, m := range indexMatches {
				fmt.Printf("%s\n", m.pattern)
			}
			continue
		}

		for _, m := range indexMatches {
			destMatches := lookup[m.index+m.length]

			if m.valid {
				fmt.Printf("%s\n", m.pattern)
				totalDesigns += len(destMatches)
			}
		}
	}

	fmt.Println()
	return totalDesigns

	//designBit := (1 << len(design)) - 1
	//bits := conversion.To(matches, func(m match) int {
	//	return m.Bits(len(design))
	//})
	//_ = designBit
	//_ = bits
	//
	//slices.SortFunc(bits, func(a, b int) int {
	//	return b - a
	//})
	//
	//for i := range bits {
	//	if ok := checkBit(designBit, bits[i], bits); ok {
	//		return true
	//	}
	//}
	//return false

	//// find all locations of a pattern
	//for _, pattern := range patterns {
	//	if bytes.HasPrefix(design, pattern) && check(design[len(pattern):], patterns) {
	//		return true
	//	}
	//}
	//return false
}

func findBit(length int, currentMatch *match, lookup map[int][]*match, bitCache map[int]int) int {
	if currentMatch.index+currentMatch.length == length {
		return 1
	} else if n, ok := bitCache[currentMatch.Bits(length)]; ok {
		return n
	}

	//fmt.Printf("%d/%d \r", currentMatch.index, length)
	nextMatches, ok := lookup[currentMatch.index+currentMatch.length]

	if !ok {
		return 0
	}
	var t int
	for _, m := range nextMatches {
		if !m.valid {
			continue
		}
		if n := findBit(length, m, lookup, bitCache); n > 0 {
			t += n
		}
	}

	bitCache[currentMatch.Bits(length)] = t

	return t
}

func checkBit(design int, pattern int, bits []int) bool {
	if design == pattern {
		return true
	}
	bits = arrays.Filter(bits, func(b int) bool {
		return pattern&b == 0
	})
	for _, b := range bits {
		if b&pattern == 0 {
			if ok := checkBit(design, b+pattern, bits); ok {
				return true
			}
		}
	}

	return false
}
