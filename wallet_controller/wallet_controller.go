package wallet_controller

import (
	"bufio"
	"bvcwallet/account"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type WalletController struct {
	walletLaunched bool
	opSys          string
	scann          *bufio.Scanner
}

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

	allowLaunchMenu := false
	for wc.walletLaunched {
		if !account.IsWalletExists() {
			wc.ClearConsole()
			println("Can't find any account on you'r PC. Create one") //TODO: or enter a mnemonic phrase
			wc.CreateAccount()
			allowLaunchMenu = true
		} else {
			allowLaunchMenu = wc.initAccount()
		}

		if allowLaunchMenu {
			wc.GetMenu()
		}
	}
}
