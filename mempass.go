package mempass

import (
	"errors"
	"math"
	"math/rand"
)

type CapRule string
type SepRule string
type SymbRule string
type SymbPos string
type PadRule string

const (
	CapRuleNone              CapRule = ""
	CapRuleAll               CapRule = "all"
	CapRuleAlternate         CapRule = "alternate"
	CapRuleWordAlternate     CapRule = "word_alternate"
	CapRuleFirstLetter       CapRule = "first_letter"
	CapRuleLastLetter        CapRule = "last_letter"
	CapRuleAllButFirstLetter CapRule = "all_but_first_letter"
	CapRuleAllButLastLetter  CapRule = "all_but_last_letter"
	CapRuleRandom            CapRule = "random"
)

const (
	SepRuleNone   SepRule = ""
	SepRuleFixed  SepRule = "fixed"
	SepRuleRandom SepRule = "random"
)

const (
	SymbRuleFixed  SymbRule = "fixed"
	SymbRuleRandom SymbRule = "random"
)

const (
	PadRuleFixed  PadRule = "fixed"
	PadRuleRandom PadRule = "random"
)

type Options struct {
	UseDict       bool     // Use dictionary. Default false
	WordCount     uint     // Number of words to generate. Using less than 2 is discouraged. Default is 2
	MinWordLength uint     // Minimum word length. O = no minimum. Using less than 4 is discouraged. Default is 0
	MaxWordLength uint     // Maximum word length. O = no maximum. Default is 0
	DigitsAfter   uint     // Number of digits to add at the end of each word. Default is 0
	DigitsBefore  uint     // Number of digits to add at the begining of each word. Default is 0
	CapRule       CapRule  // Capitalization rule
	CapRatio      float32  // Uppercase ratio. 0.0 = no uppercase, 1.0 = all uppercase, 0.3 = 1/3 uppercase, etc. Only used if `Capitalization` is `CapRandom`. Default is 0.2
	SymbRule      SymbRule // Rule for adding symbols. Default is `SymbRuleNone`
	SymbolsAfter  uint     // Number of symbols to add at the end of each word. Default is 0
	SymbolsBefore uint     // Number of symbols to add at the begining of each word. Default is 0
	SymbolPool    string   // Symbols pool. Only used if `SymbRule` is `SymbRuleRandom`. Default is "@&!-_^$*%,.;:/=+"
	Symbol        byte     // Symbol character. Only used if `SymbRule` is `SymbRuleFixed` or `SymbRulePadding`. Default is `/`
	SepRule       SepRule  // Seperator type. Default is `SepRuleNone`
	SeparatorPool string   // Seperators pool. Only used if `SepRule` is `SepRuleRandom`. Default is "@&!-_^$*%,.;:/=+"
	Separator     byte     // Separator for words. Only used if `SepRule` is `SepRuleFixed`. Default is '-'
	PadRule       PadRule  // Padding rule. Ignored if `PadLength` is 0
	PadSymbol     byte     // Padding symbol. Only used if `PadRule` si `PadRuleFixed`. Default is `.`
	PadLength     uint     // Password length to reach with padding.
}

type Generator struct {
	opt         *Options
	words       [][]byte
	size        uint
	paddingSize uint
}

func NewGenerator(opt *Options) Generator {
	return Generator{opt: opt}
}

// Generate a human memorable password
func (g *Generator) GenPassword() (string, float64, error) {
	if err := g.checkOptions(); err != nil {
		return "", 0, err
	}

	var words [][]byte
	var err error

	if g.opt.UseDict {
		if words, err = getDictWords(g.opt); err != nil {
			return "", 0, err
		}
	} else {
		words = genRandPwd(g.opt)
	}

	if g.opt.CapRule != CapRuleNone {
		words = g.capitalize(words)
	}

	g.words = g.padWord(words)

	for _, word := range g.words {
		g.size += uint(len(word))
	}

	var sep byte
	if g.opt.SepRule != SepRuleNone {
		if g.opt.SepRule == SepRuleFixed {
			sep = g.opt.Separator
		} else if g.opt.SepRule == SepRuleRandom {
			sep = g.randBytesFrom(1, g.opt.SeparatorPool)[0]
		}

		g.size += uint(len(g.words) - 1)
	}

	if g.opt.PadLength > 0 && g.size < g.opt.PadLength {
		g.paddingSize = g.opt.PadLength - g.size
		g.size += (g.opt.PadLength - g.size)
	}

	pwd := make([]byte, g.size)
	idx := 0

	for i, word := range g.words {
		copy(pwd[idx:], word)

		idx += len(word)

		if g.opt.SepRule != SepRuleNone {
			if i < len(words)-1 {
				pwd[idx] = sep
				idx++
			}
		}
	}

	if g.paddingSize >= 1 {
		pwd = g.addWordPadding(pwd, 0, g.paddingSize, g.opt.SymbolPool, g.opt.PadSymbol)
	}

	return string(pwd), g.entropy(string(pwd)), nil
}

