package definitions

import (
	"strings"
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

// ParamType is the type of a parameter
type ParamType string

// ParamType's
var (
	ParamTypeInt    = ParamType("int")
	ParamTypeUint   = ParamType("uint")
	ParamTypeString = ParamType("string")
)

// GoTypeName returns the Go TypeName (http://golang.org/ref/spec#TypeName) for the ParamType
func (pt ParamType) GoTypeName() string {
	return string(pt)
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
	str := ""
	for _, param := range ps {
		if len(str) > 0 {
			str += ","
		}
		str += param.Name + " " + param.Type.GoTypeName()
	}
	return str
}
