package hashing

import (
	"math/big"
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

func HMAC_SHA1(key, message string) string {
	keyByte := []byte(key)

	if len(keyByte) > block_size {
		keyByte = []byte(SHA1(key))
	}

	if len(keyByte) < block_size {
		keyByte = append(keyByte, make([]byte, block_size-len(keyByte))...)
	}

	var ikeypad []byte
	var okeypad []byte
	for i := 0; i < len(keyByte); i++ {
		ikeypad = append(ikeypad, keyByte[i]^ipad)
		okeypad = append(okeypad, keyByte[i]^opad)
	}

	ikeypadBig, _ := new(big.Int).SetString(SHA1(string(ikeypad)+message), 16)
	return SHA1(string(okeypad) + byteArrToStr(ikeypadBig.Bytes()))
}
