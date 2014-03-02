package definitions

import (
	"errors"
	"math"
	"strings"
)

// ParamType is the type of a parameter
// It's value should be a valid go TypeName (http://golang.org/ref/spec#TypeName)
type ParamType string

// ParamType's
//++ TODO: make ParamType an int and add switch to convert to string?
var (
	ParamTypeInt    = ParamType("int")
	ParamTypeInt8   = ParamType("int8")
	ParamTypeInt16  = ParamType("int16")
	ParamTypeInt32  = ParamType("int32")
	ParamTypeInt64  = ParamType("int64")
	ParamTypeUint   = ParamType("uint")
	ParamTypeUint8  = ParamType("uint8")
	ParamTypeUint16 = ParamType("uint16")
	ParamTypeUint32 = ParamType("uint32")
	ParamTypeUint64 = ParamType("uint64")
	ParamTypeString = ParamType("string")
	ParamTypeBool   = ParamType("bool")
)

// Param defines an argument or return parameter
type Param struct {
	Name string
	Type ParamType
}

// CapitalizedName returns the name for this param, capitalized
func (p *Param) CapitalizedName() string {
	return strings.ToUpper(p.Name[:1]) + p.Name[1:]
}

// GoTypeName returns the Go TypeName for the ParamType
func (p Param) GoTypeName() string {
	return string(p.Type)
}

// JsTypeName returns the Js TypeName for the ParamType
func (p Param) JsTypeName() string {
	switch p.Type {
	case ParamTypeString:
		return "string"
	case ParamTypeBool:
		return "boolean"
	case ParamTypeInt, ParamTypeInt8, ParamTypeInt16, ParamTypeInt32, ParamTypeInt64,
		ParamTypeUint, ParamTypeUint8, ParamTypeUint16, ParamTypeUint32, ParamTypeUint64:
		return "number"
	default:
		panic("unknown ParamType value")
	}
}

// IsNumber returns true when the type is numeric
func (p Param) IsNumber() bool {
	switch p.Type {
	case ParamTypeInt8, ParamTypeInt16, ParamTypeInt32, ParamTypeInt64, ParamTypeInt,
		ParamTypeUint8, ParamTypeUint16, ParamTypeUint32, ParamTypeUint64, ParamTypeUint:
		return true
	default:
		return false
	}
}

// NumberMax returns the max number constant for the given type.
// Or an error when the type is not a number
func (p Param) NumberMax() (uint64, error) {
	switch p.Type {
	case ParamTypeInt8:
		return math.MaxInt8, nil
	case ParamTypeInt16:
		return math.MaxInt16, nil
	case ParamTypeInt32:
		return math.MaxInt32, nil
	case ParamTypeInt64, ParamTypeInt:
		return math.MaxInt64, nil
	case ParamTypeUint8:
		return math.MaxUint8, nil
	case ParamTypeUint16:
		return math.MaxUint16, nil
	case ParamTypeUint32:
		return math.MaxUint32, nil
	case ParamTypeUint64, ParamTypeUint:
		return math.MaxUint64, nil
	default:
		return 0, errors.New("not a number")
	}
}

// NumberMax returns the max number constant for the given type.
// Or an error when the type is not a number
func (p Param) NumberMin() (int64, error) {
	switch p.Type {
	case ParamTypeInt8:
		return math.MinInt8, nil
	case ParamTypeInt16:
		return math.MinInt16, nil
	case ParamTypeInt32:
		return math.MinInt32, nil
	case ParamTypeInt64, ParamTypeInt:
		return math.MinInt64, nil
	case ParamTypeUint8:
		return 0, nil
	case ParamTypeUint16:
		return 0, nil
	case ParamTypeUint32:
		return 0, nil
	case ParamTypeUint64, ParamTypeUint:
		return 0, nil
	default:
		return 0, errors.New("not a number")
	}
}

// Params is a list of parameters
type Params []*Param

// JsParameterList returns a comma seperated string of arguments (name only)
func (ps Params) JsParameterList() string {
	names := make([]string, 0, len(ps))
	for _, p := range ps {
		names = append(names, p.Name)
	}
	return strings.Join(names, ", ")
}

// GoParameterList returns the params as go ParameterList (http://golang.org/ref/spec#ParameterList)
func (ps Params) GoParameterList() string {
	params := make([]string, 0, len(ps))
	for _, p := range ps {
		params = append(params, p.Name+" "+p.GoTypeName())
	}
	return strings.Join(params, ", ")
}
