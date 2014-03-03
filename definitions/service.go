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

// CapitalizedName returns the service name, capitalized.
// Used by ango-service.tmpl.go
func (s *Service) CapitalizedName() string {
	return strings.ToUpper(s.Name[:1]) + s.Name[1:]
}

// JsClientProceduresStringAry returns a comma seperated list of js strings
// Used by ango-service.tmpl.js
func (s *Service) JsClientProceduresStringAry() string {
	strs := make([]string, 0, len(s.ClientProcedures))
	for name := range s.ClientProcedures {
		strs = append(strs, `'`+name+`'`)
	}
	return strings.Join(strs, ", ")
}
