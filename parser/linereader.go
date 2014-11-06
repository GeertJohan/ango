package parser

import (
	"bufio"
	"io"
	"strings"
)

// lineReader is a simple line reader with peek feature
// not concurrent safe
type lineReader struct {
	ln      int
	bufline *string
	bufrd   *bufio.Reader
}

func newLineReader(rd io.Reader) *lineReader {
	return &lineReader{
		bufrd: bufio.NewReader(rd),
	}
}

// Line currently skips empty lines and removes comments.
// when the parser is more matured it should return all lines,
// especially the comments (maybe empty lines can still be skipped here)..
func (lr *lineReader) Line() (string, error) {
	if lr.bufline != nil {
		line := *lr.bufline
		lr.bufline = nil
		return line, nil
	}
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

func (lr *lineReader) Peek() (string, error) {
	if lr.bufline != nil {
		return *lr.bufline, nil
	}
	line, err := lr.Line()
	if err != nil {
		return "", err
	}
	lr.bufline = &line
	return line, nil
}

func removeComments(line string) string {
	cPos := strings.Index(line, "//")
	if cPos == -1 {
		// no comment, return complete line
		return line
	}
	return line[:cPos]
}
