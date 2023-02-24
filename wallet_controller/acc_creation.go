package wallet_controller

import (
	"bufio"
	"bvcwallet/account"
	"bvcwallet/mnemonic"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	specSign = " !”#$%&’()*+,./:;<=>?@[/^_`{}|~"
)

var scann *bufio.Scanner
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

func fieldEnterer(enterText string, errorText string, validation func(val string) bool) string {
	println("Type \"back\" to back to the menu.")
	println(enterText)
	enterValue := ""
	for true {
		if enterValue != "" {
			println(errorText)
			println(enterValue)
		}

		scann.Scan()
		enterValue = scann.Text()

		if enterValue == "back" {
			allowExit = true
			return ""
		}
		if validation(enterValue) {
			break
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

func (wc *WalletController) CreateAccount() {
	scann = bufio.NewScanner(os.Stdin)
	accountName := fieldEnterer("Enter account name", "Account name length must be in range from 4 to 25 symbols.", isAccNameValid)
	if allowExit {
		allowExit = false
		return
	}

	wc.ClearConsole()
	accountPass := fieldEnterer("Enter account password", fmt.Sprintf("Account password length must be in range from 8 to 20 symbols. It should contain upper and lower case letters, numbers and special signs (%s).\n", specSign), isAccPassValid)
	if allowExit {
		return
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
}
