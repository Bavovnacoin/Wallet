package mnemonic

import (
	"bufio"
	"os"
)

var wordsFilePath string = "mnemonic/words/"
var Langs []string = []string{"English", "Japanese", "Korean", "Spanish", "Chinese", "French", "Italian", "Czech", "Portuguese"}
var fileNameByLocale map[string]string = map[string]string{Langs[0]: "en.txt", Langs[1]: "ja.txt", Langs[2]: "ko.txt",
	Langs[3]: "es.txt", Langs[4]: "ch.txt", Langs[5]: "fr.txt", Langs[6]: "it.txt", Langs[7]: "cs.txt", Langs[8]: "pt.txt"}

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
