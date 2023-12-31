package mempass

import (
	"bufio"
	"embed"
	"errors"
	"math/rand"
)

//go:embed wordsEn.txt
var embeddedFile embed.FS

// Get random words from the dictionary file
func getDictWords(opt *Options) ([][]rune, error) {
	var words [][]rune
	dict, err := readDictFile(opt)
	if err != nil {
		return nil, err
	}

	keys := make([]int, 0, len(dict))

	for k := range dict {
		keys = append(keys, k)
	}

	for i := 0; i < int(opt.WordCount); i++ {
		var dictIdx int

		keyIdx := int(rand.Float32() * float32((len(keys) - 1)))
		dictIdx = keys[keyIdx]

		if _, exists := dict[dictIdx]; !exists {
			i--
			continue
		}

		words = append(words, dict[dictIdx][rand.Intn(len(dict[dictIdx]))])
	}

	return words, nil
}

// Read words from the dictionary file and store them in a map
// The keys of the map are the line lengths
func readDictFile(opt *Options) (map[int][][]rune, error) {
	words := make(map[int][][]rune)

	// sourceFilePath, err := getCurrentSourceFilePath()
	// if err != nil {
	// 	return nil, errors.New("Error reading dict file: " + err.Error())
	// }
	//
	// sourceFileDir := filepath.Dir(sourceFilePath)

	//file, err := os.Open(sourceFileDir + "/wordsEn.txt")
	file, err := embeddedFile.Open("wordsEn.txt")
	if err != nil {
		return nil, errors.New("Error reading dict file: " + err.Error())
	}

	defer file.Close()

	// Create a scanner to read lines from the file
	scanner := bufio.NewScanner(file)

	// Read lines and append them to the slice
	for scanner.Scan() {
		line := scanner.Text()
		runes := toRunes(line)
		lc := len(runes)

		// Don't include words that are bellow `minWl` or above `maxWl`
		if lc < int(opt.MinWordLength) || (opt.MaxWordLength > 0 && lc > int(opt.MaxWordLength)) {
			continue
		}

		// Create a copy of the line and append it to the slice
		lineCopy := make([]rune, len(line))
		copy(lineCopy, runes)
		words[lc] = append(words[lc], lineCopy)
	}

	// Check for any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		return nil, errors.New("Error while scanning dict file: " + err.Error())
	}

	return words, nil
}
