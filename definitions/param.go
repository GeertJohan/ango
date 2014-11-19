package definitions

import (
	"errors"
	"math"
	"strings"
)

// Param defines an argument or return parameter
type Param struct {
	Name string
	Type *Type
}

// CapitalizedName returns the name for this param, capitalized
func (p *Param) CapitalizedName() string {
	return strings.ToUpper(p.Name[:1]) + p.Name[1:]
}

// GoTypeName returns the type
func (p *Param) GoTypeName() string {
	switch p.Type.Category {
	case Builtin, Simple, Slice, Map:
		return p.Type.GoName()
	case Struct:
		return `*` + p.Type.CapitalizedName()
	default:
		panic("unknown type category")
	}
}

// TODO: move method to (t *Type)
// JsTypeName returns the Js TypeName for the Type
func (p *Param) JsTypeName() string {
	switch p.Type.Name {
	case TypeString.Name:
		return "string"
	case TypeBool.Name:
		return "boolean"
	case TypeInt.Name, TypeInt8.Name, TypeInt16.Name, TypeInt32.Name, TypeInt64.Name,
		TypeUint.Name, TypeUint8.Name, TypeUint16.Name, TypeUint32.Name, TypeUint64.Name:
		return "number"
	default:
		return "custom" + p.Type.CapitalizedName()
	}
}

// TODO: move method to (t *Type)
// IsNumber returns true when the type is numeric
func (p *Param) IsNumber() bool {
	switch p.Type.Name {
	case TypeInt.Name, TypeInt8.Name, TypeInt16.Name, TypeInt32.Name, TypeInt64.Name,
		TypeUint.Name, TypeUint8.Name, TypeUint16.Name, TypeUint32.Name, TypeUint64.Name:
		return true
	default:
		return false
	}
}

// TODO: move method to (t *Type)
// NumberMax returns the maximal numeric value for the given type or an error when the type is not a number
func (p Param) NumberMax() (uint64, error) {
	switch p.Type.Name {
	case TypeInt8.Name:
		return math.MaxInt8, nil
	case TypeInt16.Name:
		return math.MaxInt16, nil
	case TypeInt32.Name:
		return math.MaxInt32, nil
	case TypeInt64.Name, TypeInt.Name: // TODO: Int always Int64 ??
		return math.MaxInt64, nil
	case TypeUint8.Name:
		return math.MaxUint8, nil
	case TypeUint16.Name:
		return math.MaxUint16, nil
	case TypeUint32.Name:
		return math.MaxUint32, nil
	case TypeUint64.Name, TypeUint.Name: // TODO: Uint always Uint64 ??
		return math.MaxUint64, nil
	default:
		return 0, errors.New("not a number")
	}
}

// TODO: move method to (t *Type)
// NumberMin returns the minimal numeric value for the given type or an error when the type is not a number
func (p Param) NumberMin() (int64, error) {
	switch p.Type.Name {
	case TypeInt8.Name:
		return math.MinInt8, nil
	case TypeInt16.Name:
		return math.MinInt16, nil
	case TypeInt32.Name:
		return math.MinInt32, nil
	case TypeInt64.Name, TypeInt.Name: // TODO: Int always Int64 ??
		return math.MinInt64, nil
	case TypeUint8.Name:
		return 0, nil
	case TypeUint16.Name:
		return 0, nil
	case TypeUint32.Name:
		return 0, nil
	case TypeUint64.Name, TypeUint.Name:
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
