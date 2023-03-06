package main

import "bvcwallet/wallet_controller"

func main() {
	// TODO: create transaction
	// ecdsa.InitValues()
	// account.CurrAccount.KeyPairList = append(account.CurrAccount.KeyPairList, ecdsa.KeyPair{PrivKey: "cf5a1a10730695d9568b49bf059ecaa06a230e5222f74f2c5b9a212e02b3ad7c042337c44d1482f9f3465a978b3a03fb567d71d5960ff1168570f4767cf3fae1", PublKey: "02b0c0375c31c42a62111518a2cdf4cec059c73264da6bff35e68bdbb5250deb53"})

	// var utxo account.UTXO
	// utxo.BlockHeight = 0
	// utxo.OutAddress.SetFromHexString("92a405420017dda1ca887c3080f0b437048241bb", 20)
	// utxo.OutTxHash.SetFromHexString("123456789a123456789a123456789a123456789a", 20)
	// utxo.TxOutInd = 0
	// utxo.Value = 7000000000
	// account.CurrAccCoinDatabase = append(account.CurrAccCoinDatabase, utxo)

	// var outAddr []byteArr.ByteArr
	// outAddr = append(outAddr, byteArr.ByteArr{})
	// outAddr[0].SetFromHexString("abcabcabcaabcabcabcaabcabcabcaabcabcabca", 20)

	// var outVals []uint64
	// outVals = append(outVals, 20000)
	// tx, mes := transaction.CreateTransaction("Password_1", outAddr, outVals, 2, 0)
	// println(mes)
	// transaction.PrintTransaction(tx)
	var wc wallet_controller.WalletController
	wc.Launch()
}