func (g *Generator) padWord(words [][]byte) [][]byte {
	newWords := make([][]byte, len(words))

	for i, word := range words {
		newWord := word

		if g.opt.DigitsBefore > 0 || g.opt.DigitsAfter > 0 {
			newWord = g.addNumsPadding(newWord, g.opt.DigitsBefore, g.opt.DigitsAfter)
		}

		if g.opt.SymbolsBefore > 0 || g.opt.SymbolsAfter > 0 {
			newWord = g.addSymbolsPadding(newWord, g.opt.SymbolsBefore, g.opt.SymbolsAfter, g.opt.SymbolPool, g.opt.Symbol)
		}

		newWords[i] = newWord
	}

	return newWords
}

func (g *Generator) addNumsPadding(word []byte, nb uint, na uint) []byte {
	source := "0123456789"
	return g.addWordPadding(word, nb, na, source, 0)
}

func (g *Generator) addSymbolsPadding(word []byte, nb uint, na uint, source string, char byte) []byte {
	return g.addWordPadding(word, nb, na, source, char)
}

func (g *Generator) addWordPadding(word []byte, nb uint, na uint, source string, char byte) []byte {
	newSize := len(word) + int(nb) + int(na)
	newWord := make([]byte, newSize)
	copyPos := 0

	if nb > 0 {
		pad := g.padding(nb, char, source)
		copy(newWord[copyPos:], pad)
		copyPos += int(nb)
	}

	copy(newWord[copyPos:], word)
	copyPos += len(word)

	if na > 0 {
		pad := g.padding(na, char, source)
		copy(newWord[copyPos:], pad)
	}

	return newWord
}

func (g *Generator) padding(count uint, char byte, source string) []byte {
	if char == 0 {
		return g.randBytesFrom(count, source)
	}

	return g.paddingOfByte(count, char)
}

func (g *Generator) randBytesFrom(count uint, source string) []byte {
	res := make([]byte, count)

	for i := 0; i < int(count); i++ {
		idx := rand.Intn(10)
		res[i] = source[idx]
	}

	return res
}

func (g *Generator) paddingOfByte(count uint, char byte) []byte {
	padding := make([]byte, count)

	for i := range padding {
		padding[i] = char
	}

	return padding
}

func (g *Generator) capitalize(words [][]byte) [][]byte {
	newWords := make([][]byte, len(words))
	var newWord []byte

	for i := range words {
		switch g.opt.CapRule {
		case CapRuleAll:
			newWord = g.arrayMap(words[i], g.capChar)

		case CapRuleWordAlternate:
			if i%2 == 0 {
				newWord = g.arrayMap(words[i], g.capChar)
			} else {
				newWord = words[i]
			}

		case CapRuleAlternate:
			newWord = g.arrayMapIf(words[i], g.isAlt, g.capChar)

		case CapRuleFirstLetter:
			newWord = g.arrayMapIf(words[i], g.isFirstLetter, g.capChar)

		case CapRuleLastLetter:
			newWord = g.arrayMapIf(words[i], g.isLastLetter, g.capChar, len(words[i]))

		case CapRuleAllButFirstLetter:
			newWord = g.arrayMapIf(words[i], g.isNotFirstLetter, g.capChar)

		case CapRuleAllButLastLetter:
			newWord = g.arrayMapIf(words[i], g.isNotLastLetter, g.capChar, len(words[i]))

		case CapRuleRandom:
			newWord = g.arrayMapIf(words[i], g.isRand, g.capChar, g.opt.CapRatio)
		}

		newWords[i] = newWord
	}

	return newWords
}

func (g *Generator) capChar(char byte, idx int) byte {
	return char - 32
}

func (g *Generator) isAlt(char byte, idx int, _ ...any) bool {
	return idx%2 == 0
}

func (g *Generator) isFirstLetter(char byte, idx int, _ ...any) bool {
	return idx == 0
}

