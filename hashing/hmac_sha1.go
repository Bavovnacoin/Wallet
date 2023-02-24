package hashing

import (
	"encoding/hex"
)

const (
	block_size  = 64
	result_size = 20

	ipad = 0x36
	opad = 0x5c
)

func byteArrToStr(arr []byte) string {
	res := ""
	for i := 0; i < len(arr); i++ {
		res += string(arr[i])
	}
	return res
}

func HMAC_SHA1(key, message []byte) string {
	if len(key) > block_size {
		key = []byte(SHA1(byteArrToStr(key)))
	}

	if len(key) < block_size {
		key = append(key, make([]byte, block_size-len(key))...)
	}

	var ikeypad []byte
	var okeypad []byte
	for i := 0; i < len(key); i++ {
		ikeypad = append(ikeypad, key[i]^ipad)
		okeypad = append(okeypad, key[i]^opad)
	}

	b, _ := hex.DecodeString(SHA1(byteArrToStr(append(ikeypad, message...))))
	return SHA1(byteArrToStr(append(okeypad, b...)))
}
