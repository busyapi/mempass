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
	CapRuleNone              CapRule = "none"
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
	SepRuleNone   SepRule = "none"
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
	WordCount     uint     // Number of words to generate. Using less than 2 is discouraged. Default is 3
	MinWordLength uint     // Minimum word length. O = no minimum. Using less than 4 is discouraged. Default is 6
	MaxWordLength uint     // Maximum word length. O = no maximum. Default is 6
	DigitsAfter   uint     // Number of digits to add at the end of each word. Default is 0
	DigitsBefore  uint     // Number of digits to add at the begining of each word. Default is 0
	CapRule       CapRule  // Capitalization rule
	CapRatio      float32  // Uppercase ratio. 0.0 = no uppercase, 1.0 = all uppercase, 0.3 = 1/3 uppercase, etc. Only used if `CapRule` is `CapRandom`. Default is 0.2
	SymbRule      SymbRule // Rule for adding symbols. Default is `SymbRuleNone`
	SymbolsAfter  uint     // Number of symbols to add at the end of each word. Default is 0
	SymbolsBefore uint     // Number of symbols to add at the begining of each word. Default is 0
	SymbolPool    string   // Symbols pool. Only used if `SymbRule` is `SymbRuleRandom`. Default is "@&!-_^$*%,.;:/=+"
	Symbol        byte     // Symbol character. Only used if `SymbRule` is `SymbRuleFixed`. Default is `/`
	SepRule       SepRule  // Seperator type. Default is `SepRuleFixed`
	SeparatorPool string   // Seperators pool. Only used if `SepRule` is `SepRuleRandom`. Default is "@&!-_^$*%,.;:/=+"
	Separator     byte     // Separator for words. Only used if `SepRule` is `SepRuleFixed`. Default is '-'
	PadRule       PadRule  // Padding rule. Ignored if `PadLength` is 0
	PadSymbol     byte     // Padding symbol. Only used if `PadRule` si `PadRuleFixed`. Default is `.`
	PadLength     uint     // Password length to reach with padding.
	L33tRatio     float32  // 1337 coding ratio. 0.0 = no 1337, 1.0 = all 1337, 0.3 = 1/3 1337, etc`. Default is 0
}

type Generator struct {
	opt         *Options
	words       [][]byte
	size        uint
	paddingSize uint
	leetMap     map[byte]byte
}

func NewGenerator(opt *Options) Generator {
	leetMap := make(map[byte]byte)
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

	return Generator{opt: opt, leetMap: leetMap}
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

	g.words = g.extraProcess(words)

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
		g.size += (g.opt.PadLength - g.size)
	}

	return string(pwd), g.entropy(string(pwd)), nil
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
		idx := rand.Intn(len(source))
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

func (g *Generator) extraProcess(words [][]byte) [][]byte {
	newWords := make([][]byte, len(words))

	for i, word := range words {
		newWord := word

		if g.opt.CapRule != CapRuleNone {
			newWord = g.capWord(newWord, i)
		}

		if g.opt.DigitsBefore > 0 || g.opt.DigitsAfter > 0 {
			newWord = g.addNumsPadding(newWord, g.opt.DigitsBefore, g.opt.DigitsAfter)
		}

		if g.opt.SymbolsBefore > 0 || g.opt.SymbolsAfter > 0 {
			newWord = g.addSymbolsPadding(newWord, g.opt.SymbolsBefore, g.opt.SymbolsAfter, g.opt.SymbolPool, g.opt.Symbol)
		}

		if g.opt.L33tRatio > 0 {
			newWord = g.arrayMapIf(newWord, g.isRand, g.make1337, g.opt.L33tRatio)
		}

		g.size += uint(len(newWord))
		newWords[i] = newWord
	}

	return newWords
}

