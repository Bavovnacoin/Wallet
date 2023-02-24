package hashing

import (
	"encoding/hex"
)

func intToByteArr(blockInd int) []byte {
	var blockIndByteArr []byte = make([]byte, 4)
	blockIndByteArr[3] = byte(blockInd)

	for i := 1; i <= 3; i++ {
		blockIndByteArr[i-1] = byte(blockInd >> 24 / i)
	}
	return blockIndByteArr
}

func f(password string, salt string, iterCount int, blockInd int) []byte {
	mes := append([]byte(salt), intToByteArr(blockInd)...)
	U := HMAC_SHA1([]byte(password), mes)
	UbyteArr, _ := hex.DecodeString(U)
	resByte, _ := hex.DecodeString(U)

	for i := 2; i <= iterCount; i++ {
		U = HMAC_SHA1([]byte(password), UbyteArr)

		UbyteArr, _ = hex.DecodeString(U)
		for i := 0; i < len(UbyteArr); i++ {
			resByte[i] ^= UbyteArr[i]
		}
	}

	return resByte
}

func PBKDF2(password string, salt string, iterCount int, keyLen int) []byte {
	var blockCount int = int((keyLen + result_size - 1)) / result_size
	r := keyLen - (blockCount-1)*result_size
	var T []byte

	for i := 1; i <= blockCount; i++ {
		T = append(T, f(password, salt, iterCount, i)...)
	}

	return T[:(blockCount-1)*result_size+r]
}
