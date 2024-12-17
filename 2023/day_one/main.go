package main

import (
	"adventofcode"
	"adventofcode/toolbox/assert"
	"adventofcode/toolbox/conversion"
	"adventofcode/toolbox/fs"
	"adventofcode/toolbox/text"
	"fmt"
)

func main() {
	adventofcode.Time(func() {
		fmt.Printf("Part 1: %d\n", CalibrationScore(fs.LoadFile("2023/day_one/input.txt"), false))
	})
	adventofcode.Time(func() {
		fmt.Printf("Part 2: %d\n", CalibrationScore(fs.LoadFile("2023/day_one/input.txt"), true))
	})
}

func CalibrationScore(d []byte, nonNumeric bool) int {
	data := text.Text(d)

	var total int
	// for each line calculate the digits
	for _, line := range data.Lines() {
		digits := line.FindDigits(nonNumeric)

		assert.LengthMin(digits, 1)
		concatDigits := digits[0] + digits[len(digits)-1]

		n, _ := conversion.ToInt(concatDigits)

		total += n
	}

	return total
}
