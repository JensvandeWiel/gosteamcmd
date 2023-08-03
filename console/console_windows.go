//go:build windows

package console

import (
	"context"
	"github.com/UserExistsError/conpty"
	"io"
)

//Todo console should only be used if not headless otherwise just use os.exec

type Console struct {
	stdout      io.Writer
	commandLine string
	conPTY      *conpty.ConPty
	ExitCode    uint32
	Parser      *Parser
}

func New(exePath string, stdout io.Writer) *Console {
	return &Console{
		commandLine: exePath,
		stdout:      stdout,
	}
}

func (c *Console) Run() error {
	var err error
	c.conPTY, err = conpty.Start(c.commandLine)
	if err != nil {
		return err
	}
	defer c.conPTY.Close()

	go io.Copy(io.MultiWriter(c.stdout, c.Parser), c.conPTY)

	c.ExitCode, err = c.conPTY.Wait(context.Background())
	if err != nil {
		return err
	}

	return nil
}
