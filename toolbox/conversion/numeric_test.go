package conversion

import (
	"adventofcode/toolbox/tests"
	"testing"
)

func TestToInt(t *testing.T) {
	tests.TestTables(
		t,
		tests.TestTable[any, int]{
			{
				Name:   "numeric to number",
				Input:  "4",
				Expect: 4,
			},
			{
				Name:   "non numeric to error",
				Input:  "48",
				Expect: 48,
			},
			{
				Name:   "bytes to numeric",
				Input:  []byte("29"),
				Expect: 29,
			},
		},
		func(test tests.Test[any, int], t *testing.T) {
			value, err := ToInt(test.GetInput())

			if err != nil {
				t.Errorf("did not get expected error, got %v, want %v", err, nil)
			}
			if value != test.GetExpecting() {
				t.Errorf("did not get expected output, got %d, want %d", value, test.GetExpecting())
			}
		},
	)
}
