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

	"github.com/sqweek/dialog"
)

func (wc *WalletController) createTransaction() (transaction.Transaction, string) {
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
				return tx, "back"
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
				return tx, "back"
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
				return tx, "back"
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
				return tx, "back"
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
				return tx, ""
			} else { // Remove last keypair if tx is incorrect
				kpLen := len(account.CurrAccount.KeyPairList)
				account.CurrAccount.KeyPairList = append(account.CurrAccount.KeyPairList, account.CurrAccount.KeyPairList[:kpLen-1]...)
			}
		}
	}
	return tx, "err"
}

func genTxFileName(tx transaction.Transaction) string {
	txHash := hashing.SHA1(transaction.GetCatTxFields(tx))
	return fmt.Sprintf("tx-%s.bvctx", txHash)
}

func WriteTx(tx transaction.Transaction, path string) bool {
	txByteArr, isConv := byteArr.ToByteArr(tx)
	if !isConv {
		return false
	}

	f, isOpen := os.Create(path + "/" + genTxFileName(tx))
	if isOpen != nil {
		return false
	}

	f.Write(txByteArr)
	f.Close()
	return true
}

func getFileSavePath() string {
	dirPath, err := dialog.Directory().Title("Select a directory to save a tx").Browse()
	if err != nil {
		return ""
	}
	return strings.ReplaceAll(dirPath, "\\", "/")
}

func isAddrInAccount(addr byteArr.ByteArr) bool {
	for i := 0; i < len(account.CurrAccount.KeyPairList); i++ {
		var accAddr byteArr.ByteArr
		accAddr.SetFromHexString(hashing.SHA1(account.CurrAccount.KeyPairList[i].PublKey), 20)

		if accAddr.IsEqual(addr) {
			return true
		}
	}
	return false
}

func (wc *WalletController) createReceiveTransaction() (transaction.Transaction, string) {
	reader := bufio.NewReader(os.Stdin)

	var tx transaction.Transaction
	for true {
		wc.ClearConsole()
		println("Transaction creation for coins receiving. To stop creation type \"back\"")
		var outAddr []byteArr.ByteArr
		var outValue []uint64
		println("Type in your address and value to be received separated by a space. Or continue by typing next.")

		for true {
			text, _ := reader.ReadString('\n')
			text = strings.Trim(text, wc.getLineSeparator())
			inputArr := strings.Split(text, " ")
			if text == "next" {
				break
			}
			if text == "back" {
				return tx, "back"
			}

			var inpAddr byteArr.ByteArr
			if len(inputArr[0]) != 40 {
				println("You have typed wrong address")
			} else if len(inputArr) == 2 && len(inputArr[0]) == 40 {
				if inpAddr.SetFromHexString(inputArr[0], 20) {
					println("You have typed wrong address")
					continue
				} else if !isAddrInAccount(inpAddr) {
					println("You don't have such an address")
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

		println("Type in how many blocks your transaction will be freezed")
		var locktime uint
		for true {
			text, _ := reader.ReadString('\n')
			text = strings.Trim(text, wc.getLineSeparator())
			if text == "back" {
				return tx, "back"
			}

			locktimeInp, err := strconv.ParseUint(text, 10, 64)
			if err == nil {
				locktime = uint(locktimeInp)
				break
			} else {
				println("Error. Expected numerical value.")
			}
		}

		var tx transaction.Transaction
		tx.Locktime = locktime
		for i := 0; i < len(outAddr); i++ {
			tx.Outputs = append(tx.Outputs, transaction.Output{Address: outAddr[i], Value: outValue[i]})
		}

		savePath := getFileSavePath()
		if savePath == "" {
			break
		}

		if WriteTx(tx, savePath) {
			return tx, ""
		} else {
			break
		}

	}
	return tx, "err"
}
