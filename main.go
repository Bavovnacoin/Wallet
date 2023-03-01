package main

import "bvcwallet/wallet_controller"

// type Listener int
// type Reply struct {
// 	Data string
// }

// func (l *Listener) GetLine(line []byte, reply *Reply) error {
// 	rv := string(line)
// 	fmt.Printf("Receive: %v\n", rv)
// 	*reply = Reply{rv}
// 	time.Sleep(10000)
// 	time.Sleep(10 * time.Second)
// 	return nil
// }
// func main() {
// 	addy, err := net.ResolveTCPAddr("tcp", "localhost:8080")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	inbound, err := net.ListenTCP("tcp", addy)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	listener := new(Listener)
// 	rpc.Register(listener)
// 	rpc.Accept(inbound)
// }

func main() {
	// _, err := rpc.Dial("tcp", "localhost:12345")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	var wc wallet_controller.WalletController
	wc.Launch()
}
