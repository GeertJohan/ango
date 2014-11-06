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
	Source Source
}

// CapitalizedName returns the name, capitalized.
// Used by ango-service.tmpl.go
func (p *Procedure) CapitalizedName() string {
	return strings.ToUpper(p.Name[:1]) + p.Name[1:]
}

// JsArgs returns the js arguments
// Used by ango-service.tmpl.js
func (p *Procedure) JsArgs() string {
	return p.Args.JsParameterList()
}

// GoArgs returns the go function definition argument ParameterList
// Used by ango-service.tmpl.go
func (p *Procedure) GoArgs() string {
	return p.Args.GoParameterList()
}

// GoRets returns the go function definition return ParameterList
// Used by ango-service.tmpl.go
func (p *Procedure) GoRets() string {
	retStr := p.Rets.GoParameterList()
	if !p.Oneway {
		if len(retStr) > 0 {
			retStr += ", "
		}
		retStr += "err error"
	}
	return retStr
}

// JsCallArgs returns the procedure call argument values
// Used by ango-service.tmpl.js
func (p *Procedure) JsCallArgs() string {
	str := ""
	for _, param := range p.Args {
		if len(str) > 0 {
			str += ", "
		}
		str += "messageObj.data." + param.Name
	}
	return str
}

// GoCallArgs returns the procedure call argument values
// Used by ango-service.tmpl.go
func (p *Procedure) GoCallArgs() string {
	str := ""
	for _, param := range p.Args {
		if len(str) > 0 {
			str += ", "
		}
		str += "procArgs." + param.CapitalizedName()
	}
	return str
}

// GoCallRets returns the procedure call return values
// Used by ango-service.tmpl.go
func (p *Procedure) GoCallRets() string {
	str := ""
	for _, param := range p.Rets {
		if len(str) > 0 {
			str += ", "
		}
		str += "procRets." + param.CapitalizedName()
	}
	if !p.Oneway {
		if len(str) > 0 {
			str += ","
		}
		str += "procErr"
	}
	return str
}
