package adventofcode

func SingleByteToInt(data byte) int {
	var value int
	switch data {
	case '1':
		value = 1
	case '2':
		value = 2
	case '3':
		value = 3
	case '4':
		value = 4
	case '5':
		value = 5
	case '6':
		value = 6
	case '7':
		value = 7
	case '8':
		value = 8
	case '9':
		value = 9
	}
	return value
}
