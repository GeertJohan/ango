package definitions

import (
	"strings"
)

// Param defines an argument or return parameter
type Param struct {
	Name string
	Type ParamType
}

// ParamType is the type of a parameter
type ParamType string

// ParamType's
var (
	ParamTypeInt    = ParamType("int")
	ParamTypeUint   = ParamType("uint")
	ParamTypeString = ParamType("string")
)

// Params is a list of parameters
type Params []*Param

// CommaSeperatedString returns a comma seperated string of arguments
func (ps Params) CommaSeperatedString() string {
	names := make([]string, 0, len(ps))
	for _, p := range ps {
		names = append(names, p.Name)
	}
	return strings.Join(names, ", ")
}
