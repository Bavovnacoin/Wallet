package wallet_controller

import (
	"bufio"
	"bvcwallet/account"
	"bvcwallet/networking"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type WalletController struct {
	walletLaunched bool
	opSys          string
	scann          *bufio.Scanner

	menuMessage string
}

var allowLaunchMenu bool = false
var isMnemonicEntered bool = false

func (wc *WalletController) ClearConsole() {
	if wc.opSys == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else if wc.opSys == "linux" || wc.opSys == "darwin" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Println("\n\n")
	}
}

func (wc *WalletController) getLineSeparator() string {
	if wc.opSys == "windows" {
		return "\r\n"
	} else if wc.opSys == "darwin" {
		return "\r"
	} else {
		return "\n"
	}
}

func (wc *WalletController) Launch() {
	wc.walletLaunched = true
	wc.opSys = runtime.GOOS
	wc.scann = bufio.NewScanner(os.Stdin)

	for wc.walletLaunched {
		account.IsWalletExists()

		allowLaunchMenu = wc.initAccount()

		var connection networking.Connection
		isEstablished := connection.Establish()

		if isEstablished {
			connection.GetMyUtxo(account.GetAccAddresses())
			account.GetBalance()
			connection.Close()
		} else {
			println("Can't connect to any Bavovnacoin node. Please, try again later")
			return
		}

		if allowLaunchMenu {
			wc.GetMenu()
		}
	}
}
