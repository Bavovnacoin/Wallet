package main

func stubNetwCheck(addr string) bool {
	netwAddr := []string{"e", "a", "c", "d", "b"}
	for _, v := range netwAddr {
		if v == addr {
			return true
		}
	}
	return false
}

func checkAddr(addresses []string) int {
	l := 0
	r := len(addresses) - 1
	mid := 0
	for true {
		mid = (r + l) / 2
		checkRes := stubNetwCheck(addresses[mid])

		if checkRes && l <= r {
			return mid
		} else if !checkRes && l <= r {
			return -1
		}

		if checkRes { // TODO: change to network communication
			l = mid + 1
		} else if !checkRes {
			r = mid - 1
		}
	}
	return -1
}

func main() {
	var arr []string = []string{"a", "d"}
	println(checkAddr(arr))
	// var wc wallet_controller.WalletController
	// wc.Launch()
}
