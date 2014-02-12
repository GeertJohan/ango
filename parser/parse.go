package parser

import (
	"errors"
	"fmt"
	"io"
)

var (
	// ErrNotImlemented indicates a feature has not been implemented yet.
	ErrNotImlemented = errors.New("not implemented")
)

// Verbose, when true this package will send verbose information to stdout.
var Verbose bool

func verbosef(format string, data ...interface{}) {
	if Verbose {
		fmt.Printf(format, data...)
	}
}

// Parse parses an ango definition stream and returns a *Service or an error.
func Parse(rd io.Reader) (*Service, error) {
	verbosef("do stuff with reader\n")
	return nil, ErrNotImlemented
}
