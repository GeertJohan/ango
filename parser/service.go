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

func newService() *Service {
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
	Args   []*Param
	Rets   []*Param
}

// ProcedureType indicates wether a procedure is server- or client-side
type ProcedureType string

var (
	ServerProcedure = ProcedureType("string")
	ClientProcedure = ProcedureType("client")
)

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
