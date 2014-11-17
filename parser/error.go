package parser

import "fmt"

// ParseError holds information about an error on a given line.
// ParseError implements the error interface.
type ParseError struct {
	Line  int
	Type  string
	Extra string
}

func (pe *ParseError) Error() string {
	if len(pe.Extra) > 0 {
		return fmt.Sprintf("%s at line %d: %s", pe.Type, pe.Line, pe.Extra)
	}
	return fmt.Sprintf("%s at line %d", pe.Type, pe.Line)
}

func (p *Parser) newError(errType string) *ParseError {
	return &ParseError{
		Line: p.lr.ln,
		Type: errType,
	}
}

func (p *Parser) newErrorExtra(errType string, format string, args ...interface{}) *ParseError {
	return &ParseError{
		Line:  p.lr.ln,
		Type:  errType,
		Extra: fmt.Sprintf(format, args...),
	}
}
