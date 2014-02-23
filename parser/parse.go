package parser

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// regular expression

var (
	regexpName      *regexp.Regexp
	regexpProcedure *regexp.Regexp
	regexpParameter *regexp.Regexp
)

func init() {
	var (
		// partials
		rIdentifier     = `[a-z][a-zA-Z0-9]*`
		rTypes          = `(?:int|uint|string)`
		rMustWhitespace = `[ \t]+`
		rOptWhitepsace  = `[ \t]*`

		// used to capture name definition
		rNameCapture = `^name` + rMustWhitespace + `(` + rIdentifier + `)$`

		// used to capture procedure definition (and partials)
		rParam            = rOptWhitepsace + rIdentifier + rMustWhitespace + rTypes + rOptWhitepsace
		rParameters       = `\((?:` + rParam + `)?(?:,` + rParam + `)*\)`
		rProcedureCapture = `^(server|client)` + rMustWhitespace + `(?:(oneway)` + rMustWhitespace + `)?(` + rIdentifier + `)` + rOptWhitepsace + `(` + rParameters + `)` + rOptWhitepsace + `((?:` + rParameters + `)?)` + rOptWhitepsace + `$`

		// used to capture parameters
		rParamCapture = rOptWhitepsace + `(` + rIdentifier + `)` + rMustWhitespace + `(` + rTypes + `)` + rOptWhitepsace
	)
	regexpName = regexp.MustCompile(rNameCapture)
	regexpProcedure = regexp.MustCompile(rProcedureCapture)
	regexpParameter = regexp.MustCompile(rParamCapture)
}

var (
	// ErrNotImlemented indicates a feature has not been implemented yet.
	ErrNotImlemented = errors.New("not implemented")

	// ErrUnexpectedEOF indicates the file had not ended yet
	ErrUnexpectedEOF = errors.New("unexpected end of file")
)

// ParseError.Type values
var (
	// ParseErrInvalidNameDefinition indicates an invalid name definition
	ParseErrInvalidNameDefinition = "invalid name definition"

	// ParseErrInvalidProcDefinition indicates an invalid procedure definition
	ParseErrInvalidProcDefinition = "invalid procedure definition"

	// ParseErrInvalidParameter indicates an invalid parameter definition (argument or return value)
	ParseErrInvalidParameter = "invalid parameter definition (argument or return value)"

	// ParseErrDuplicateProcedureIdentifier indicates a duplicate identifier for a procedure
	ParseErrDuplicateProcedureIdentifier = "duplicate procedure identifier"

	// ParseErrDuplicateParameterIdentifier indicates a duplicate parameter identifier (argument or return value)
	ParseErrDuplicateParameterIdentifier = "duplicate parameter identifier (argument or return value)"

	// ParseErrUnexpectedReturnParameters indicates that return parameters were given.
	// This is probably unexpected because the procedure is a oneway procedure.
	ParseErrUnexpectedReturnParameters = "unexpected return parameters (oneway procedure?)"

	// ParseErrEmptyReturnGroup indicates parenthesis for return values are given, but no actual return parameters inside them.
	ParseErrEmptyReturnGroup = "empty return group"
)

// Verbose, when true this package will send verbose information to stdout.
var Verbose = false

func verbosef(format string, data ...interface{}) {
	if Verbose {
		fmt.Printf(format, data...)
	}
}

// PrintParseErrors can be set to false to disable the printing of a parse error
var PrintParseErrors = true

func printParseErrorf(format string, data ...interface{}) {
	if PrintParseErrors {
		fmt.Printf(format, data...)
	}
}

