package account

import (
	"bvcwallet/byteArr"
	"bvcwallet/cryption"
	"bvcwallet/ecdsa"
	"bvcwallet/hashing"
	"bvcwallet/mnemonic"
	"strings"

	"fmt"
	"sort"
)

var CurrAccount Account

type Account struct {
	Id           int
	AccName      string
	HashPass     string          // TODO: byte arr
	MnemonicEncr byteArr.ByteArr // TODO: Check if correct, because it stores it really weirdly
	KeyPairList  []ecdsa.KeyPair
	ArrId        int    `json:"-"`
	Balance      uint64 `json:"-"`
}

// Generates new account and set up a password to encode a private key
func GenAccount(password string, accName string, mnemPhrase byteArr.ByteArr) Account {
	ecdsa.InitValues()
	var newAcc Account

	newAcc.AccName = accName
	newAcc.HashPass = hashing.SHA1(password)

	mnemEncrString := cryption.AES_encrypt(mnemPhrase.ToHexString(), password)
	var mnemEncr byteArr.ByteArr
	mnemEncr.SetFromHexString(mnemEncrString, len(mnemEncrString)/2)
	newAcc.MnemonicEncr = mnemEncr

	var mn mnemonic.Mnemonic
	mnemString := string(mnemPhrase.ByteArr)
	seed := mn.GenSeed(strings.Split(mnemString, " "), "")

	newKeyPair := ecdsa.GenKeyPair(byteArr.ByteArr{ByteArr: seed}, 0)
	newKeyPair.PrivKey = cryption.AES_encrypt(newKeyPair.PrivKey, password)
	newAcc.KeyPairList = append(newAcc.KeyPairList, newKeyPair)

	newAcc.Id = RightBoundAccNum + 1
	RightBoundAccNum++
	return newAcc
}

func AddKeyPairToAccount(password string, allowWrite bool) string {
	if CurrAccount.HashPass == hashing.SHA1(password) {
		ecdsa.InitValues()

		mnemStrDecr := cryption.AES_decrypt(CurrAccount.MnemonicEncr.ToHexString(), password)
		var mnemDecr byteArr.ByteArr
		mnemDecr.SetFromHexString(mnemStrDecr, len(mnemStrDecr)/2)

		var mn mnemonic.Mnemonic
		seed := mn.GenSeed(strings.Split(string(mnemDecr.ByteArr), " "), "")

		newKeyPair := ecdsa.GenKeyPair(byteArr.ByteArr{ByteArr: seed}, len(CurrAccount.KeyPairList))
		newKeyPair.PrivKey = cryption.AES_encrypt(newKeyPair.PrivKey, password)
		CurrAccount.KeyPairList = append(CurrAccount.KeyPairList, newKeyPair)
		Wallet[CurrAccount.ArrId] = CurrAccount

		if allowWrite {
			WriteAccounts()
		}
	} else {
		return "Wrong password!"
	}
	return ""
}

func GetAccUtxo() []UTXO {
	var accUtxo []UTXO
	for i := 0; i < len(CurrAccount.KeyPairList); i++ {
		for j := 0; j < len(CurrAccCoinDatabase); j++ {
			var currAccAddr byteArr.ByteArr
			currAccAddr.SetFromHexString(hashing.SHA1(CurrAccount.KeyPairList[i].PublKey), 20)
			if CurrAccCoinDatabase[j].OutAddress.IsEqual(currAccAddr) {
				accUtxo = append(accUtxo, CurrAccCoinDatabase[j])
			}
		}
	}
	sort.Slice(accUtxo, func(i, j int) bool {
		return accUtxo[i].Value > accUtxo[j].Value
	})
	return accUtxo
}

func GetBalHashOutInd(txHash byteArr.ByteArr, outInd int) uint64 {
	for j := 0; j < len(CurrAccCoinDatabase); j++ {
		if txHash.IsEqual(CurrAccCoinDatabase[j].OutTxHash) && CurrAccCoinDatabase[j].TxOutInd == uint64(outInd) {
			return CurrAccCoinDatabase[j].Value
		}
	}
	return 0
}

func GetBalByAddress(address byteArr.ByteArr) uint64 {
	var Value uint64
	for i := 0; i < len(CurrAccCoinDatabase); i++ {
		if address.IsEqual(CurrAccCoinDatabase[i].OutAddress) {
			Value += CurrAccCoinDatabase[i].Value
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

func InitAccount(accountId int) {
	CurrAccount = Wallet[accountId]
	GetBalance()
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

func GetAccAddresses() []byteArr.ByteArr {
	var addresses []byteArr.ByteArr
	for i := 0; i < len(CurrAccount.KeyPairList); i++ {
		var currAddr byteArr.ByteArr
		currAddr.SetFromHexString(hashing.SHA1(CurrAccount.KeyPairList[i].PublKey), 20)
		addresses = append(addresses, currAddr)
	}
	return addresses
}
