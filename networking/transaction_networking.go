package networking

import (
	"bvcwallet/transaction"
	"log"
)

func (c *Connection) SendTransaction(tx transaction.Transaction, isAccepted *bool) bool {
	byteArr, isConv := c.ToByteArr(tx)
	if !isConv {
		return false
	}

	var repl Reply
	err := c.client.Call("Listener.AddNewTxToMemp", byteArr, &repl)
	if err != nil {
		log.Println(err)
		return false
	}

	if repl.Data[0] == 1 {
		*isAccepted = true
	} else {
		*isAccepted = false
	}
	return true
}
