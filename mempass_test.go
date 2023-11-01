package mempass

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
)

func TestDefault(t *testing.T) {
	testPwd(nil, `^[a-z]{6,8}-[a-z]{6,8}-[a-z]{6,8}$`, t)
}

func TestSperatorFixedSet(t *testing.T) {
	testPwd(&Options{
		WordCount: 2,
		SepRule:   SepRuleFixed,
		Separator: '_',
	}, `^[a-z]{6,8}_[a-z]{6,8}$`, t)
}

func TestSperatorRandomDefault(t *testing.T) {
	testPwd(&Options{
		WordCount: 2,
		SepRule:   SepRuleRandom,
	}, `^[a-z]{6,8}[@&!-_^$*%,.;:/=+][a-z]{6,8}$`, t)
}

func TestSperatorRandomSet(t *testing.T) {
	testPwd(&Options{
		WordCount:     2,
		SepRule:       SepRuleRandom,
		SeparatorPool: "@&!",
	}, `^[a-z]{6,8}[@&!][a-z]{6,8}$`, t)
}

func Test2Digits2FixedSymbolDefault(t *testing.T) {
	testPwd(&Options{
		WordCount:     3,
		MinWordLength: 4,
		MaxWordLength: 4,
		SepRule:       SepRuleNone,
		DigitsBefore:  2,
		DigitsAfter:   2,
		SymbolsBefore: 2,
		SymbolsAfter:  2,
		SymbRule:      SymbRuleFixed,
	}, `^(//\d{2}[a-z]{4}\d{2}//){3}$`, t)
}

func Test2Digits2FixedSymbolSet(t *testing.T) {
	testPwd(&Options{
		UseRand:       true,
		WordCount:     2,
		MinWordLength: 4,
		MaxWordLength: 4,
		DigitsBefore:  2,
		DigitsAfter:   2,
		SepRule:       SepRuleNone,
		SymbolsBefore: 2,
		SymbolsAfter:  2,
		SymbRule:      SymbRuleFixed,
		Symbol:        '!',
	}, `^(!!\d{2}[a-z]{4}\d{2}!!){2}$`, t)
}

func Test2SymbolsRandomDefault(t *testing.T) {
	testPwd(&Options{
		WordCount:     2,
		MinWordLength: 4,
		MaxWordLength: 4,
		SymbolsBefore: 2,
		SymbolsAfter:  2,
		SepRule:       SepRuleNone,
	}, `^([@&!-_^$*%,.;:/=+]{2}[a-z]{4}[@&!-_^$*%,.;:/=+]{2}){2}$`, t)
}

func Test2SymbolsRandomSet(t *testing.T) {
	testPwd(&Options{
		WordCount:     2,
		SymbolsBefore: 2,
		SymbolsAfter:  2,
		SymbRule:      SymbRuleRandom,
		SymbolPool:    "@&!",
		SepRule:       SepRuleNone,
	}, `^([@&!]{2}[a-z]{6,8}[@&!]{2}){2}$`, t)
}

func TestPaddingFixed(t *testing.T) {
	testPwd(&Options{
		WordCount: 2,
		PadRule:   PadRuleFixed,
		PadSymbol: '@',
		PadLength: 20,
		SepRule:   SepRuleNone,
	}, `^[a-z]{12,16}@{4,8}$`, t)
}

func TestCapFirst(t *testing.T) {
	testPwd(&Options{
		WordCount:     2,
		SepRule:       SepRuleNone,
		MinWordLength: 6,
		MaxWordLength: 6,
		CapRule:       CapRuleFirstLetter,
	}, `^([A-Z][a-z]{5}){2}$`, t)
}

func TestCapLast(t *testing.T) {
	testPwd(&Options{
		WordCount:     2,
		MinWordLength: 4,
		MaxWordLength: 6,
		CapRule:       CapRuleLastLetter,
		SepRule:       SepRuleNone,
	}, `^([a-z]{3,5}[A-Z]){2}$`, t)
}

func TestCapAllButFirst(t *testing.T) {
	testPwd(&Options{
		WordCount:     2,
		MinWordLength: 4,
		MaxWordLength: 6,
		CapRule:       CapRuleAllButFirstLetter,
		SepRule:       SepRuleNone,
	}, `^([a-z][A-Z]{3,5}){2}$`, t)
}

func TestCapAllButLast(t *testing.T) {
	testPwd(&Options{
		WordCount:     2,
		MinWordLength: 4,
		MaxWordLength: 6,
		CapRule:       CapRuleAllButLastLetter,
		SepRule:       SepRuleNone,
	}, `^([A-Z]{3,5}[a-z]){2}$`, t)
}

func TestCapAlternate(t *testing.T) {
	testPwd(&Options{
		WordCount:     1,
		MinWordLength: 4,
		MaxWordLength: 4,
		CapRule:       CapRuleAlternate,
		SepRule:       SepRuleNone,
	}, `^([A-Z][a-z]){2}$`, t)
}

func TestCapWordAlternate(t *testing.T) {
	testPwd(&Options{
		WordCount: 2,
		CapRule:   CapRuleWordAlternate,
	}, `^[A-Z]{6,8}-[a-z]{6,8}$`, t)
}

func TestCapRandomDefault(t *testing.T) {
	testPwd(&Options{
		WordCount: 2,
		CapRule:   CapRuleRandom,
	}, `^[a-zA-Z]{6,8}-[a-zA-Z]{6,8}$`, t)
}

func TestCapRandomSet(t *testing.T) {
	testPwd(&Options{
		WordCount: 2,
		CapRule:   CapRuleRandom,
		CapRatio:  .8,
	}, `^[a-zA-Z]{6,8}-[a-zA-Z]{6,8}$`, t)
}

func Test1337(t *testing.T) {
	testPwd(&Options{
		WordCount: 2,
		L33tRatio: 0.5,
	}, `^[a-zA-Z0-9]{6,8}-[a-zA-Z0-9]{6,8}$`, t)
}

func testPwd(opt *Options, pattern string, t *testing.T) {
	gen := NewGenerator(opt)
	pwd, _, err := gen.GenPassword()

	fmt.Println("===       Testing password", pwd)

	if err != nil {
		printError(err, t)
	}

	var re = regexp.MustCompile(pattern)
	if !re.Match([]byte(pwd)) {
		printError(errors.New("regex failed"), t)
	}
}

func printError(err error, t *testing.T) {
	t.Errorf("Test failed: %v\n", err)
}
