package console

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Parser contains info about the running steamcmd/console instance.
type Parser struct {
	// OnInfoUpdated is called when the status of the download/verifying changes. action is what it is doing (downloading/verifying), bytesCompleted is the amount of bytes downloaded/verified and bytesTotal is the total amount of bytes to download/verify. This won't be called when command is not app_update or anything related
	OnInfoUpdated func(action string, bytesCompleted uint64, bytesTotal uint64, progress float32)
	// OnSteamCMDOutput is called when steamcmd prints a new line.
	OnSteamCMDOutput func(o []byte)

	// RawOutput is called whenever there is written to the pseudoterminal by steamcmd
	RawOutput func(o []byte)

	buffer bytes.Buffer
}

// TODO refactor complexity
func (p *Parser) Write(b []byte) (n int, err error) {
	n, err = p.buffer.Write(b)
	if err != nil {
		return n, err
	}

	// Process complete lines from the buffer
	for {
		line, err := p.buffer.ReadString('\n')
		if err != nil {
			// If the error is not EOF, it means we encountered an issue
			// during reading, so return the error.
			if err != io.EOF {
				return n, err
			}
			// If EOF is encountered, it might mean that the last line is not
			// terminated with a newline. In this case, we can break the loop.
			break
		}
		p.OnSteamCMDOutput([]byte(line))
		//do the parsing here
		if strings.Contains(line, "downloading") {
			progressAndNumbers, current, total, err := extractProgressAndNumbers(line)
			if err != nil {
				return 0, err
			}
			p.OnInfoUpdated("downloading", current, total, float32(progressAndNumbers))
		} else if strings.Contains(line, "verifying") {
			progressAndNumbers, current, total, err := extractProgressAndNumbers(line)
			if err != nil {
				return 0, err
			}
			p.OnInfoUpdated("verifying", current, total, float32(progressAndNumbers))
		}

		n += len(line)
	}

	// Reset the buffer with any remaining data (not a complete line).
	// The unread data will be used for processing in the next Write call.
	p.buffer.Reset()
	return n, nil
}

// extractProgressAndNumbers extracts the progress, current and total bytes from a line
func extractProgressAndNumbers(input string) (progress float64, current, total uint64, err error) {
	progressPattern := `progress:\s+(\d+\.\d+)\s+\((\d+)\s+/\s+(\d+)\)`

	re := regexp.MustCompile(progressPattern)
	match := re.FindStringSubmatch(input)

	if len(match) != 4 {
		err = fmt.Errorf("values not found in the provided line")
		return
	}

	progressStr := match[1]
	currentStr := match[2]
	totalStr := match[3]

	progress, err = strconv.ParseFloat(progressStr, 64)
	if err != nil {
		err = fmt.Errorf("error parsing progress value: %w", err)
		return
	}

	current, err = strconv.ParseUint(currentStr, 10, 64)
	if err != nil {
		err = fmt.Errorf("error parsing current value: %w", err)
		return
	}

	total, err = strconv.ParseUint(totalStr, 10, 64)
	if err != nil {
		err = fmt.Errorf("error parsing total value: %w", err)
		return
	}

	return progress, current, total, nil
}
