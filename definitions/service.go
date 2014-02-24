package definitions

import (
	"strings"
)

// Service holds information about a ango service
type Service struct {
	// Name is the name given to the service
	Name string

	// ServiceProceduers holds all server-side procedures, by their name
	ServerProcedures map[string]*Procedure

	// ClientProcedures holds all client-side procedures, by their name
	ClientProcedures map[string]*Procedure
}

// NewService creates a new service instance and sets up maps and defaults
func NewService() *Service {
	return &Service{
		ServerProcedures: make(map[string]*Procedure),
		ClientProcedures: make(map[string]*Procedure),
	}
}

// Procedure defines a remote method/function
type Procedure struct {
	Type   ProcedureType
	Oneway bool
	Name   string
	Args   Params
	Rets   Params
}

// ProcedureType indicates wether a procedure is server- or client-side
type ProcedureType string

// ProcedureType's
var (
	ServerProcedure = ProcedureType("string")
	ClientProcedure = ProcedureType("client")
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
