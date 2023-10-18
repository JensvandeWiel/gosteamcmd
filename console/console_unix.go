//go:build unix

package console

import (
	"fmt"
	"io"
	"os/exec"
)

type Console struct {
	stdout      io.Writer // stdout is the writer where the output will be written to.
	commandLine string    // commandLine is the command that will be executed.
	proc        *exec.Cmd
	exitCode    uint32
	args        []string
	Parser      *Parser
}

// Run executes the command in steamcmd and returns the exit code. Exit code does not need to be 0 to return no errors (error is for executing the pseudoconsole)
func (c *Console) Run() (uint32, error) {
	var err error

	if c.commandLine == "" {
		_, err := exec.LookPath(c.commandLine)
		if err != nil {
			return 1, fmt.Errorf("Steamcmd not found: %v\n", err)
		}
	}

	var a []string
	a = append(a, c.commandLine)
	a = append(a, c.args...)

	c.proc = exec.Command("sh", a...)
	if err != nil {
		return 1, err
	}

	d := &duplicateWriter{
		writer1: c.Parser,
		writer2: c.stdout,
	}

	c.proc.Stdout = d
	c.proc.Stderr = d

	err = c.proc.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			c.exitCode = uint32(exitErr.ExitCode())
		} else {
			return 1, err
		}
	}
	c.exitCode = uint32(c.proc.ProcessState.ExitCode())

	return c.exitCode, nil
}
