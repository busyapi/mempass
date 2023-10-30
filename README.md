# Human Memorable Password Generator

## Introduction

This module generates strong, memorable passwords.

Source of password can be either a dictionary of English words or derived from statistics of threee-letter sequences in English. The latter is a direct port in Go of https://www.multicians.org/thvv/gpw.js

This modules also allows you add digits and symbols to passwords. While this makes passwords even stronger, it decreases password readability so you probably will have to choose balanced settings. Idea for this come from https://xkpasswd.net/s/.

The password generator also returns the password generation [entropy](#entropy)

## Installation

```sh
go get github.com/busyapi/mempass
```

## Usgae

### Example

```go
gen := mempass.NewGenerator(&mempass.Options{
		UseDict:       true,
		WordCount:     4,
		MinWordLength: 6,
		MaxWordLength: 8,
		CapRule:       mempass.CapRuleFirstLetter,
		SepRule:       mempass.SepRuleFixed,
		Separator: '!',
	})

	password, entropy, err := gen.GenPassword()
```

This will produce a password liek `Auroral!Tallied!Couture!Crewmen`

### Options

```go
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
```

<a id="entropy"></a>

## Entropy

Entropy (measured in bits) is an indicator of the strenght of the password **generation method**, not of the password itself. That's why you will always get the same entropy for different passwords if the generation rules remain the same.

Basically, the longer your password is and the wider your characters pool is, the higher the entropy is.

It's generally considered an entrepy above 120 bits provide a very strong generation strengh.

## TODO

- More options?
- Provide more password strengh indicator
- Tests
