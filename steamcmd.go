package gosteamcmd

import (
	"github.com/jensvandewiel/gosteamcmd/console"
	"io"
)

type SteamCMD struct {
	// Prompts contains all the commands that will be executed.
	Prompts []*Prompt
	console *console.Console

	Stdout io.Writer
}

// New creates a new SteamCMD instance.
func New(stdout io.Writer) *SteamCMD {
	return &SteamCMD{
		Prompts: make([]*Prompt, 0),
		Stdout:  stdout,
	}
}

// Run puts all the prompts together and executes them.
func (s *SteamCMD) Run() error {
	cmd := "steamcmd"

	for _, prompt := range s.Prompts {
		cmd += " +" + prompt.FullPrompt
	}

	cmd += " +quit"
	s.console = console.New(cmd, s.Stdout)
	err := s.console.Run()
	if err != nil {
		return err
	}
	return nil
}
