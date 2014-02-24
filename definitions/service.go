package definitions

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