// Parse parses an ango definition stream and returns a *Service or an error.
// When and error occured during the parsing of a line, it is of type *ParseError.
func Parse(rd io.Reader) (*Service, error) {
	verbosef("do stuff with reader\n")

	service := newService()

	lr := newLineReader(rd)

	line, err := lr.Line()
	if err != nil {
		if err == io.EOF {
			return nil, ErrUnexpectedEOF
		}
		return nil, err
	}

	var perr *ParseError
	service.Name, perr = findName(line)
	if perr != nil {
		perr.Line = lr.ln
		printParseErrorf(perr.Error())
		return nil, perr
	}

	for {
		line, err := lr.Line()
		if err != nil {
			if err == io.EOF {
				// end of file, parsing is completed
				break
			}
			return nil, err
		}

		proc, perr := findProcedure(line)
		if perr != nil {
			perr.Line = lr.ln
			printParseErrorf(perr.Error())
			return nil, perr
		}
		var procMap map[string]*Procedure
		switch proc.Type {
		case ClientProcedure:
			procMap = service.ClientProcedures
		case ServerProcedure:
			procMap = service.ServerProcedures
		default:
			panic("unreachable")
		}
		if _, exists := procMap[proc.Name]; exists {
			perr := &ParseError{
				Line:  lr.ln,
				Type:  ParseErrDuplicateProcedureIdentifier,
				Extra: fmt.Sprintf(`"%s"`, proc.Name),
			}
			printParseErrorf(perr.Error())
			return nil, perr
		}
		procMap[proc.Name] = proc
	}

	// all done
	return service, nil
}

func findName(line string) (string, *ParseError) {
	matches := regexpName.FindStringSubmatch(line)
	if len(matches) != 2 {
		return "", &ParseError{
			Type: ParseErrInvalidNameDefinition,
		}
	}

	return matches[1], nil
}

func findProcedure(line string) (*Procedure, *ParseError) {
	matches := regexpProcedure.FindStringSubmatch(line)
	if len(matches) == 0 {
		return nil, &ParseError{
			Type: ParseErrInvalidProcDefinition,
		}
	}

	proc := &Procedure{}
	switch matches[1] {
	case "server":
		proc.Type = ServerProcedure
	case "client":
		proc.Type = ClientProcedure
	default:
		panic("unreachable")
	}

	proc.Oneway = (matches[2] == "oneway")
	if proc.Oneway && len(matches[5]) > 0 {
		return nil, &ParseError{
			Type: ParseErrUnexpectedReturnParameters,
		}
	}

	proc.Name = matches[3]

	perr := parseParams(matches[4], &proc.Args)
	if perr != nil {
		return nil, perr
	}
	perr = parseParams(matches[5], &proc.Rets)
	if perr != nil {
		return nil, perr
	}
	if len(matches[5]) > 0 && len(proc.Rets) == 0 {
		return nil, &ParseError{
			Type: ParseErrEmptyReturnGroup,
		}
	}

	return proc, nil
}

func parseParams(text string, list *Params) *ParseError {
	if len(text) < 3 {
		// fast return for no params or ()
		return nil
	}

	// map holding taken identifiers for this param set
	taken := make(map[string]bool)

	// remove ( and )
	text = text[1 : len(text)-1]
	// split on comma
	paramStrings := strings.Split(text, ",")

	// count and create slice for params
	paramCount := len(paramStrings)
	*list = make([]*Param, 0, paramCount)

	// loop over params
	for i, paramString := range paramStrings {
		// find match and verify
		matches := regexpParameter.FindStringSubmatch(paramString)
		if len(matches) != 3 {
			return &ParseError{
				Type:  ParseErrInvalidParameter,
				Extra: fmt.Sprintf(`at position %d: "%s"`, i+1, paramString),
			}
		}

		// get name and type
		name := matches[1]
		tipe := matches[2]

		// check if name (identifier) is taken
		if taken[name] {
			return &ParseError{
				Type:  ParseErrDuplicateParameterIdentifier,
				Extra: fmt.Sprintf(`at position %d: "%s"`, i+1, name),
			}
		}
		taken[name] = true

		// create new param
		p := &Param{
			Name: name,
		}
		// set typed param type
		switch tipe {
		case "int":
			p.Type = ParamTypeInt
		case "uint":
			p.Type = ParamTypeUint
		case "string":
			p.Type = ParamTypeString
		default:
			panic("unreachable")
		}

		// append param to params slice on procedure
		*list = append(*list, p)
	}

	return nil
}
