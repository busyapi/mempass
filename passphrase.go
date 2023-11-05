package mempass

import (
	"math/rand"
	"regexp"
	"strings"
	"unicode"
)

type FromPassphrase struct {
	l33t *L33t
}

func NewFromPassphrase() *FromPassphrase {
	return &FromPassphrase{l33t: NewL33t()}
}

func (f *FromPassphrase) Generate(input string) []rune {
	input = f.preProcessPassphrase(input)
	l, ucCount, numCount, scCount, lcPos := f.countChars(input)

	// Minimum number of uppercase, numbers and specials chars
	ratio := 1.0 / 8.0
	min := int(float64(l) * ratio)
	if min == 0 {
		min = 1
	}

	// Minimum number of uppercase, numbers and specials chars to add
	addUc := min - ucCount
	addNum := min - numCount
	addSc := min - scCount

	// Create a RW rune array from the input string
	runes := toRunes(input)
	runes = f.addUc(runes, addUc, lcPos)
	runes = f.addNums(runes, addNum, lcPos)
	runes = f.addSc(runes, addSc)

	return runes
}

func (f *FromPassphrase) countChars(s string) (l, uc, num, sc int, lcPos []int) {
	runes := toRunes(s)
	l = len(runes)

	for i, char := range runes {
		switch {
		case unicode.IsLower(char):
			lcPos = append(lcPos, i)
		case unicode.IsNumber(char):
			num++
		case !unicode.IsLetter(char) && !unicode.IsNumber(char):
			sc++
		}
	}

	sc = sc - uc - num

	return
}

func (f *FromPassphrase) preProcessPassphrase(input string) string {
	// First replace all spaces by an hyphen
	input = strings.ReplaceAll(input, " ", "-")

	// Then, replace all multiple hyphens with a single hyphen.
	regMultipleHyphens := regexp.MustCompile(`\-+`)
	input = regMultipleHyphens.ReplaceAllString(input, "-")

	return input
}

func (f *FromPassphrase) addUc(runes []rune, count int, lcPos []int) []rune {
	done := 0

	if len(lcPos) > 0 {
		for i := 0; i < count; i++ {
			// Pick a random run
			idx := rand.Intn(len(lcPos))
			pos := lcPos[idx]
			char := runes[pos]

			// Don't process chars that are not lowercase letter
			if !unicode.IsLower(char) {
				continue
			}

			// Transform the char to uppercase
			runes[pos] = unicode.ToUpper(char)

			// Remove the uppercased character from the lc array
			lcPos = append(lcPos[:idx], lcPos[idx+1:]...)

			done++
		}
	}

	if done < count {
		for i := done + 1; i <= count; i++ {
			idx := rand.Intn(26)
			runes = append(runes, rune(ALPHABET_UPPER[idx]))
		}
	}

	return runes
}

func (f *FromPassphrase) addNums(runes []rune, count int, lcPos []int) []rune {
	done := 0

	// Find all l33table characters positions
	l33table := f.find1337able(runes)

	if len(l33table) > 0 {
		for i := 0; i < count; i++ {
			// Get a random position from the l33table characters positions array
			idx := rand.Intn(len(l33table))
			pos := l33table[idx]

			// Transform the character
			runes[pos] = rune(f.l33t.make1337(runes[pos], 0))

			// Remove the l33ted character from the l33table array
			l33table = append(l33table[:idx], l33table[idx+1:]...)

			done++
		}
	}

	if done < count {
		for i := done + 1; i <= count; i++ {
			idx := rand.Intn(10)
			runes = append(runes, rune(NUMBERS[idx]))
		}
	}

	return runes
}

func (f *FromPassphrase) addSc(runes []rune, count int) []rune {
	for i := 0; i < count; i++ {
		runes = append(runes, '-')
	}

	return runes
}

func (f *FromPassphrase) find1337able(input []rune) (l33table []int) {
	for i, char := range input {
		if f.l33t.can1337(char) {
			l33table = append(l33table, i)
		}
	}

	return l33table
}
