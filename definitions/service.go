package definitions

import (
	"strings"
)

// Service holds information about a ango service
type Service struct {
	// Name is the name given to the service
	Name string

	// Types defined on the service
	Types map[string]*Type

	// ServiceProceduers holds all server-side procedures, by their name
	ServerProcedures map[string]*Procedure

	// ClientProcedures holds all client-side procedures, by their name
	ClientProcedures map[string]*Procedure
}

// NewService creates a new service instance and sets up maps and defaults
func NewService() *Service {
	s := &Service{
		Types:            make(map[string]*Type),
		ServerProcedures: make(map[string]*Procedure),
		ClientProcedures: make(map[string]*Procedure),
	}
	// s.Types[TypeInt.Name] = TypeInt
	// s.Types[TypeInt8.Name] = TypeInt8
	// s.Types[TypeInt16.Name] = TypeInt16
	// s.Types[TypeInt32.Name] = TypeInt32
	// s.Types[TypeInt64.Name] = TypeInt64
	// s.Types[TypeUint.Name] = TypeUint
	// s.Types[TypeUint8.Name] = TypeUint8
	// s.Types[TypeUint16.Name] = TypeUint16
	// s.Types[TypeUint32.Name] = TypeUint32
	// s.Types[TypeUint64.Name] = TypeUint64
	// s.Types[TypeString.Name] = TypeString
	// s.Types[TypeBool.Name] = TypeBool

	return s
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

// LookupType searches for a type in the service.Types map or BuiltinTypes map.
// When a type is builtin and is not in service.Types yet, it is added.
// When a type cannot be found, nil is returned.
func (s *Service) LookupType(name string) *Type {
	var t *Type
	if t = s.Types[name]; t != nil {
		return t
	}
	if t = BuiltinTypes[name]; t != nil {
		s.Types[t.Name] = t
		return t
	}
	return nil
}
