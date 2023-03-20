package mnemonic

import (
	"bvcwallet/hashing"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	mnemLenBits = 256
	iterC       = 512
	seedLen     = 64
)

type Mnemonic struct {
	source rand.Source
	random *rand.Rand

	entropyLowBound *big.Int
	entropyHghBound *big.Int
	entropy         *big.Int

	words  []string
	locale string
}

func (mn *Mnemonic) InitMnemonic(locale string) {
	mn.source = rand.NewSource(time.Now().Unix())
	mn.random = rand.New(mn.source)
	mn.entropyLowBound, _ = new(big.Int).SetString("1"+strings.Repeat("0", mnemLenBits-1), 2)
	mn.entropyHghBound, _ = new(big.Int).SetString(strings.Repeat("1", mnemLenBits), 2)
	mn.locale = locale
}

func (mn *Mnemonic) GenEntropy() *big.Int {
	boundSub := new(big.Int).Sub(mn.entropyHghBound, mn.entropyLowBound)
	randEntrDiff := new(big.Int).Rand(mn.random, boundSub)
	return new(big.Int).Add(mn.entropyLowBound, randEntrDiff)
}

func (mn *Mnemonic) GenMnemonicPhrase() []string {
	entropyBits := fmt.Sprintf("%b", mn.GenEntropy())

	hashEntropy, _ := new(big.Int).SetString(hashing.SHA1(mn.entropy.String()), 16)
	hashEntropyBits := fmt.Sprintf("%b", hashEntropy)
	entropyBits += hashEntropyBits[:len(entropyBits)/32]

	if !mn.InitWords() {
		println("En error occured when initializing mnemonic phrase words.")
		return nil
	}

	var mnemonicPhrase []string
	for i := 11; i < len(entropyBits); i += 11 {
		wordInd, _ := strconv.ParseInt(entropyBits[i-11:i], 2, 64)
		mnemonicPhrase = append(mnemonicPhrase, mn.words[wordInd])
	}

	return mnemonicPhrase
}

func (mn *Mnemonic) GenSeed(mnemonic []string, password string) []byte {
	return hashing.PBKDF2(strings.Join(mnemonic, " "), "mnemonic"+password, iterC, seedLen)
}
