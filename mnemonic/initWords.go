package mnemonic

import (
	"bufio"
	"os"
)

var wordsFilePath string = "mnemonic/words/"
var fileNameByLocale map[string]string = map[string]string{"en": "en.txt", "ja": "ja.txt", "ko": "ko.txt",
	"es": "es.txt", "ch": "ch.txt", "fr": "fr.txt", "it": "it.txt", "cs": "cs.txt", "pt": "pt.txt"}

func (mn *Mnemonic) InitWords() bool {
	fname := wordsFilePath + fileNameByLocale[mn.locale]
	if fname == wordsFilePath {
		return false
	}
	file, ferr := os.Open(fname)

	if ferr != nil {
		return false
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		mn.words = append(mn.words, scanner.Text())
	}

	return true
}
