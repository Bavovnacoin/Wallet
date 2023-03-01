package networking

import (
	"bufio"
	"log"
	"net/rpc"
	"os"
)

type Connection struct {
	client *rpc.Client
}

type Reply struct {
	Data string
}

func (c *Connection) EstablishConnection() {
	var err error
	c.client, err = rpc.Dial("tcp", "localhost:12345")
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Connection) PingPong() {
	var reply Reply
	err := c.client.Call("Listener.PingPong", []byte("ping"), &reply)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Reply: %v, Data: %v", reply, reply.Data)
}

func abc() {
	client, err := rpc.Dial("tcp", "localhost:12345")
	if err != nil {
		log.Fatal(err)
	}
	in := bufio.NewReader(os.Stdin)
	for {
		line, _, err := in.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		var reply Reply
		err = client.Call("Listener.GetLine", line, &reply)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Reply: %v, Data: %v", reply, reply.Data)
	}
}
