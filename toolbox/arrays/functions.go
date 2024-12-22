package arrays

func Filter[S any](s []S, f func(S) bool) []S {
	matches := make([]bool, len(s))

	for i, v := range s {
		matches[i] = f(v)
	}

	var j int
	var v bool
	for _, v = range matches {
		if v {
			j++
		}
	}

	r := make([]S, j)
	j = 0
	for i, v := range matches {
		if v {
			r[j] = s[i]
			j++
		}
	}

	return r
}

func Count[S comparable](haystack []S, needle S) int {
	var count int
	for _, v := range haystack {
		if v == needle {
			count++
		}
	}
	return count
}
