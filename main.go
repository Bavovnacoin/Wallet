package main

func binarySearch(needle int, haystack []int) bool {

	low := 0
	high := len(haystack) - 1

	for low <= high {
		median := (low + high) / 2

		if haystack[median] < needle {
			low = median + 1
		} else {
			high = median - 1
		}
	}

	if low == len(haystack) || haystack[low] != needle {
		return false
	}

	return true
}

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

		if checkRes && l >= r {
			return mid
		} else if !checkRes && l >= r {
			return mid - 1
		}

		if checkRes {
			l = mid + 1
		} else if !checkRes {
			r = mid - 1
		}
	}
	return -1
}

func foo() string {
	defer println("ayayyaya")
	for i := 0; i < 100; i++ {
		if i == 10 {
			return "a"
		}
	}
	return "b"
}

func main() {
	foo()
	// var wc wallet_controller.WalletController
	// wc.Launch()
}
