package account

/*
	The wallet data (list of Account structure - wallet.json) is stored
	locally on user's device. There could be a multiple ammount of accounts
	in one wallet (many users can use one device, so there should be a way to
	distinguish data). Private key is encrypted by user's own password
	using AES algorithm. The password is stored in wallet.json as a hash
	value that is created using SHA-1.
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	WalletName = "wallet.json"
)

var Wallet []Account

func IsWalletExists() bool {
	file, err := os.Open(WalletName)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	file.Close()
	return true
}

func WriteAccounts() {
	byteData, _ := json.MarshalIndent(Wallet, "", "    ")
	os.WriteFile(WalletName, byteData, 0777)
}

func GetWalletData() {
	jsonFile, err := os.Open(WalletName)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &Wallet)
	jsonFile.Close()

}
