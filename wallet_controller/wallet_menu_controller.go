package wallet_controller

import (
	"bvcwallet/account"
	"bvcwallet/hashing"
	"fmt"
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

// TODO: add view addresses??
func (wc *WalletController) handleInput(input string) {
	if input == "0" { // Update balance

	} else if input == "1" { // Create transaction

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