func (g *Generator) capWord(word []byte, i int) []byte {
	var newWord []byte

	switch g.opt.CapRule {
	case CapRuleAll:
		newWord = g.arrayMap(word, g.capChar)

	case CapRuleWordAlternate:
		if i%2 == 0 {
			newWord = g.arrayMap(word, g.capChar)
		} else {
			newWord = word
		}

	case CapRuleAlternate:
		newWord = g.arrayMapIf(word, g.isAlt, g.capChar)

	case CapRuleFirstLetter:
		newWord = g.arrayMapIf(word, g.isFirstLetter, g.capChar)

	case CapRuleLastLetter:
		newWord = g.arrayMapIf(word, g.isLastLetter, g.capChar, len(word))

	case CapRuleAllButFirstLetter:
		newWord = g.arrayMapIf(word, g.isNotFirstLetter, g.capChar)

	case CapRuleAllButLastLetter:
		newWord = g.arrayMapIf(word, g.isNotLastLetter, g.capChar, len(word))

	case CapRuleRandom:
		newWord = g.arrayMapIf(word, g.isRand, g.capChar, g.opt.CapRatio)
	}

	return newWord
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

func (g *Generator) make1337(char byte, idx int) byte {
	if _, exists := g.leetMap[char]; exists {
		return g.leetMap[char]
	}

	return char
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
		g.opt.WordCount = 3
	}

	if g.opt.MinWordLength == 0 {
		g.opt.MinWordLength = 6
	}

	if g.opt.MaxWordLength == 0 {
		g.opt.MaxWordLength = 6
	}

	if g.opt.SeparatorPool == "" {
		g.opt.SeparatorPool = "@&!-_^$*%,.;:/=+"
	}

	if g.opt.SymbolPool == "" {
		g.opt.SymbolPool = "@&!-_^$*%,.;:/=+"
	}

	if g.opt.MinWordLength > 28 || g.opt.MaxWordLength > 28 {
		return errors.New("`MinWordLength` and `MaxWordLength` cannot be greater than 28")
	}

	if g.opt.MinWordLength > 0 && g.opt.MaxWordLength > 0 && g.opt.MinWordLength > g.opt.MaxWordLength {
		return errors.New("`MinWordLength` cannot be greater than `MaxWordLength`")
	}

	if g.opt.CapRule == CapRuleRandom && g.opt.CapRatio == 0 {
		g.opt.CapRatio = .2
	}

	if g.opt.CapRule == CapRuleRandom && (g.opt.CapRatio <= 0 || g.opt.CapRatio >= 1) {
		return errors.New("`CapRatio` must be between 0 and 1 excluded")
	}

	if g.opt.SymbRule == SymbRuleFixed && g.opt.Symbol == 0 {
		g.opt.Symbol = '/'
	}

	if g.opt.SymbRule == SymbRuleRandom && g.opt.Symbol != 0 {
		g.opt.Symbol = 0
	}

	if g.opt.CapRule == "" {
		g.opt.CapRule = CapRuleNone
	}

	if g.opt.SepRule == "" {
		g.opt.SepRule = SepRuleFixed
	}

	if g.opt.SepRule == SepRuleFixed && g.opt.Separator == 0 {
		g.opt.Separator = '-'
	}

	if g.opt.SepRule == SepRuleRandom && g.opt.Separator != 0 {
		g.opt.Separator = 0
	}

	if g.opt.PadRule == PadRuleFixed && g.opt.PadSymbol == 0 {
		g.opt.PadSymbol = '.'
	}

	if g.opt.L33tRatio < 0 || g.opt.L33tRatio > 1 {
		return errors.New("`L33tRatio` must be between 0 and 1 included")
	}

	return nil
}

func (g *Generator) entropy(pass string) float64 {
	charRange := 26
	var usedSymbols string

	if g.opt.CapRule != CapRuleNone {
		charRange *= 2
	}

	if g.opt.DigitsAfter > 0 || g.opt.DigitsBefore > 0 || g.opt.L33tRatio > 0 {
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
