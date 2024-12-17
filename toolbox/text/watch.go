package text

type WordMatch func(byte) (WordMatch, bool, Text)

func WatchWord(word Text, pos int) WordMatch {
	return func(c byte) (WordMatch, bool, Text) {
		if c == word[pos] {
			if pos == len(word)-1 {
				return nil, true, word
			}
			return WatchWord(word, pos+1), false, word
		}
		return nil, false, word
	}
}
