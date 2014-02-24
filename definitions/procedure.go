package definitions

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
