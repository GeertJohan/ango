package parser

// Service holds information about a ango service
type Service struct {
	// Name is the name given to the service
	Name string

	// ServiceProceduers holds all server-side procedures, by their name
	ServerProcedures map[string]*Procedure

	// ClientProcedures holds all client-side procedures, by their name
	ClientProcedures map[string]*Procedure
}

// Procedure defines a remote method/function
type Procedure struct {
	Oneway bool
	Name   string
	Args   []*Param
	Rets   []*Param
}

// Param defines an argument or return parameter
type Param struct {
	Name string
	Type ParamType
}

// ParamType is the type of a parameter
type ParamType string

var (
	ParamTypeInt    = ParamType("int")
	ParamTypeUint   = ParamType("uint")
	ParamTypeString = ParamType("string")
)
