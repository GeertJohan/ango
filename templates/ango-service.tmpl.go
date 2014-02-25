package {{.PackageName}}

import (
	"github.com/gorilla/websocket"
	"net/http"
)

//++ create interface that implements procedures

// {{.Service.CapitalizedName}}SessionInterface types all methods that can be called by the client
type {{.Service.CapitalizedName}}SessionInterface interface {
	{{range .Service.ServerProcedures}}
		{{.CapitalizedName}} ( {{.GoArgs}} )( {{.GoRets}} )
	{{end}}
}

// New{{.Service.CapitalizedName}}SessionInterface must return a new instance implementing {{.Service.CapitalizedName}}SessionInterface
type New{{.Service.CapitalizedName}}SessionInterface func()(handler {{.Service.CapitalizedName}}SessionInterface)

// {{.Service.CapitalizedName}}Server handles incomming http requests
type {{.Service.CapitalizedName}}Server struct {
	NewSession               New{{.Service.CapitalizedName}}SessionInterface //++ inline type?
	ErrorIncommingConnection func(err error)
}

//++ TODO: what to do with errors?
//++ add fields to Server? ErrorIncommingConnection(err error)
//++ when error occurs and non-nil: call the function with the error

// ServeHTTP hijacks incomming http connections and sets up the websocket communication
func (server *{{.Service.CapitalizedName}}Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); ok {
			http.Error(w, "Not a websocket handshake", 400)
			return
		}
		if server.ErrorIncommingConnection != nil {
			server.ErrorIncommingConnection(err)
		}
		return
	}

	//++ use conn
}