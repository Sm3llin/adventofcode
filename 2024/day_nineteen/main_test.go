package main

import (
	"adventofcode/toolbox/tests"
	"testing"
)

var (
	example1 = `r, wr, b, g, bwu, rb, gb, br

brwrr
bggr
gbbr
rrbgbr
ubwu
bwurrg
brgr
bbrgwb`
)

func TestTowels_Valid(t *testing.T) {
	tests.TestTables(t, tests.TestTable[string, int]{
		{
			Name:   "example 1",
			Input:  example1,
			Expect: 6,
		},
	}, func(test tests.Test[string, int], t *testing.T) {
		designs, patterns := loadTowels([]byte(test.GetInput()))

		var validTowels int
		for _, design := range designs {
			allowed := valid(design, patterns)
			if allowed > 0 {
				validTowels++
			}
			t.Logf("%s: %d", design, allowed)
		}

		if validTowels != test.GetExpecting() {
			t.Errorf("did not get expected output, got %d, want %d", validTowels, test.GetExpecting())
		}
	})
}
