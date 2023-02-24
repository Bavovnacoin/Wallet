package wallet_controller

import (
	"bvcwallet/account"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type WalletController struct {
	walletLaunched bool
	opSys          string
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

func (wc *WalletController) Launch() {
	wc.opSys = runtime.GOOS

	if !account.IsWalletExists() {
		wc.ClearConsole()
		println("Can't find any account on you'r PC. Create one") //TODO: or enter a mnemonic phrase
		wc.CreateAccount()
	}

	for wc.walletLaunched {
		return
	}
}
