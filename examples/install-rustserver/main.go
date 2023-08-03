package main

import (
	"github.com/jensvandewiel/gosteamcmd"
	"github.com/jensvandewiel/gosteamcmd/console"
	"os"
)

func main() {
	//this code follows the steps of: https://www.rustafied.com/how-to-host-your-own-rust-server

	prompts := []*gosteamcmd.Prompt{
		gosteamcmd.ForceInstallDir("c:\\rustserver\\"),
		gosteamcmd.Login("", ""),
		gosteamcmd.AppUpdate(258550, "", true),
	}

	cmd := gosteamcmd.New(os.Stdout, prompts)

	cmd.Console.Parser.OnInformationReceived = func(action console.Action, progress float64, currentWritten, total uint64) {
		println("")
	}

	err := cmd.Run()

	if err != nil {
		panic(err)
	}
}