func (g *Generator) isNotFirstLetter(char byte, idx int, _ ...any) bool {
	return idx != 0
}

func (g *Generator) isLastLetter(char byte, idx int, o ...any) bool {
	return idx == o[0].(int)-1
}

func (g *Generator) isNotLastLetter(char byte, idx int, o ...any) bool {
	return idx != o[0].(int)-1
}

func (g *Generator) isRand(char byte, idx int, o ...any) bool {
	return rand.Float32() <= o[0].(float32)
}

func (g *Generator) arrayMap(slice []byte, fn func(byte, int) byte) []byte {
	result := make([]byte, len(slice))

	for i, v := range slice {
		result[i] = fn(v, i)
	}

	return result
}

func (g *Generator) arrayMapIf(slice []byte, ifFn func(byte, int, ...any) bool, fn func(byte, int) byte, ifFnArgs ...any) []byte {
	result := make([]byte, len(slice))

	for i, v := range slice {
		if ifFn(v, i, ifFnArgs...) {
			result[i] = fn(v, i)
		} else {
			result[i] = v
		}

	}

	return result
}

func (g *Generator) checkOptions() error {
	if g.opt.WordCount == 0 {
		g.opt.WordCount = 2
	}

	if g.opt.SeparatorPool == "" {
		g.opt.SeparatorPool = "@&!-_^$*%,.;:/=+"
	}

	if g.opt.SymbolPool == "" {
		g.opt.SymbolPool = "@&!-_^$*%,.;:/=+"
	}

	if g.opt.MaxWordLength == 0 {
		g.opt.MaxWordLength = 28
	}

	if g.opt.MinWordLength > 28 || g.opt.MaxWordLength > 28 {
		return errors.New("`MinWordLength` and `MaxWordLength` cannot be greater than 28")
	}

	if g.opt.MinWordLength > 0 && g.opt.MaxWordLength > 0 && g.opt.MinWordLength > g.opt.MaxWordLength {
		return errors.New("`MinWordLength` cannot be greater than `MaxWordLength`")
	}

	if g.opt.CapRule == CapRuleRandom && (g.opt.CapRatio <= 0 || g.opt.CapRatio >= 1) {
		return errors.New("`CapRatio` must be between 0 and 1 excluded")
	}

	if g.opt.SymbRule == SymbRuleFixed && g.opt.Symbol == 0 {
		g.opt.Symbol = '/'
	}

	if g.opt.SepRule == SepRuleFixed && g.opt.Separator == 0 {
		g.opt.Separator = '-'
	}

	if g.opt.PadRule == PadRuleFixed && g.opt.PadSymbol == 0 {
		g.opt.PadSymbol = '.'
	}

	return nil
}

func (g *Generator) entropy(pass string) float64 {
	charRange := 26
	var usedSymbols string

	if g.opt.CapRule != CapRuleNone {
		charRange *= 2
	}

	if g.opt.DigitsAfter > 0 || g.opt.DigitsBefore > 0 {
		charRange += 10
	}

	if g.opt.SymbolsAfter > 0 || g.opt.SymbolsBefore > 0 {
		if g.opt.SymbRule == SymbRuleFixed {
			usedSymbols = string(g.opt.Symbol)
		} else {
			usedSymbols = g.opt.SymbolPool
		}
	}

	if g.opt.SepRule == SepRuleFixed {
		usedSymbols = g.mergeStrings(usedSymbols, string(g.opt.Separator))
	} else if g.opt.SepRule == SepRuleRandom {
		usedSymbols = g.mergeStrings(usedSymbols, g.opt.SeparatorPool)
	}

	if g.opt.PadRule == PadRuleFixed {
		usedSymbols = g.mergeStrings(usedSymbols, string(g.opt.PadSymbol))
	} else if g.opt.PadRule == PadRuleRandom {
		usedSymbols = g.mergeStrings(usedSymbols, g.opt.SymbolPool)
	}

	charRange += len(usedSymbols)
	len := float64(len(pass))

	return math.Log2(float64(math.Pow(float64(charRange), len)))
}

func (g *Generator) mergeStrings(dest, src string) string {
	// Create a map to store characters in the destination string
	destChars := make(map[rune]bool)

	// Populate the map with characters from the destination string
	for _, char := range dest {
		destChars[char] = true
	}

	// Create a result string
	result := dest

	// Append characters from the source string that are not in the destination
	for _, char := range src {
		if !destChars[char] {
			result += string(char)
		}
	}

	return result
}
