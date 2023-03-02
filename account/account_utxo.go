package account

import "bvcwallet/byteArr"

type UTXO struct {
	OutTxHash   byteArr.ByteArr
	TxOutInd    uint64
	Value       uint64
	OutAddress  byteArr.ByteArr
	BlockHeight uint64
}

var CurrAccCoinDatabase []UTXO
