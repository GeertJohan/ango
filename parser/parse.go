package parser

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/GeertJohan/ango/definitions"
)

// regular expression

var (
	regexpName            *regexp.Regexp
	regexpProcedure       *regexp.Regexp
	regexpParameter       *regexp.Regexp
	regexpSimpleType      *regexp.Regexp
	regexpSliceType       *regexp.Regexp
	regexpMapType         *regexp.Regexp
	regexpStructFieldType *regexp.Regexp
)

func init() {
	var (
		// partials
		rIdentifier     = `[a-z][a-zA-Z0-9]*`
		rBits           = `(?:8|16|32|64)`
		rTypes          = `(?:int` + rBits + `?|uint` + rBits + `?|string)`
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

	regexpSimpleType = regexp.MustCompile(`^` + rIdentifier + `$`)
	regexpSliceType = regexp.MustCompile(`^\[\]` + rIdentifier + `$`)
	regexpMapType = regexp.MustCompile(`^map\[(` + rIdentifier + `)\](` + rIdentifier + `)$`)
	regexpStructFieldType = regexp.MustCompile(`^(` + rIdentifier + `)` + rMustWhitespace + `(` + rIdentifier + `)$`)
}

// ParseError.Type values
var (
	// ParseErrInvalidNameDefinition indicates an invalid name clause
	ParseErrInvalidNameClause = "invalid name clause"

	// ParseErrInvalidProcDefinition indicates an invalid procedure definition
	ParseErrInvalidProcDefinition = "invalid procedure definition"

	// ParseErrInvalidTypeDefinition indicates an invalid type definition
	ParseErrInvalidTypeDefinition = "invalid type definition"

	// ParseErrInvalidStructFieldDefinition indicates an invalid struct field definition
	ParseErrInvalidStructFieldDefinition = "invalid struct field definition"

	// ParseErrMissingType indicates that a type was expected but not found
	ParseErrMissingType = "missing type"

	// ParseErrInvalidDefinition indicates an invalid statement
	ParseErrInvalidStatement = "invalid statement"

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

	// ParseErrUnexpectedEOF indicates that the parse expected more lines, but got EOF
	ParseErrUnexpectedEOF = "unexpected EOF"

	// ParseErrReader indicates there was a problem with reading the syntax
	ParseErrReader = "reader error"
)

var (
	// ErrNotImlemented indicates a feature has not been implemented yet.
	ErrNotImlemented = errors.New("not implemented")

	// ErrUnexpectedEOF indicates the file had not ended yet
	ErrUnexpectedEOF = errors.New(ParseErrUnexpectedEOF)
)

func (parser *Parser) verbosef(format string, data ...interface{}) {
	if parser.config.Verbose {
		fmt.Printf(format, data...)
	}
}

func (parser *Parser) printParseErrorf(format string, data ...interface{}) {
	if parser.config.PrintParseErrors {
		fmt.Printf(format, data...)
	}
}

type Parser struct {
	used    bool
	lr      *lineReader
	config  *Config
	service *definitions.Service
}

type Config struct {
	Verbose          bool
	PrintParseErrors bool
}

func NewParser(config *Config) *Parser {
	return &Parser{
		config: config,
	}
}

// Parse parses an ango definition stream and returns a *Service or an error.
// When and error occured during the parsing of a line, it is of type *ParseError.
func (parser *Parser) Parse(rd io.Reader) (*definitions.Service, error) {
	if parser.used {
		return nil, errors.New("parser can be used only once right now")
	}
	parser.used = true

	// create new service definition
	parser.service = definitions.NewService()

	parser.lr = newLineReader(rd)

	perr := parser.parseName()
	if perr != nil {
		parser.printParseErrorf(perr.Error())
		return nil, perr
	}

	for {
		peekLine, err := parser.lr.Peek()
		if err != nil {
			if err == io.EOF {
				// end of file, parsing is completed
				break
			}
			return nil, err
		}

		switch {
		case strings.HasPrefix(peekLine, "server") || strings.HasPrefix(peekLine, "client"):
			perr := parser.parseProcedure()
			if perr != nil {
				parser.printParseErrorf(perr.Error())
				return nil, perr
			}
		case strings.HasPrefix(peekLine, "type"):
			perr := parser.parseTypeDefinition()
			if perr != nil {
				parser.printParseErrorf(perr.Error())
				return nil, perr
			}
		default:
			perr := &ParseError{
				Type: ParseErrInvalidStatement,
				Line: parser.lr.ln,
			}
			parser.printParseErrorf(perr.Error())
			return nil, perr
		}
	}

	// all done
	return parser.service, nil
}

func (parser *Parser) parseName() *ParseError {
	line, err := parser.lr.Line()
	if err != nil {
		if err == io.EOF {
			return &ParseError{
				Line: parser.lr.ln,
				Type: ParseErrUnexpectedEOF,
			}
		}
		return &ParseError{
			Line:  parser.lr.ln,
			Type:  ParseErrReader,
			Extra: err.Error(),
		}
	}
	matches := regexpName.FindStringSubmatch(line)
	if len(matches) != 2 {
		return &ParseError{
			Line: parser.lr.ln,
			Type: ParseErrInvalidNameClause,
		}
	}

	parser.service.Name = matches[1]

	return nil
}

func (parser *Parser) parseProcedure() *ParseError {
	line, _ := parser.lr.Line() // don't check error, previous line was peeked
	matches := regexpProcedure.FindStringSubmatch(line)
	if len(matches) == 0 {
		return &ParseError{
			Type: ParseErrInvalidProcDefinition,
		}
	}

	proc := &definitions.Procedure{
		Source: definitions.Source{
			Linenumber: parser.lr.ln,
		},
	}
	switch matches[1] {
	case "server":
		proc.Type = definitions.ServerProcedure
	case "client":
		proc.Type = definitions.ClientProcedure
	default:
		panic("unreachable")
	}

	proc.Oneway = (matches[2] == "oneway")
	if proc.Oneway && len(matches[5]) > 0 {
		return &ParseError{
			Type: ParseErrUnexpectedReturnParameters,
		}
	}

	proc.Name = matches[3]

	perr := parser.parseParams(matches[4], &proc.Args)
	if perr != nil {
		return perr
	}
	perr = parser.parseParams(matches[5], &proc.Rets)
	if perr != nil {
		return perr
	}
	if len(matches[5]) > 0 && len(proc.Rets) == 0 {
		return &ParseError{
			Type: ParseErrEmptyReturnGroup,
		}
	}

	var procMap map[string]*definitions.Procedure
	switch proc.Type {
	case definitions.ClientProcedure:
		procMap = parser.service.ClientProcedures
	case definitions.ServerProcedure:
		procMap = parser.service.ServerProcedures
	default:
		panic("unreachable")
	}
	if _, exists := procMap[proc.Name]; exists {
		perr := &ParseError{
			Line:  parser.lr.ln,
			Type:  ParseErrDuplicateProcedureIdentifier,
			Extra: fmt.Sprintf(`"%s"`, proc.Name),
		}
		return perr
	}
	procMap[proc.Name] = proc

	return nil
}

func (parser *Parser) parseParams(text string, list *definitions.Params) *ParseError {
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
	*list = make([]*definitions.Param, 0, paramCount)

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
		p := &definitions.Param{
			Name: name,
		}
		// set typed param type
		switch tipe {
		case "int":
			p.Type = definitions.TypeInt
		case "int8":
			p.Type = definitions.TypeInt8
		case "int16":
			p.Type = definitions.TypeInt16
		case "int32":
			p.Type = definitions.TypeInt32
		case "int64":
			p.Type = definitions.TypeInt64
		case "uint":
			p.Type = definitions.TypeUint
		case "uint8":
			p.Type = definitions.TypeUint8
		case "uint16":
			p.Type = definitions.TypeUint16
		case "uint32":
			p.Type = definitions.TypeUint32
		case "uint64":
			p.Type = definitions.TypeUint64
		case "string":
			p.Type = definitions.TypeString
		default:
			panic("unreachable")
		}

		// append param to params slice on procedure
		*list = append(*list, p)
	}

	return nil
}

