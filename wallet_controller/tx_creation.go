package wallet_controller

import (
	"bufio"
	"bvcwallet/account"
	"bvcwallet/byteArr"
	"bvcwallet/cryption"
	"bvcwallet/hashing"
	"bvcwallet/networking"
	"bvcwallet/transaction"
	"fmt"
	"os"
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
			account.AddKeyPairToAccount(command, true)
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

func SendTransaction(tx transaction.Transaction, isSent *bool) {
	var connection networking.Connection
	isEstablished := connection.Establish()

	if isEstablished {
		var isAccepted bool = false
		*isSent = connection.SendTransaction(tx, &isAccepted)
		if isAccepted {
			println("INFO: Transaction is accepted by the node")
		} else {
			println("INFO: An error occured when the node was trying to accept your transaction. Try again later.")
		}
		connection.Close()
	}
}

func UpdateBalance(res *bool) {
	var connection networking.Connection
	isEstablished := connection.Establish()

	if isEstablished {
		*res = connection.GetMyUtxo(account.GetAccAddresses())
		account.GetBalance()
		println("INFO: Balance updated. Type anything to refresh it.")
		connection.Close()
	}
}

func (wc *WalletController) ShowMyAddresses() {
	for true {
		wc.ClearConsole()
		println("Type \"back\" to back to the main menu.")
		fmt.Printf("You have %d addresses on your account\n", len(account.CurrAccount.KeyPairList))

		for _, kp := range account.CurrAccount.KeyPairList {
			println(hashing.SHA1(kp.PublKey))
		}

		wc.scann.Scan()
		command := wc.scann.Text()
		if command == "back" {
			return
		}
	}
}

func (wc *WalletController) ShowMnemonicPhrase() {
	var password string
	reader := bufio.NewReader(os.Stdin)
	for true {
		wc.ClearConsole()
		println("Type \"back\" to back to the main menu.")
		println("Enter your password to confirm the showing.")

		password, _ = reader.ReadString('\n')
		password = strings.Trim(password, wc.getLineSeparator())
		if password == "back" {
			return
		}

		if hashing.SHA1(password) != account.CurrAccount.HashPass {
			println("Password is incorrect. Try again")
		} else {
			break
		}
	}

	mnemStrDecr := cryption.AES_decrypt(account.CurrAccount.MnemonicEncr.ToHexString(), password)
	var mnemDecr byteArr.ByteArr
	mnemDecr.SetFromHexString(mnemStrDecr, len(mnemStrDecr)/2)

	wc.ClearConsole()
	println("Your mnemonic phrase:")
	println(string(mnemDecr.ByteArr))

	println("\nType anything to back")
	reader.ReadString('\n')
}

func (wc *WalletController) handleInput(input string) {
	if input == "0" { // Update balance
		var res bool = true
		go UpdateBalance(&res)
		if res {
			wc.menuMessage = "Updating your balance, it may take some time..."
		} else {
			wc.menuMessage = "An error occured when updating your balance. Try again later"
		}
	} else if input == "1" { // Create transaction
		newTx, isCreated := wc.createTransaction()

		if isCreated == "" {
			var isSent bool
			go SendTransaction(newTx, &isSent)
			if isSent {
				wc.menuMessage = "Your transaction is sent to the node."
			} else {
				wc.menuMessage = "An error occured when sending your transaction. Try again later"
			}
		} else if isCreated == "err" {
			wc.menuMessage = "An error occured when creating your transaction. Try again later"
		}
	} else if input == "2" { // Receive coins
		_, isCreated := wc.createReceiveTransaction()
		if isCreated == "" {
			wc.menuMessage = "Your transaction is created. Now you can send it to the receiver."
		} else if isCreated == "err" {
			wc.menuMessage = "An error occured when creating your transaction. Try again later"
		}
	} else if input == "3" { // Create new address
		wc.createNewAddress()
	} else if input == "4" { // Show my addresses
		wc.ShowMyAddresses()
	} else if input == "5" { // Choose other account
		menuLaunched = false
	} else if input == "6" { // Show mnemonic phrase
		wc.ShowMnemonicPhrase()
	} else if input == "7" { // Exit
		wc.ClearConsole()
		wc.walletLaunched = false
		menuLaunched = false
	}
}

func (wc *WalletController) GetMenu() {
	menuLaunched = true
	wrongInput := false
	for menuLaunched {
		wc.ClearConsole()
		if wc.menuMessage != "" {
			println(wc.menuMessage)
			wc.menuMessage = ""
		}

		fmt.Printf("You balance: %0.8f BVC.\n", float64(account.CurrAccount.Balance)/float64(100000000))
		println("0. Update balance")
		println("1. Send coins (TODO: sign tx)") // TODO: recieve tx and sign it
		println("2. Receive coins")
		println()
		println("3. Create new address")
		println("4. Show my addresses")
		println()
		println("5. Choose the other account")
		println("6. Show mnemonic phrase")
		println("7. Exit")
		println()
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
