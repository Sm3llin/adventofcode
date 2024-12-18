package assert

import (
	"fmt"
)

func Assert(cond bool) {
	if !cond {
		panic("assertion failed")
	}
}

func Assertf(cond bool, format string, args ...interface{}) {
	if !cond {
		panic(fmt.Sprintf(format, args...))
	}
}

func Equal[T comparable](a, b T) {
	if a != b {
		panic(fmt.Sprintf("expected %v, got %v", a, b))
	}
}

func NotEqual[T comparable](a, b T) {
	if a == b {
		panic(fmt.Sprintf("expected %v, got %v", a, b))
	}
}

func LessThan[T int | int64 | byte | float64](a, b T) bool {
	return a < b
}
func Length[T any](v []T) {
	if len(v) == 0 {
		panic("expected non-empty slice")
	}
}

func NotLength[T any](v []T) {
	if len(v) != 0 {
		panic("expected empty slice")
	}
}

func LengthMin[T any](v []T, min int) {
	if len(v) < min {
		panic(fmt.Sprintf("expected slice length >= %d", min))
	}
}

func NoError(err error) {
	if err != nil {
		panic(err)
	}
}

func NotNil(v any) {
	if v == nil {
		panic("expected non-nil value")
	}
}

func Nil(v any) {
	if v != nil {
		panic("expected nil value")
	}
}
