package tests

import (
	"reflect"
	"testing"
)

func TestTables[I, T any](t *testing.T, tables TestTable[I, T], f ...func(test Test[I, T], t *testing.T)) {
	fallback := func(test Test[I, T], t *testing.T) {
		if test.GetResult == nil {
			t.Errorf("Test %q failed: no result", test.GetName())
			return
		}

		value, err := test.GetResult(test.GetInput())

		if test.Error != nil {
			if err == nil {
				t.Errorf("test %q did not throw error", test.GetName())
				return
			} else if test.ErrorMsg != "" && test.ErrorMsg != err.Error() {
				t.Errorf("test %q failed: expected error %s but got %s", test.GetName(), test.ErrorMsg, err.Error())
			}
		} else {
			if err != nil {
				t.Errorf("test %q failed: %s", test.GetName(), err.Error())
				return
			}

			if !reflect.DeepEqual(value, test.GetExpecting()) {
				t.Errorf("test %q failed: expected %v but got %v", test.GetName(), test.GetExpecting(), value)
			}
		}
	}

	run := fallback
	if len(f) >= 1 {
		run = f[0]
	}
	for _, table := range tables {
		t.Run(table.GetName(), func(t *testing.T) {
			run(table, t)
		})
	}
}

type Tester interface {
	GetName() string
	GetInput() string
}

type TestTable[I any, T any] []Test[I, T]
type Test[I, T any] struct {
	Name  string
	Input I

	Expect T

	// Should we have an auto option for test runners
	GetResult func(I) (T, error)

	// Error allow option
	Error    error
	ErrorMsg string
}

func (t Test[I, T]) GetName() string {
	return t.Name
}

func (t Test[I, T]) GetInput() I {
	return t.Input
}

func (t Test[I, T]) GetExpecting() T {
	return t.Expect
}

func (t Test[I, T]) GetError() error {
	return t.Error
}

func (t Test[I, T]) GetErrorMsg() string {
	return t.ErrorMsg
}
