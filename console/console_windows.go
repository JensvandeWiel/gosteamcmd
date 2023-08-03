//go:build windows

package console

import (
	"context"
	"github.com/UserExistsError/conpty"
	"io"
)

type Console struct {
	stdout      io.Writer // stdout is the writer where the output will be written to.
	commandLine string    // commandLine is the command that will be executed.
	conPTY      *conpty.ConPty
	exitCode    uint32
	Parser      *Parser
}

func New(exePath string, stdout io.Writer) *Console {
	p := NewParser()

	return &Console{
		commandLine: exePath,
		stdout:      stdout,
		Parser:      p,
	}
}

// Run executes the command in steamcmd and returns the exit code. Exit code does not need to be 0 to return no errors (error is for executing the pseudoconsole)
func (c *Console) Run() (uint32, error) {
	var err error
	c.conPTY, err = conpty.Start(c.commandLine)
	if err != nil {
		return 1, err
	}
	defer c.conPTY.Close()

	d := &duplicateWriter{
		writer1: c.Parser,
		writer2: c.stdout,
	}

	go io.Copy(d, c.conPTY)

	c.exitCode, err = c.conPTY.Wait(context.Background())
	if err != nil {
		return 1, err
	}

	return c.exitCode, nil
}
