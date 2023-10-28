package gosteamcmd

import "strconv"

type PromptType int

const (
	AppUpdatePrompt PromptType = iota
	LoginPrompt
	ForceInstallDirPrompt
)

type Prompt struct {
	FullPrompt string
	Type       PromptType
}

// AppUpdate updates the given appID. If beta is not empty, it will update to the given beta. If validate is true, it will validate the files.
func AppUpdate(appID int, beta string, validate bool) *Prompt {
	cmd := "app_update "
	cmd += strconv.Itoa(appID) + " "
	if beta != "" {
		cmd += "-beta " + beta + " "
	}
	if validate {
		cmd += "validate"
	}
	return &Prompt{cmd, AppUpdatePrompt}
}

// Login logs into SteamCMD with the given username and password. If the arguments are empty strings, it will login as anonymous.
func Login(username string, password string, authcode string) *Prompt {
	cmd := "login "
	if username != "" {
		cmd += username + " " + password + " " + authcode + " "
	} else {
		cmd += "anonymous"
	}
	return &Prompt{cmd, LoginPrompt}
}

func ForceInstallDir(directory string) *Prompt {
	cmd := "force_install_dir "
	cmd += directory

	return &Prompt{cmd, ForceInstallDirPrompt}
}
