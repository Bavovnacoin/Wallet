package wallet_controller

import (
	"bvcwallet/account"
	"bvcwallet/byteArr"
	"bvcwallet/cryption"
	"bvcwallet/ecdsa"
	"bvcwallet/hashing"
	"bvcwallet/mnemonic"
	"bvcwallet/networking"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	specSign = " !”#$%&’()*+,./:;<=>?@[/^_`{}|~"
)

var allowExit bool = false

func isAccNameValid(name string) bool {
	if len(name) < 4 || len(name) > 25 {
		return false
	}
	return true
}

func isAccPassValid(pass string) bool {
	// Regex for checking at least one lower, one upper, one num and spec sign
	reg := fmt.Sprintf("^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[%s])[A-Za-z\\d%s]$", specSign, specSign)
	isPassMatch, _ := regexp.MatchString(reg, pass)
	if (len(pass) < 8 || len(pass) > 20) && !isPassMatch {
		return false
	}
	return true
}

func (wc *WalletController) fieldEnterer(enterText string, errorText string, validation func(val string) bool) string {
	println("Type \"back\" to back to the menu.")
	println(enterText)
	enterValue := ""
	validErrText := ""
	for true {
		if validErrText != "" {
			println(validErrText)
			validErrText = ""
		}

		wc.scann.Scan()
		enterValue = wc.scann.Text()

		if enterValue == "back" {
			allowExit = true
			return ""
		}

		if validation(enterValue) {
			break
		} else {
			validErrText = errorText
		}
	}
	return enterValue
}

func byteArrToString(arr []byte) string {
	res := ""
	for i := 0; i < len(arr); i++ {
		res += string(arr[i])
	}
	return res
}

func (wc *WalletController) CreateAccount() bool {
	wc.ClearConsole()
	accountName := wc.fieldEnterer("Enter an account name", "Account name length must be in range from 4 to 25 symbols.", isAccNameValid)
	if allowExit {
		allowExit = false
		return false
	}

	wc.ClearConsole()
	accountPass := wc.fieldEnterer("Enter an account password", fmt.Sprintf("Account password length must be in range from 8 to 20 symbols. It should contain upper and lower case letters, numbers and special signs (%s).\n", specSign), isAccPassValid)
	if allowExit {
		return false
	}

	var mn mnemonic.Mnemonic
	mn.InitMnemonic("en")
	phrase := mn.GenMnemonicPhrase()

	command := ""
	for true {
		wc.ClearConsole()
		println("You need to save a mnemonic phrase to restore your keys in the future.")
		println("To show it type in \"show\" command.")
		fmt.Scan(&command)
		if command == "show" {
			break
		}
	}
	println()
	println(strings.Join(phrase, " "))

	seed := byteArr.ByteArr{ByteArr: mn.GenSeed(phrase, "")}
	account.Wallet = append(account.Wallet, account.GenAccount(accountPass, accountName, seed))
	account.CurrAccount = account.Wallet[len(account.Wallet)-1]
	account.WriteAccounts()
	println("Account successfully created.")
	println("To continue enter any word.")
	fmt.Scan(&command)
	return true
}

func (wc *WalletController) validateUser(accId int) bool {
	isPassErr := false
	var password string
	for true {
		wc.ClearConsole()
		if isPassErr {
			println("Wrong password. Try again.")
			isPassErr = false
		}

		println("Enter password or \"back\" to get back.")
		wc.scann.Scan()
		password = wc.scann.Text()
		hashPass := hashing.SHA1(password)

		if hashPass == account.Wallet[accId].HashPass {
			return true
		} else if password == "back" {
			break
		} else if hashPass != account.Wallet[accId].HashPass {
			isPassErr = true
		}
	}
	return false
}

