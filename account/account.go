package account

import (
	"bvcwallet/byteArr"
	"bvcwallet/cryption"
	"bvcwallet/ecdsa"
	"bvcwallet/hashing"

	"bvcwallet/networking"
	"fmt"
	"sort"
)

var CurrAccount Account

type Account struct {
	Id          int
	AccName     string
	HashPass    string          // TODO: byte arr
	SeedEncr    byteArr.ByteArr // TODO: Check if correct, because it stores it really weirdly
	KeyPairList []ecdsa.KeyPair
	ArrId       int    `json:"-"`
	Balance     uint64 `json:"-"`
}

// Generates new account and set up a password to encode a private key
func GenAccount(password string, accName string, seed string) Account {
	ecdsa.InitValues()
	var newAcc Account

	newAcc.AccName = accName
	newAcc.HashPass = hashing.SHA1(password)

	newKeyPair := ecdsa.GenKeyPair()
	newKeyPair.PrivKey = cryption.AES_encrypt(newKeyPair.PrivKey, password)

	newAcc.Id = RightBoundAccNum + 1
	RightBoundAccNum++
	newAcc.KeyPairList = append(newAcc.KeyPairList, newKeyPair)

	seedEncrString := cryption.AES_encrypt(seed, password)
	var seedEncr byteArr.ByteArr
	seedEncr.SetFromHexString(seedEncrString, len(cryption.AES_encrypt(seed, password))/2) //
	newAcc.SeedEncr = seedEncr

	return newAcc
}

func AddKeyPairToAccount(password string) string {
	if CurrAccount.HashPass == hashing.SHA1(password) {
		ecdsa.InitValues()
		newKeyPair := ecdsa.GenKeyPair()
		newKeyPair.PrivKey = cryption.AES_encrypt(newKeyPair.PrivKey, password)
		CurrAccount.KeyPairList = append(CurrAccount.KeyPairList, newKeyPair)
		Wallet[CurrAccount.ArrId] = CurrAccount
	} else {
		return "Wrong password!"
	}
	return ""
}

func GetAccUtxo() []networking.TXO {
	var accUtxo []networking.TXO
	for i := 0; i < len(CurrAccount.KeyPairList); i++ {
		for j := 0; j < len(networking.CoinDatabase); j++ {
			var currAccAddr byteArr.ByteArr
			currAccAddr.SetFromHexString(hashing.SHA1(CurrAccount.KeyPairList[i].PublKey), 20)
			if networking.CoinDatabase[j].OutAddress.IsEqual(currAccAddr) {
				accUtxo = append(accUtxo, networking.CoinDatabase[j])
			}
		}
	}
	sort.Slice(accUtxo, func(i, j int) bool {
		return accUtxo[i].Value > accUtxo[j].Value
	})
	return accUtxo
}

func GetBalHashOutInd(txHash byteArr.ByteArr, outInd int) uint64 {
	for j := 0; j < len(networking.CoinDatabase); j++ {
		if txHash.IsEqual(networking.CoinDatabase[j].OutTxHash) && networking.CoinDatabase[j].TxOutInd == uint64(outInd) {
			return networking.CoinDatabase[j].Value
		}
	}
	return 0
}

func GetBalByAddress(address byteArr.ByteArr) uint64 {
	var Value uint64
	for i := 0; i < len(networking.CoinDatabase); i++ {
		if address.IsEqual(networking.CoinDatabase[i].OutAddress) {
			Value += networking.CoinDatabase[i].Value
		}
	}
	return Value
}

// A function counts all the UTXOs that is on specific public keys on user's account
func GetBalance() uint64 {
	CurrAccount.Balance = 0
	for i := 0; i < len(CurrAccount.KeyPairList); i++ {
		var address byteArr.ByteArr
		address.SetFromHexString(hashing.SHA1(CurrAccount.KeyPairList[i].PublKey), 20)
		CurrAccount.Balance += GetBalByAddress(address)
	}
	return CurrAccount.Balance
}

func PrintBalance() {
	GetBalance()
	var bal float64 = float64(CurrAccount.Balance) / 100000000.
	fmt.Printf("Balance: %.8f BVC\n", bal)
}

func getAccountInd(accountId int) int {
	for i := 0; i < len(Wallet); i++ {
		if Wallet[i].Id == accountId {
			Wallet[i].ArrId = i
			return i
		}
	}
	return -1
}

func InitAccount(accountId int) bool {
	ecdsa.InitValues()
	accInd := getAccountInd(accountId)
	if accInd != -1 {
		CurrAccount = Wallet[accInd]
		return true
	}
	return false
}

func SignData(hashMes string, kpInd int, pass string) (string, bool) {
	if CurrAccount.HashPass != hashing.SHA1(pass) {
		return "", true
	}
	kp := CurrAccount.KeyPairList[kpInd]
	kp.PrivKey = cryption.AES_decrypt(kp.PrivKey, pass)

	return ecdsa.Sign(hashMes, kp.PrivKey), false
}

func VerifData(hashMes string, kpInd int, signature string) bool {
	kp := CurrAccount.KeyPairList[kpInd]
	return ecdsa.Verify(kp.PublKey, signature, hashMes)
}
