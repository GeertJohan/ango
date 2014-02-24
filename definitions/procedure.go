package definitions

import (
	"strings"
)

// ProcedureType indicates wether a procedure is server- or client-side
type ProcedureType string

// ProcedureType's
var (
	ServerProcedure = ProcedureType("string")
	ClientProcedure = ProcedureType("client")
)

// Procedure defines a remote method/function
type Procedure struct {
	Type   ProcedureType
	Oneway bool
	Name   string
	Args   Params
	Rets   Params
}

// CapitalizedName returns the name, capitalized.
// Used by ango-service.tmpl.go
func (p *Procedure) CapitalizedName() string {
	return strings.ToUpper(p.Name[:1]) + p.Name[1:]
}

// GoArgs returns the go arguments
// Used by ango-service.tmpl.go
func (p *Procedure) GoArgs() string {
	return p.Args.GoParameterList()
}

// GoRets returns the go return values
// Used by ango-service.tmpl.go
func (p *Procedure) GoRets() string {
	retStr := p.Rets.GoParameterList()
	if !p.Oneway {
		if len(retStr) > 0 {
			retStr += ","
		}
		retStr += "err error"
	}
	return retStr
}
