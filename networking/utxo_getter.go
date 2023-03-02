package networking

import (
	"bvcwallet/account"
	"bvcwallet/byteArr"
	"log"
)

func (c *Connection) GetMyUtxo(addresses []byteArr.ByteArr) bool {
	byteArr, isConv := c.ToByteArr(addresses)
	if !isConv {
		return false
	}

	var repl Reply
	err := c.client.Call("Listener.GetUtxoByAddr", byteArr, &repl)
	if err != nil {
		log.Println(err)
		return false
	}

	c.FromByteArr(repl.Data, &account.CurrAccCoinDatabase) //TODO: problem here (wrong byte array -> struct array conversion)!
	for i := 0; i < len(account.CurrAccCoinDatabase); i++ {
		println(account.CurrAccCoinDatabase[i].Value)
	}
	return true

}
