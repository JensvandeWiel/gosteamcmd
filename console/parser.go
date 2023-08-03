package console

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type Action int

const (
	Downloading Action = iota
	Verifying
	Preallocating
)

// Parser is a parser for the steamcmd console output. Right now it only parses the progress of downloading, verifying and preallocating. To capture other output use the stdout.
type Parser struct {
	OnInformationReceived func(action Action, progress float64, currentWritten, total uint64)
}

// NewParser creates a new Parser instance. this must be used, if not the parser will not work.
func NewParser() *Parser {
	// Set empty function because otherwise it will panic, because not every user will use this.
	return &Parser{
		OnInformationReceived: func(action Action, progress float64, currentWritten, total uint64) {},
	}
}

func (p *Parser) Write(data []byte) (int, error) {
	switch {
	case strings.Contains(string(data), "verifying"):
		progress, current, total, err := extractProgressAndNumbers(string(data))
		if err != nil {
			return 0, err
		}
		go p.OnInformationReceived(Verifying, progress, current, total)
	case strings.Contains(string(data), "downloading"):
		progress, current, total, err := extractProgressAndNumbers(string(data))
		if err != nil {
			return 0, err
		}
		go p.OnInformationReceived(Downloading, progress, current, total)
	case strings.Contains(string(data), "preallocating"):
		progress, current, total, err := extractProgressAndNumbers(string(data))
		if err != nil {
			return 0, err
		}
		go p.OnInformationReceived(Preallocating, progress, current, total)
	}
	return len(data), nil
}

func extractProgressAndNumbers(input string) (progress float64, currentWritten, total uint64, err error) {
	progressPattern := `progress:\s+(\d+\.\d+)\s+\((\d+)\s+/\s+(\d+)\)`

	re := regexp.MustCompile(progressPattern)
	match := re.FindStringSubmatch(input)

	if len(match) != 4 {
		err = fmt.Errorf("Progress not found in the provided line.")
		return
	}

	progressStr := match[1]
	currentStr := match[2]
	totalStr := match[3]

	progress, err = strconv.ParseFloat(progressStr, 64)
	if err != nil {
		err = fmt.Errorf("Error parsing progress value: %w", err)
		return
	}

	currentWritten, err = strconv.ParseUint(currentStr, 10, 64)
	if err != nil {
		err = fmt.Errorf("Error parsing currentWritten value: %w", err)
		return
	}

	total, err = strconv.ParseUint(totalStr, 10, 64)
	if err != nil {
		err = fmt.Errorf("Error parsing total value: %w", err)
		return
	}

	return progress, currentWritten, total, nil
}

// region duplicateWriter

// duplicateWriter duplicates the output to two writers. Io.MultiWriter does not work for some reason. (prob because parser takes longer)
type duplicateWriter struct {
	writer1 io.Writer
	writer2 io.Writer
}

// Write writes the data to both writers 1 and 2.
func (d *duplicateWriter) Write(p []byte) (n int, err error) {
	go func() {
		d.writer1.Write(p)
	}()
	go func() {
		d.writer2.Write(p)
	}()
	return len(p), nil
}

// endregion
