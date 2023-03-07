package wallet_controller

import (
	"bvcwallet/account"
	"bvcwallet/byteArr"
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

	seedString := byteArrToString(mn.GenSeed(phrase, accountPass))
	account.Wallet = append(account.Wallet, account.GenAccount(accountName, accountPass, seedString))
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

		var command string
		wc.scann.Scan()
		command = wc.scann.Text()

		accIdNum, err := strconv.Atoi(command)
		if command == "exit" {
			wc.walletLaunched = false
			return false
		} else if err != nil || accIdNum < 0 || accIdNum >= len(account.Wallet) {
			printErr = true
		} else {
			validResult := wc.validateUser(accIdNum)
			if validResult {
				account.InitAccount(accIdNum)
				return true
			}
		}
	}
	return false
}

func isMnemPhraseValid(phrase string) bool {
	phraseArr := strings.Split(phrase, " ")
	if len(phraseArr) != 24 {
		return false
	}
	return true
}

// Binary search for checking existing keys
func checkExistingKeyPairs(kps []ecdsa.KeyPair) int {
	l := 0
	r := len(kps) - 1
	mid := 0

	var con networking.Connection
	con.Establish()
	defer con.Close()

	if con.IsAddrExist(hashing.SHA1(kps[r].PublKey)) {
		return r
	}

	for true {
		mid = (r + l) / 2
		checkRes := con.IsAddrExist(hashing.SHA1(kps[mid].PublKey))

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

// Checking an existing keys in the network
func getExistingKeyPairs(seed byteArr.ByteArr, keysCheckAmmount int) []ecdsa.KeyPair {
	var existingKeys []ecdsa.KeyPair
	kpInd := 0

	for true {
		var currKPs []ecdsa.KeyPair
		for ; kpInd < keysCheckAmmount; kpInd += 1 {
			currKPs = append(currKPs, ecdsa.GenKeyPair(seed, kpInd))
		}

		currKPsInd := checkExistingKeyPairs(currKPs)
		if currKPsInd < 0 {
			break
		} else if currKPsInd == len(currKPs)-1 {
			existingKeys = append(existingKeys, currKPs...)
		} else {
			existingKeys = append(existingKeys, currKPs[:currKPsInd]...)
			break
		}
	}
	return existingKeys
}

func (wc *WalletController) EnterMnemonic() bool {
	wc.ClearConsole()
	phraseArr := strings.Split(wc.fieldEnterer("Enter a mnemonic phrase", "Mnemonic phrase must be 24 words long.", isMnemPhraseValid), " ")
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

	println(fmt.Sprint(phraseArr))
	println(accountName)
	println(password)

	var mnem mnemonic.Mnemonic
	seed := byteArr.ByteArr{ByteArr: mnem.GenSeed(phraseArr, password)}
	println(seed.ToString())
	newAcc := account.GenAccount(password, accountName, seed.ToString())
	newAcc.KeyPairList = getExistingKeyPairs(seed, 10)
	return true
}
