package gosteamcmd

import (
	"github.com/jensvandewiel/gosteamcmd/console"
	"io"
)

type SteamCMD struct {
	// Prompts contains all the commands that will be executed.
	Prompts []*Prompt
	console *console.Console

	stdout io.Writer
}

func New() *SteamCMD {
	return &SteamCMD{
		Prompts: make([]*Prompt, 0),
	}
}

func (s *SteamCMD) RunHeadless() error {
	cmd := "steamcmd"

	for _, prompt := range s.Prompts {
		cmd += " +" + prompt.FullPrompt
	}

	cmd += " +quit"

	s.console = console.New(cmd, s.stdout)
	err := s.console.Run()
	if err != nil {
		return err
	}
	return nil
}
