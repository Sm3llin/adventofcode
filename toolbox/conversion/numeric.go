package conversion

import (
	"fmt"
	"strconv"
)

func To[I, O any](s []I, f func(i I) O) []O {
	o := make([]O, len(s))
	for x := range s {
		o[x] = f(s[x])
	}
	return o
}

func ToInt[T any](v T) (int, error) {
	switch val := any(v).(type) {
	case byte:
		return strconv.Atoi(string(val))
	case int:
		return val, nil
	case float64:
		return int(val), nil
	case fmt.Stringer:
		return strconv.Atoi(val.String())
	case string:
		return strconv.Atoi(val)
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case []byte:
		return strconv.Atoi(string(val))
	default:
		return 0, fmt.Errorf("unsupported ToInt type: %T", v)
	}
}
