package networking

import (
	"log"
)

func (c *Connection) IsAddrExist(address string) bool {
	byteArr, isConv := c.ToByteArr(address)
	if !isConv {
		return false
	}

	var repl Reply
	err := c.client.Call("Listener.IsAddrExist", byteArr, &repl)
	if err != nil {
		log.Println(err)
		return false
	}

	if repl.Data[0] == 0 {
		return false
	}

	return true
}
