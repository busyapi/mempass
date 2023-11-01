# Human Memorable Password Generator

## Introduction and feeatures

This module generates strong passwords based on memorabled words. Current features are:

- Number of words
- Words length can be from 1 to 28
- Add separators beetween words
- Mulitple letter capitalization options
- Add digits before/after each word
- Add symbols before/after each word
- Add 1337 encoding for letters a, e, i, o, s, t
- Choice between dictionary of English words or randomly generated memorable words.
- Calculate the password generation [entropy](#entropy)

This modules is inspired by the great work of:

- https://www.multicians.org/thvv/gpw.js (the random memorable password generator is a direct port in Go)
- https://xkpasswd.net/s/

## Installation

```sh
go get github.com/busyapi/mempass
```

## Usgae

### Example

```go
gen := mempass.NewGenerator(nil)
password, entropy, err := gen.GenPassword()
```

This will produce a password like `tildes-brazen-quezals`

### Options and default values

```go
type Options struct {
	UseRand          bool     // Use randomly generated words instead of dictionary words . Default false
	WordCount        uint     // Number of words to generate. Using less than 2 is discouraged. Default is 3
	MinWordLength    uint     // Minimum word length. O = no minimum. Using less than 4 is discouraged. Default is 6
	MaxWordLength    uint     // Maximum word length. O = no maximum. Default is 8
	DigitsAfter      uint     // Number of digits to add at the end of each word. Default is 0
	DigitsBefore     uint     // Number of digits to add at the begining of each word. Default is 0
	CapRule          CapRule  // Capitalization rule
	CapRatio         float32  // Uppercase ratio. 0.0 = no uppercase, 1.0 = all uppercase, 0.3 = 1/3 uppercase, etc. Only used if `CapRule` is `CapRandom`. Default is 0.2
	SymbRule         SymbRule // Rule for adding symbols. Default is `SymbRuleNone`
	SymbolsAfter     uint     // Number of symbols to add at the end of each word. Default is 0
	SymbolsBefore    uint     // Number of symbols to add at the begining of each word. Default is 0
	SymbolPool       string   // Symbols pool. Only used if `SymbRule` is `SymbRuleRandom`. Default is "@&!-_^$*%,.;:/=+"
	Symbol           byte     // Symbol character. Only used if `SymbRule` is `SymbRuleFixed`. Default is `/`
	SepRule          SepRule  // Seperator type. Default is `SepRuleFixed`
	SeparatorPool    string   // Seperators pool. Only used if `SepRule` is `SepRuleRandom`. Default is "@&!-_^$*%,.;:/=+"
	Separator        byte     // Separator for words. Only used if `SepRule` is `SepRuleFixed`. Default is '-'
	PadRule          PadRule  // Padding rule. Ignored if `PadLength` is 0
	PadSymbol        byte     // Padding symbol. Only used if `PadRule` si `PadRuleFixed`. Default is `.`
	PadLength        uint     // Password length to reach with padding.
	L33tRatio        float32  // 1337 coding ratio. 0.0 = no 1337, 1.0 = all 1337, 0.3 = 1/3 1337, etc`. Default is 0
	CalculateEntropy bool     // Calculate entropy. Default is false
}
```

<a id="entropy"></a>

## Entropy

Entropy (measured in bits) is an indicator of the strength of the password **generation method**, not of the password itself. That's why you will always get the same entropy for different passwords if the generation rules remain the same.

Basically, the longer your password is and the wider your characters pool is, the higher the entropy is.

It's generally considered that an entropy above 120 bits provide a very strong generation strength.

## TODO

- More options?
- Provide more password strengh indicator
