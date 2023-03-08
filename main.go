package main

import "bvcwallet/wallet_controller"

func main() {
	// seedEncrString := cryption.AES_encrypt("96fbc728b291049dff549a95bd3c821881ee38dce6a5cb34ece64df879908eccad61d97321b7e5e06a4411e3889d21bd62b1c5d07dce21a56c5c7ca651c28fec", "Password_1")
	// var seedEncr byteArr.ByteArr
	// seedEncr.SetFromHexString(seedEncrString, len(seedEncrString)/2)
	// //println(seedEncr.ToHexString())

	// decrString := cryption.AES_decrypt(seedEncr.ToHexString(), "Password_1")
	// var seed byteArr.ByteArr
	// seed.SetFromHexString(decrString, len(decrString)/2)

	var wc wallet_controller.WalletController
	wc.Launch()
}
