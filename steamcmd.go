package gosteamcmd

import (
	"github.com/jensvandewiel/gosteamcmd/console"
	"io"
)

type SteamCMD struct {
	// prompts contains all the commands that will be executed.
	prompts []*Prompt
	Console *console.Console

	Stdout io.Writer
}

// New creates a new SteamCMD instance.
func New(stdout io.Writer, prompts []*Prompt) *SteamCMD {

	s := &SteamCMD{
		prompts: prompts,
		Stdout:  stdout,
	}

	//prepare command
	cmd := "steamcmd"
	for _, prompt := range s.prompts {
		cmd += " +" + prompt.FullPrompt
	}
	cmd += " +quit"
	s.Console = console.New(cmd, s.Stdout)

	return s
}

// Run executes the SteamCMD instance.
func (s *SteamCMD) Run() (uint32, error) {
	return s.Console.Run()
}
