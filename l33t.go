package mempass

import "unicode"

type L33t struct {
	leetMap map[rune]rune
}

func NewL33t() *L33t {
	leetMap := make(map[rune]rune)
	leetMap['a'] = '4'
	leetMap['A'] = '4'
	leetMap['e'] = '3'
	leetMap['E'] = '3'
	leetMap['i'] = '1'
	leetMap['I'] = '1'
	leetMap['o'] = '0'
	leetMap['O'] = '0'
	leetMap['s'] = '5'
	leetMap['S'] = '5'
	leetMap['t'] = '7'
	leetMap['T'] = '7'

	return &L33t{leetMap}
}

func (l *L33t) can1337(char rune) bool {
	if !unicode.IsLower((char)) {
		return false
	}

	_, exists := l.leetMap[char]

	return exists
}

func (l *L33t) make1337(char rune, idx int) rune {
	if _, exists := l.leetMap[char]; exists {
		return l.leetMap[char]
	}

	return char
}
