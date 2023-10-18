package console

import "io"

func New(exePath string, args []string, stdout io.Writer) *Console {
	p := NewParser()

	return &Console{
		commandLine: exePath,
		stdout:      stdout,
		args:        args,
		Parser:      p,
	}
}
