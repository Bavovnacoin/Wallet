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
			if len(inputArr[0]) != 40 {
				println("You have typed wrong address")
			} else if len(inputArr) == 2 && len(inputArr[0]) == 40 {
				isInpCorrect := inpAddr.SetFromHexString(inputArr[0], 20)
				if !isInpCorrect {
					println("You have typed wrong address")
					continue
				}
				value, err := strconv.ParseUint(inputArr[1], 10, 64)
				if err == nil {
					outValue = append(outValue, value)
					outAddr = append(outAddr, inpAddr)
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
		for isTxCorrect := false; !isTxCorrect; {
			tx, res = transaction.CreateTransaction(password, outAddr, outValue, fee, uint(locktime))
			if !res {
				incorrectAmmount = true
				break
			}
			isTxCorrect = transaction.VerifyTransaction(tx)

			if isTxCorrect {
				account.WriteAccounts()
				return tx, true
			} else { // Remove last keypair if tx is incorrect
				kpLen := len(account.CurrAccount.KeyPairList)
				account.CurrAccount.KeyPairList = append(account.CurrAccount.KeyPairList, account.CurrAccount.KeyPairList[:kpLen-1]...)
			}
		}
	}
	return tx, false
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
		println("INFO: Balance updated. Press any button to refresh it.")
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

// TODO: add view addresses??
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

		if isCreated {
			var isSent bool
			go SendTransaction(newTx, &isSent)
			if isSent {
				wc.menuMessage = "Your transaction is sent to the node."
			} else {
				wc.menuMessage = "An error occured when sending your transaction. Try again later"
			}
		} else {
			wc.menuMessage = "An error occured when creating your transaction. Try again later"
		}
	} else if input == "2" { // Create new address
		wc.createNewAddress()
	} else if input == "3" { // Show my addresses
		wc.ShowMyAddresses()
	} else if input == "4" { // Choose other account
		menuLaunched = false
	} else if input == "5" { // Show mnemonic phrase
		wc.ShowMnemonicPhrase()
	} else if input == "6" { // Exit
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
		if wc.menuMessage != "" {
			println(wc.menuMessage)
			wc.menuMessage = ""
		}

		fmt.Printf("You balance: %0.8f BVC.\n", float64(account.CurrAccount.Balance)/float64(100000000))
		println("0. Update balance")
		println("1. Send coins")
		println("2. Create new address")
		println("3. Show my addresses")
		println("4. Choose the other account")
		println("5. Show mnemonic phrase")
		println("6. Exit")

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
