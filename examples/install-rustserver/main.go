package main

import "github.com/jensvandewiel/gosteamcmd"

func main() {
	//this code follows the steps of: https://www.rustafied.com/how-to-host-your-own-rust-server
	cmd := gosteamcmd.New()
	cmd.Prompts = append(cmd.Prompts, gosteamcmd.ForceInstallDir("c:\\rustserver\\"))
	cmd.Prompts = append(cmd.Prompts, gosteamcmd.Login("", ""))
	cmd.Prompts = append(cmd.Prompts, gosteamcmd.AppUpdate(258550, "", false))
	//running it headless means it will not output anything
	err := cmd.RunHeadless()

	if err != nil {
		panic(err)
	}
}
