package networking

import (
	"bvcwallet/byteArr"
	"log"
)

func (c *Connection) IsAddrExist(addr byteArr.ByteArr) bool {
	var repl Reply
	err := c.client.Call("Listener.IsAddrExist", addr.ByteArr, &repl)
	if err != nil {
		log.Println(err)
		return false
	}

	if repl.Data[0] == 0 {
		return false
	}

	return true
}
