package gosteamcmd

type SteamCMD struct {
	prompts *[]Prompt
}

func New() *SteamCMD {
	return &SteamCMD{}
}
