package main

import (
	"bufio"
	"math"
	"os"
	"regexp"
	"strconv"
)

var inputFile = "2024/day_two/input.txt"

type Report []int

var reNumbers = regexp.MustCompile(`\d+`)

func NewReport(line []byte) Report {
	// split line by spaces
	numbers := reNumbers.FindAll(line, -1)

	var report = make(Report, len(numbers))

	for i, number := range numbers {
		num, err := strconv.Atoi(string(number))
		if err != nil {
			panic("Error parsing integers from input")
		}
		report[i] = num
	}

	return report
}

var safetyThreshold = 3.0

func (r Report) CreateSlice(ignore int) Report {
	newReport := make(Report, len(r)-1)

	var x int
	for i := range r {
		if i == ignore {
			continue
		}
		newReport[x] = r[i]
		x++
	}
	return newReport
}

func Clamp(x, max int) int {
	if x > max {
		return max
	}
	return x
}

// safety check is change in direction
// level changes more than 3 values
func (r Report) CalculateSafety(problemDampener bool) bool {
	safe := true

	// test step size limit
	for i := 0; i < len(r)-1; i++ {
		if !(CheckDelta(r[i], r[i+1]) && CheckDirection(r[0:Clamp(i+3, len(r))])) {
			safe = false
		}

		if !safe && problemDampener {
			return r.CreateSlice(i).CalculateSafety(false) || r.CreateSlice(i+1).CalculateSafety(false)
		} else if !safe {
			return false
		}
	}

	return safe
}

func (r Report) CalculateSafetyBrute(problemDampener bool) bool {
	safe := false
	for x := range r {
		innerSafe := true

		inspectReport := r.CreateSlice(x)
		// test step size limit
		for i := 0; i < len(inspectReport)-1; i++ {
			if !(CheckDelta(inspectReport[i], inspectReport[i+1]) && CheckDirection(inspectReport[0:i+2])) {
				innerSafe = false
			}
		}

		safe = safe || innerSafe
	}

	return safe
}
func CheckDelta(a, b int) bool {
	delta := a - b
	return math.Abs(float64(delta)) <= safetyThreshold
}

func CheckDirection(r Report) bool {
	comparison := func(a, b int) bool {
		return b > a
	}
	if r[0] > r[1] {
		comparison = func(a, b int) bool {
			return b < a
		}
	}
	for i := 1; i < len(r); i++ {
		if !comparison(r[i-1], r[i]) {
			return false
		}
	}
	return true
}

func main() {
	file, err := os.OpenFile(inputFile, os.O_RDONLY, 0644)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	//safeFile, err := os.OpenFile("2024/day_two/safe_reports_2.txt", os.O_RDONLY, 0644)

	safeReports := 0
	unsafeReports := 0
	for line, _, err := reader.ReadLine(); err == nil; line, _, err = reader.ReadLine() {
		report := NewReport(line)

		if !report.CalculateSafetyBrute(true) {
			//print("unsafe\n")
			unsafeReports++
		} else {
			//print("safe\n")
			//fmt.Fprintln(safeFile, line)
			safeReports++
		}
	}

	println(safeReports)
	println(unsafeReports)
}
