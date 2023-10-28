package mempass

import (
	"errors"
	"math"
	"math/rand"
)

type Options struct {
	UseDict        bool    // Use dictionary. Default false
	WordCount      uint    // Number of words to generate. Using less than 4 is discouraged. Default 4
	Separator      byte    // Separator for words. Default "" (no separator)
	WordLength     uint    // Word length. For dict password, ignored if either `MinWordLength` or `MaxWordLength` is set. For random password, default is 8
	UppercaseRatio float32 // Uppercase ratio. 0.0 = no uppercase, 1.0 = all uppercase, 0.3 = 1/3 uppercase, etc. Default 0
	DigitsCount    uint    // Number of digits to add at the end of each word. Default 0
	MinWordLength  uint    // Minimum word length. O = no minimum. Using less than 4 is discouraged. Default 0
	MaxWordLength  uint    // Maximum word length. O = no maximum. Default 0
}

// Generate a human memorable password
func GenPassword(opt *Options) (string, float64, error) {
	var words [][]byte
	var err error

	if opt.WordCount == 0 {
		opt.WordCount = 4
	}

	if opt.UseDict {
		useMinMax := (opt.MinWordLength > 0 || opt.MaxWordLength > 0) || (opt.WordLength == 0 && (opt.MinWordLength == 0 || opt.MaxWordLength == 0))

		if useMinMax {
			if opt.MaxWordLength == 0 {
				opt.MaxWordLength = 28
			}

			if opt.MinWordLength > 28 || opt.MaxWordLength > 28 {
				return "", 0, errors.New("`MinWordLength` and `MaxWordLength` cannot be greater than 28")
			}

			if opt.MinWordLength > 0 && opt.MaxWordLength > 0 && opt.MinWordLength > opt.MaxWordLength {
				return "", 0, errors.New("`MinWordLength` cannot be greater than `MaxWordLength`")
			}
		} else {
			if opt.WordLength > 28 || opt.WordLength == 24 || opt.WordLength == 26 || opt.WordLength == 27 {
				return "", 0, errors.New("`WordLength` must be beetween 1 and 28 and not 24, 26 or 27")
			}

			if opt.WordLength == 0 {
				opt.WordLength = 8
			}
		}

		if words, err = getDictWords(opt, useMinMax); err != nil {
			return "", 0, err
		}
	} else {
		if opt.WordLength == 0 {
			opt.WordLength = 8
		}

		words = genRandPwd(opt)
	}

	words = capAndAddNum(words, opt.UppercaseRatio, opt.DigitsCount)
	pwd := ""

	for i, word := range words {
		pwd += string(word)

		if i < len(words)-1 {
			pwd += string(opt.Separator)
		}
	}

	return pwd, entropy(pwd, opt), nil
}

func capAndAddNum(words [][]byte, up float32, nc uint) [][]byte {
	newWords := make([][]byte, len(words))

	for i, word := range words {
		newWord := make([]byte, len(word))
		if up > 0 {
			for i := 0; i < len(word); i++ {
				if rand.Float32() <= up {
					newWord[i] = word[i] - ('a' - 'A')
				} else {
					newWord[i] = word[i]
				}
			}
		} else {
			newWord = word
		}

		addNums(&newWord, nc)
		newWords[i] = newWord
	}

	return newWords
}

func addNums(str *[]byte, nc uint) {
	numbers := "0123456789"

	for i := 0; i < int(nc); i++ {
		idx := rand.Intn(10)
		*str = append(*str, numbers[idx])
	}
}

func entropy(pass string, opt *Options) float64 {
	charRange := 26.0

	if opt.UppercaseRatio > 0 {
		charRange *= 2
	}

	// TODO: use a character set for symbols
	if opt.Separator != 0 {
		charRange += 1
	}

	len := float64(len(pass))

	return math.Log2(float64(math.Pow(charRange, len)))
}