func (parser *Parser) parseTypeDefinition() *ParseError {

	line, _ := parser.lr.Line()

	fields := strings.Fields(line)

	if len(fields) < 3 || fields[0] != "type" {
		return parser.newError(ParseErrInvalidTypeDefinition)
	}
	t := &definitions.Type{
		Name: fields[1],
	}

	typetext := fields[2:]

	_, err := parser.parseType(t, typetext)
	if err != nil {
		return err
	}

	parser.service.Types[t.Name] = t
	return nil
}

// parseType parses a type starting with given typetext
// it will read more lines if required (struct)
// if t is non-nil, it will add the type data to t and return t.
// if it is nil, a new *Type will be created (anonymous type).
func (parser *Parser) parseType(t *definitions.Type, typetext []string) (*definitions.Type, *ParseError) {
	if t == nil {
		// anonymous type (no name)
		t = &definitions.Type{}
	}
	if len(typetext) == 0 {
		return nil, parser.newError(ParseErrMissingType)
	}

	switch {
	// struct type
	case typetext[0] == "struct{",
		len(typetext) == 2 && typetext[0] == "struct" && typetext[1] == "{":
		for {
			structLine, err := parser.lr.Line()
			if err != nil {
				return nil, parser.newErrorExtra(ParseErrUnexpectedEOF, "unexpected EOF when parsing struct type `%s`", t.Name)
			}
			if structLine == "}" {
				t.Category = definitions.Struct
				break
			}
			fieldFields := strings.Fields(structLine)
			if len(fieldFields[0]) == 0 || len(fieldFields[1]) == 0 {
				return nil, parser.newError(ParseErrInvalidStructFieldDefinition)
			}
			//++ TODO: verify that fieldFields[0] is a valid field identifier!
			sf := definitions.StructField{
				Name: fieldFields[0],
			}
			var perr *ParseError
			sf.Type, perr = parser.parseType(nil, fieldFields[1:])
			if perr != nil {
				return nil, perr
			}
			t.StructFields = append(t.StructFields, sf)
		}

	// simple type
	case regexpSimpleType.MatchString(typetext[0]):
		simpleTypeName := typetext[0]
		t.SimpleType = parser.service.LookupType(simpleTypeName)
		if t.SimpleType == nil {
			return nil, parser.newErrorExtra(ParseErrInvalidTypeDefinition, "unknown type `%s`", simpleTypeName)
		}
		t.Category = definitions.Simple

	// slice type
	case regexpSliceType.MatchString(typetext[0]):
		elementTypeName := typetext[0][2:]
		t.SliceElementType = parser.service.LookupType(elementTypeName)
		if t.SliceElementType == nil {
			return nil, parser.newErrorExtra(ParseErrInvalidTypeDefinition, "unknown element type `%s`", elementTypeName)
		}
		t.Category = definitions.Slice

	// map type
	case regexpMapType.MatchString(typetext[0]):
		keyValueTypeName := regexpMapType.FindStringSubmatch(typetext[0])
		if len(keyValueTypeName) != 3 {
			return nil, parser.newError(ParseErrInvalidTypeDefinition)
		}
		t.MapKeyType = parser.service.LookupType(keyValueTypeName[1])
		if t.MapKeyType == nil {
			return nil, parser.newErrorExtra(ParseErrInvalidTypeDefinition, "unknown map key type `%s`", keyValueTypeName[0])
		}
		t.MapValueType = parser.service.LookupType(keyValueTypeName[2])
		if t.MapValueType == nil {
			return nil, parser.newErrorExtra(ParseErrInvalidTypeDefinition, "unknown map value type `%s`", keyValueTypeName[1])
		}
		t.Category = definitions.Map

	default:
		// unknown/invalid type definition
		return nil, parser.newError(ParseErrInvalidTypeDefinition)
	}

	return t, nil
}
