package main

import (
	"fmt"
	"math/big"
)

func main() {
	//println(hashing.HMAC_SHA1("Hello World", "a"))
	a, _ := new(big.Int).SetString("707172737475767778797a7b7c7d7e7f80818283", 16)
	fmt.Printf("%x", a.Bytes())
}
