package hashing

import (
	"fmt"
	"strconv"
	"strings"
)

func binary(s string) string {
	res := ""
	for _, c := range s {
		res = fmt.Sprintf("%s%.8b", res, c)
	}
	return res
}

func addZerosTo(mes string, to int) string {
	return strings.Repeat("0", to-len(mes)) + mes
}

func filler(binMes string) string {
	length := len(binMes)
	binLen := strconv.FormatInt(int64(length), 2)
	return binMes + "1" + strings.Repeat("0", 512-((length+65)%512)) + addZerosTo(binLen, 64)
}

func leftShift(val uint32, k int) uint32 {
	return uint32((uint64(val) << k) % 4294967295)
}

func f1(arr []uint32, f *uint32, k *uint32) {
	*f = (arr[1] & arr[2]) | ((^arr[1]) & arr[3])
	*k = 0x5A827999
}

func f2(arr []uint32, f *uint32, k *uint32) {
	*f = arr[1] ^ arr[2] ^ arr[3]
	*k = 0x6ED9EBA1
}

func f3(arr []uint32, f *uint32, k *uint32) {
	*f = (arr[1] & arr[2]) | (arr[1] & arr[3]) | (arr[2] & arr[3])
	*k = 0x8F1BBCDC
}

func f4(arr []uint32, f *uint32, k *uint32) {
	*f = arr[1] ^ arr[2] ^ arr[3]
	*k = 0xCA62C1D6
}

func SHA1(message string) string {
	var arr []uint32
	var f uint32
	var k uint32
	var functions []func([]uint32, *uint32, *uint32)
	var sha1_h = []uint32{0x67452301, 0xEFCDAB89, 0x98BADCFE, 0x10325476, 0xC3D2E1F0}

	functions = []func([]uint32, *uint32, *uint32){f1, f2, f3, f4}
	arr = make([]uint32, len(sha1_h))
	f = 0
	k = 0

	binForm := binary(message)
	binForm = filler(binForm)
	for i := 0; i < len(binForm)/512; i += 1 {
		var w []uint32
		chunk := binForm[i*512 : i*512+512]

		for j := 0; j < 16; j += 1 {
			num64, _ := strconv.ParseUint(chunk[j*32:j*32+32], 2, 32)
			w = append(w, uint32(num64))
		}

		for j := 16; j < 80; j += 1 {
			w = append(w, leftShift(w[j-3]^w[j-8]^w[j-14]^w[j-16], 1))
		}

		copy(arr, sha1_h)

		f = 0
		k = 0
		for j := 0; j < 80; j += 1 {
			functions[(j-(j/19))/19](arr, &f, &k)
			temp := leftShift(arr[0], 5) + f + arr[4] + k + w[j]
			arr[4] = arr[3]
			arr[3] = arr[2]
			arr[2] = leftShift(arr[1], 30)
			arr[1] = arr[0]
			arr[0] = temp
		}

		for i := 0; i < 5; i++ {
			sha1_h[i] += arr[i]
		}
	}
	a := addZerosTo(fmt.Sprintf("%x", sha1_h[0]), 8)
	b := addZerosTo(fmt.Sprintf("%x", sha1_h[1]), 8)
	c := addZerosTo(fmt.Sprintf("%x", sha1_h[2]), 8)
	d := addZerosTo(fmt.Sprintf("%x", sha1_h[3]), 8)
	e := addZerosTo(fmt.Sprintf("%x", sha1_h[4]), 8)
	return a + b + c + d + e
}
