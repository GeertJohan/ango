package parser

import (
	"bufio"
	"io"
	"strings"
)

type lineReader struct {
	ln    int
	bufrd *bufio.Reader
}

func newLineReader(rd io.Reader) *lineReader {
	return &lineReader{
		bufrd: bufio.NewReader(rd),
	}
}

func (lr *lineReader) Line() (string, error) {
	for {
		lr.ln++
		line, err := lr.bufrd.ReadString('\n')
		if err != nil {
			return "", err
		}

		line = removeComments(line)

		// remove whitespace
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		return line, nil
	}
}

func removeComments(line string) string {
	cPos := strings.Index(line, "//")
	if cPos == -1 {
		// no comment, return complete line
		return line
	}
	return line[:cPos]
}
