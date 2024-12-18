package text

import (
	"adventofcode/toolbox/conversion"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

type Text string
type Texts []Text

// Lines split the text into Text lines
func (t Text) Lines() []Text {
	lines := strings.Split(string(t), "\n")
	texts := make([]Text, len(lines))

	for i, line := range lines {
		texts[i] = Text(line)
	}

	return texts
}

var (
	numericDigits = Texts{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	}
	allNumericDigits = Texts{
		// position is important for converting to int
		"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	}
	reNumericWord = regexp.MustCompile("(" + string(allNumericDigits.Join("|")) + ")")
)

func (t Text) String() string {
	return string(t)
}

func (t Text) Bytes() []byte {
	return []byte(t)
}

// DigitFrom get the digit of the matching word in text
func DigitFrom(word Text) Text {
	if len(word) == 1 && byte(word[0]) >= '0' && word[0] <= '9' {
		return word
	}

	return Text(fmt.Sprintf("%d", slices.Index(allNumericDigits, word)))
}

// FindDigits returns all numbers within the text block
func (t Text) FindDigits(nonNumeric ...bool) Texts {
	var includeNonNumeric bool
	if len(nonNumeric) > 0 && nonNumeric[0] {
		includeNonNumeric = nonNumeric[0]
	}

	digits := []Text{}
	if includeNonNumeric {
		matches := make(map[Text][]int)
		var matchCount int
		for _, digit := range allNumericDigits {
			for i := 0; i < len(t); {
				j := strings.Index(string(t[i:]), string(digit))

				if j >= 0 {
					if _, ok := matches[digit]; !ok {
						matches[digit] = []int{}
					}
					matches[digit] = append(matches[digit], j+i)
					i += j + 1
					matchCount++
				} else {
					break
				}
			}
		}

		m := make([]struct {
			digit Text
			index int
		}, matchCount)
		var i int
		for word, match := range matches {
			for _, index := range match {
				m[i] = struct {
					digit Text
					index int
				}{digit: word, index: index}
				i++
			}
		}

		slices.SortFunc(m, func(a, b struct {
			digit Text
			index int
		}) int {
			return a.index - b.index
		})

		return conversion.To(m, func(m struct {
			digit Text
			index int
		}) Text {
			if slices.Contains(numericDigits, m.digit) {
				return Text(fmt.Sprintf("%d", slices.Index(numericDigits, m.digit)))
			}

			return Text(fmt.Sprintf("%d", slices.Index(allNumericDigits, m.digit)))
		})
	}

	for _, c := range t {
		if c >= '0' && c <= '9' {
			digits = append(digits, Text(c))
		}
	}
	return digits
}

func (t Text) TrimSpace() Text {
	return Text(strings.TrimSpace(string(t)))
}

func (t Text) Split(sep string, count int) Texts {
	return conversion.To(strings.Split(string(t), sep), func(s string) Text {
		return Text(s)
	})
}

func (t Texts) Join(sep string) Text {
	var text Text
	if len(t) == 1 {
		return t[0]
	} else {
		for i, line := range t {
			if i > 0 {
				text += Text(sep)
			}
			text += line
		}
	}
	return text
}

func (t Texts) Trim() Texts {
	return conversion.To(t, func(s Text) Text {
		return s.TrimSpace()
	})
}
