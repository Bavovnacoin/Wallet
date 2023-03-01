package wallet_controller

import (
	"bufio"
	"bvcwallet/account"
	"bvcwallet/byteArr"
	"bvcwallet/hashing"
	"bvcwallet/transaction"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var menuLaunched bool

func (wc *WalletController) createNewAddress() {
	isInputIncorrect := false
	for true {
		wc.ClearConsole()
		if isInputIncorrect {
			println("The password is incorrect")
			isInputIncorrect = false
		}

		println("Confirm new address creation (by entering your password) or type \"back\" to back to the menu.")

		wc.scann.Scan()
		command := wc.scann.Text()
		if hashing.SHA1(command) == account.CurrAccount.HashPass {
			wc.ClearConsole()
			account.AddKeyPairToAccount(command)
			fmt.Printf("New address created: %s\n",
				hashing.SHA1(account.CurrAccount.KeyPairList[len(account.CurrAccount.KeyPairList)-1].PublKey))

			println("Type anything to continue")
			wc.scann.Scan()
			wc.scann.Text()
			break
		} else if command == "back" {
			break
		} else {
			isInputIncorrect = true
		}
	}
}

func (wc *WalletController) createTransaction() (transaction.Transaction, bool) {

	reader := bufio.NewReader(os.Stdin)

	var tx transaction.Transaction
	incorrectAmmount := false
	for true {
		wc.ClearConsole()
		if incorrectAmmount {
			println("You do not have enought coins to make a transaction. Try again.")
			incorrectAmmount = false
		}
		println("Transaction creation. To stop creation type \"back\"")
		var outAddr []byteArr.ByteArr
		var outValue []uint64
		println("Type in address and value to be sent separated by a space. Or continue by typing next.")

		for true {
			text, _ := reader.ReadString('\n')
			text = strings.Trim(text, wc.getLineSeparator())
			inputArr := strings.Split(text, " ")
			if text == "next" {
				break
			}
			if text == "back" {
				return tx, false
			}

			var inpAddr byteArr.ByteArr
			isInpCorrect := inpAddr.SetFromHexString(inputArr[0], 20)
			if !isInpCorrect || len(inputArr[0]) != 40 {
				println("You have typed wrong address")
			} else if len(inputArr) == 2 {
				value, err := strconv.ParseUint(inputArr[1], 10, 64)
				if err == nil {
					outValue = append(outValue, value)
				} else {
					println("You have typed wrong coins ammount")
				}
			} else {
				println("You have typed wrong input")
			}
		}

		// TODO: select min, high, middle
		println("Type in tx fee per byte.")
		var fee int
		for true {
			text, _ := reader.ReadString('\n')
			text = strings.Trim(text, wc.getLineSeparator())
			if text == "back" {
				return tx, false
			}

			feeInp, err := strconv.ParseInt(text, 10, 64)
			if err == nil {
				fee = int(feeInp)
				break
			} else {
				println("Error. Expected numerical value.")
			}
		}

		println("Type in how many blocks your transaction will be freezed")
		var locktime uint64
		for true {
			text, _ := reader.ReadString('\n')
			text = strings.Trim(text, wc.getLineSeparator())
			if text == "back" {
				return tx, false
			}

			locktimeInp, err := strconv.ParseUint(text, 10, 64)
			if err == nil {
				locktime = locktimeInp
				break
			} else {
				println("Error. Expected numerical value.")
			}
		}

		println("Enter your password to confirm creation")
		var password string
		for true {
			password, _ = reader.ReadString('\n')
			password = strings.Trim(password, wc.getLineSeparator())
			if password == "back" {
				return tx, false
			}

			if hashing.SHA1(password) != account.CurrAccount.HashPass {
				println("Password is incorrect. Try again")
			} else {
				break
			}
		}

		var res bool
		for !transaction.VerifyTransaction(tx) {
			tx, res = transaction.CreateTransaction(password, outAddr, outValue, fee, uint(locktime))
			if !res {
				incorrectAmmount = true
				break
			}
		}
	}
	return tx, true
}

// TODO: add view addresses??
func (wc *WalletController) handleInput(input string) {
	if input == "0" { // Update balance

	} else if input == "1" { // Create transaction
		_, res := wc.createTransaction()
		if res {
			// wc.sendTransaction(newTx)??
		}
	} else if input == "2" { // Create new address
		wc.createNewAddress()
	} else if input == "3" { // Choose other account
		menuLaunched = false
	} else if input == "4" { // Exit
		wc.ClearConsole()
		println("Thank you for using our wallet. See you!")
		wc.walletLaunched = false
		menuLaunched = false
	}
}

func (wc *WalletController) GetMenu() {
	menuLaunched = true
	wrongInput := false
	for menuLaunched {
		wc.ClearConsole()
		fmt.Printf("You balance: %d BVC.\n", account.CurrAccount.Balance/100000000)
		println("0. Update balance")
		println("1. Send coins")
		println("2. Create new address")
		println("3. Choose the other account")
		println("4. Exit")

		println("Type in a number, to select a function.")
		if wrongInput {
			println("You have typed in wrong value.")
			wrongInput = false
		}
		wc.scann.Scan()
		inValue := wc.scann.Text()
		wc.handleInput(inValue)
	}
}
