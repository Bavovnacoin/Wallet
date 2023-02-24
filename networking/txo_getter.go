package networking

import "bvcwallet/byteArr"

type TXO struct {
	OutTxHash   byteArr.ByteArr
	TxOutInd    uint64
	Value       uint64
	OutAddress  byteArr.ByteArr
	BlockHeight uint64
}

var CoinDatabase []TXO