func (wc *WalletController) initAccount() bool {
	printErr := false
	account.GetWalletData()
	for true {
		wc.ClearConsole()
		if printErr {
			println("You have typed in wrong value. Try again.")
			printErr = false
		}
		println("Choose account and type in the right number or type \"exit\" to exit:")

		for i := 0; i < len(account.Wallet); i++ {
			fmt.Printf("%d. %s\n", i, account.Wallet[i].AccName)
		}
		fmt.Printf("%d. Create a new one\n", len(account.Wallet))
		fmt.Printf("%d. Add via a mnemonic phrase\n", len(account.Wallet)+1)

		var command string
		wc.scann.Scan()
		command = wc.scann.Text()

		accIdNum, err := strconv.Atoi(command)
		if command == "exit" {
			wc.walletLaunched = false
			return false
		} else if err != nil || accIdNum < 0 || accIdNum > len(account.Wallet)+1 {
			printErr = true
		} else if err != nil || accIdNum < 0 || accIdNum <= len(account.Wallet)+1 {
			if accIdNum == len(account.Wallet) {
				wc.CreateAccount()
			} else if accIdNum == len(account.Wallet)+1 {
				wc.EnterMnemonic()
			} else {
				validResult := wc.validateUser(accIdNum)
				if validResult {
					account.InitAccount(accIdNum)
					return true
				}
			}
		}
	}
	return false
}

func isMnemPhraseValid(phrase string) bool {
	phraseArr := strings.Split(phrase, " ")
	if len(phraseArr) == 0 {
		return false
	}
	return true
}

// Binary search for checking existing keys
func checkExistingKeyPairs(kps []ecdsa.KeyPair) int {
	var con networking.Connection
	con.Establish()
	defer con.Close()

	for i := len(kps) - 1; i >= 0; i-- {
		var addr byteArr.ByteArr
		addr.SetFromHexString(hashing.SHA1(kps[i].PublKey), 20)
		if con.IsAddrExist(addr) {
			return i
		}
	}
	return -1
}

// Checking an existing keys in the network
func getExistingKeyPairs(seed byteArr.ByteArr, keysCheckAmmount int, password string) []ecdsa.KeyPair {
	var existingKeys []ecdsa.KeyPair

	kpInd := 0
	for true {
		var currKPs []ecdsa.KeyPair
		for ; kpInd < keysCheckAmmount; kpInd += 1 {
			newKp := ecdsa.GenKeyPair(seed, kpInd)
			newKp.PrivKey = cryption.AES_encrypt(newKp.PrivKey, password)
			currKPs = append(currKPs, newKp)
		}

		currKPsInd := checkExistingKeyPairs(currKPs)
		if currKPsInd < 0 {
			break
		} else if currKPsInd == len(currKPs)-1 {
			existingKeys = append(existingKeys, currKPs...)
		} else {
			existingKeys = append(existingKeys, currKPs[:currKPsInd+1]...)
			break
		}
	}
	return existingKeys
}

func (wc *WalletController) EnterMnemonic() bool {
	wc.ClearConsole()
	phraseArr := strings.Split(wc.fieldEnterer("Enter a mnemonic phrase", "Mnemonic phrase must contain words.", isMnemPhraseValid), " ")
	if allowExit {
		allowExit = false
		return false
	}

	println()
	accountName := wc.fieldEnterer("Enter a new account name", "Account name length must be in range from 4 to 25 symbols.", isAccNameValid)
	if allowExit {
		allowExit = false
		return false
	}

	println()
	password := wc.fieldEnterer("Please, enter a new password.", "Error: password type is not correct.", isAccPassValid)
	if allowExit {
		allowExit = false
		return false
	}

	println("Restoring your keys. Please, wait for a while...")

	var mnem mnemonic.Mnemonic
	seed := byteArr.ByteArr{ByteArr: mnem.GenSeed(phraseArr, "")}
	newAcc := account.GenAccount(password, accountName, seed)

	newAcc.KeyPairList = getExistingKeyPairs(seed, 5, password)
	account.Wallet = append(account.Wallet, newAcc)
	account.CurrAccount = account.Wallet[len(account.Wallet)-1]
	account.WriteAccounts()
	return true
}
