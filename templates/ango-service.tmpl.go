// WARNING This is generated code by the ango tool (github.com/GeertJohan/ango)
// DO NOT EDIT unless you know what you're doing!

package {{.PackageName}}

import (
	"errors"
	"fmt"
	"encoding/json"
	"net/http"

	"github.com/GeertJohan/go.wstext"
	"github.com/GeertJohan/go.incremental"
	"github.com/gorilla/websocket"
)

const ProtocolVersion = "{{.ProtocolVersion}}"

var (
	ErrInvalidVersionString = errors.New("invalid version string")
	ErrInvalidMessageType   = errors.New("invalid message type")
	ErrUnkonwnProcedure     = errors.New("unknown procedure")
	ErrNotImplementedYet    = errors.New("not implemented yet")
	ErrInvalidCallbackID    = errors.New("callbackID is inavlid")
)

const (
	msgTypeRequest  = "req"
	msgTypeResponse = "res"
)

// root structure for incoming message json
type angoInMsg struct {
	Type       string          `json:"type"`      // "req" or "res"
	Procedure  string          `json:"procedure"` // name for the procedure 9when "req"
	CallbackID uint64          `json:"cb_id"`     // callback ID for request or response
	Data       json.RawMessage `json:"data"`      // remain raw, depends on procedure
	Error      json.RawMessage `json:"error"`     // remain raw, depens on ??
}

// root structure for outgoing message json
type angoOutMsg struct {
	Type       string        `json:"type"`                // "req" or "res"
	Procedure  string        `json:"procedure,omitempty"` // name for the procedure 9when "req"
	CallbackID uint64        `json:"cb_id,omitempty"`     // callback ID for request or response
	Data       interface{}   `json:"data,omitempty"`      // remain raw, depends on procedure
	Error      *angoOutError `json:"error,omitempty"`     // when not-nil, an error ocurred
}

type angoOutError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

{{range .Service.ServerProcedures}}
	type angoServerArgsData{{.CapitalizedName}} struct {
		{{range .Args}}
			{{.CapitalizedName}} {{.GoTypeName}} `json:"{{.Name}}"` {{end}}
	}
	{{if not .Oneway}}
		type angoServerRetsData{{.CapitalizedName}} struct {
			{{range .Rets}}
				{{.CapitalizedName}} {{.GoTypeName}} `json:"{{.Name}}"` {{end}}
		}
	{{end}}
{{end}}


{{range .Service.ClientProcedures}}
	type angoClientArgsData{{.CapitalizedName}} struct {
		{{range .Args}}
			{{.CapitalizedName}} {{.GoTypeName}} `json:"{{.Name}}"` {{end}}
	}
	{{if not .Oneway}}
		type angoClientRetsData{{.CapitalizedName}} struct {
			{{range .Rets}}
				{{.CapitalizedName}} {{.GoTypeName}} `json:"{{.Name}}"` {{end}}
		}
	{{end}}
{{end}}

// SessionHandler defines all methods that can be called by the client
type SessionHandler interface {
	// Stop is called when the session is closing (websocket closed)
	Stop(err error)

	{{range .Service.ServerProcedures}}
		// {{.CapitalizedName}} is a ango procedure defined in the .ango file
		{{.CapitalizedName}}( {{.GoArgs}} )( {{.GoRets}} )
	{{end}}
}

// // NewSessionFunc must return a new instance implementing {{.Service.CapitalizedName}}SessionHandler
// //++ TODO: rename to NewSessionFunc (?) when generated code gets it's own package
// type NewSessionFunc func(*Client)(handler SessionHandler)

// Server handles incomming http requests
type Server struct {
	NewSession               func(*Client)(handler SessionHandler)
	ErrorIncommingConnection func(err error)
}

// ServeHTTP hijacks incomming http connections and sets up the websocket communication
func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	// wrap for simple text read/write
	textconn := wstext.Conn{conn}

	receivedVersion, err := textconn.ReadText()
	if err != nil {
		if server.ErrorIncommingConnection != nil {
			server.ErrorIncommingConnection(err)
		}
		return
	}
	if receivedVersion != ProtocolVersion {
		_ = textconn.WriteText("invalid")
		fmt.Printf("err: %s\n", err)
		fmt.Printf("in: '%s'\n", receivedVersion)
		fmt.Printf("hv: '%s'\n", ProtocolVersion)
		if server.ErrorIncommingConnection != nil {
			server.ErrorIncommingConnection(ErrInvalidVersionString)
		}
		return
	}
	err = textconn.WriteText("good")
	if err != nil {
		if server.ErrorIncommingConnection != nil {
			server.ErrorIncommingConnection(err)
		}
		return
	}

	fmt.Println("Valid protocol version detected")

	// create new client instance with conn
	client := &Client{
		ws:               conn,
		callbackInc:      &incremental.Uint64{},
		callbackChannels: make(map[uint64]chan *angoInMsg),
	}

	// create session on server
	session := server.NewSession(client)
	
	// run protocol
	err = run{{.Service.CapitalizedName}}Protocol(conn, client, session)
	// err can be nil, but we want to call .Stop always
	session.Stop(err)
}

func run{{.Service.CapitalizedName}}Protocol(conn *websocket.Conn, client *Client, session SessionHandler) error {
	for {
		{{/* unmarshal root message structure */}}
		inMsg := &angoInMsg{}
		err := conn.ReadJSON(inMsg)
		if err != nil {
			return err
		}

		switch inMsg.Type {
		case msgTypeRequest:
			fmt.Printf("Have request: %s\n", inMsg.Procedure)
			switch inMsg.Procedure {
			{{range .Service.ServerProcedures}}
				case "{{.Name}}":
					{{/* unmarshal procedure arguments */}}
					procArgs := &angoServerArgsData{{.CapitalizedName}}{} {{/* var procArgs is referenced by .GoCallArgs */}}
					err = json.Unmarshal(inMsg.Data, procArgs)
					if err != nil {
						return err
					}

					{{/* prepare for return values */}}
					{{if not .Oneway}}
						procRets := &angoServerRetsData{{.CapitalizedName}}{} {{/* var procRets is referenced by .GoCallRets */}}
						var procErr error {{/* var procErr is referenced by .GoCallRets */}}
					{{end}}

					{{/* call procedure, accept return values when not oneway */}}
					{{if not .Oneway}}{{.GoCallRets}} = {{end}}session.{{.CapitalizedName}}( {{.GoCallArgs}} )

					{{/* return message with procedure return values */}}
					{{if not .Oneway}}
						outMsg := &angoOutMsg{
							Type:       "res",
							CallbackID: inMsg.CallbackID,
						}
						if procErr != nil {
							outMsg.Error = &angoOutError{
								Type: "errorReturned",
								Message: procErr.Error(),
							}
							err = conn.WriteJSON(outMsg)
							if err != nil {
								return err
							}
							break
						}
						outMsg.Data = procRets
						err = conn.WriteJSON(outMsg)
						if err != nil {
							return err
						}
					{{end}}
			{{end}}
			default:
				return ErrUnkonwnProcedure
			}
		case msgTypeResponse:
			callbackCh := client.callbackChannels[inMsg.CallbackID]
			if callbackCh == nil {
				return ErrInvalidCallbackID
			}
			delete(client.callbackChannels, inMsg.CallbackID)

			callbackCh <- inMsg
		default:
			return ErrInvalidMessageType
		}
	}
}

// Client is a reference to the client end-point and available methods defined on the client
type Client struct {
	ws               *websocket.Conn
	callbackInc      *incremental.Uint64
	callbackChannels map[uint64]chan *angoInMsg
}

{{range .Service.ClientProcedures}}
	// {{.CapitalizedName}} is a ango procedure defined in the .ango file
	{{if .Oneway}}
		func (c *Client) {{.CapitalizedName}}( {{.GoArgs}} )( err error ) {
			fmt.Println("Called oneway service {{.CapitalizedName}}")
			outMsg := angoOutMsg{
				Type:      "req",
				Procedure: "{{.Name}}",
				Data:      &angoClientArgsData{{.CapitalizedName}}{
				{{range .Args}}
					{{.CapitalizedName}}: {{.Name}},{{end}}
				},
			}

			// write message
			err = c.ws.WriteJSON(outMsg)
			if err != nil {
				return {{/* when service is not oneway, this will return the error using named return values */}}
			}

			// all done
			return
		}
	{{else}}
		type {{.CapitalizedName}}Result struct {
			// error
			Err error

			{{range .Rets}}
				{{.CapitalizedName}} {{.Type}}{{end}}
		}
		func (c *Client) {{.CapitalizedName}}( {{.GoArgs}} )( retCh <-chan *{{.CapitalizedName}}Result ) {
			ch := make(chan *{{.CapitalizedName}}Result, 1)
			retCh = ch
			response := &{{.CapitalizedName}}Result{}
			fmt.Println("Called returning service {{.CapitalizedName}}")
			
			go func() {
				defer func() {
					ch <- response
				}()

				outMsg := angoOutMsg{
					Type:      "req",
					Procedure: "{{.Name}}",
					Data:      &angoClientArgsData{{.CapitalizedName}}{
					{{range .Args}}
						{{.CapitalizedName}}: {{.Name}},{{end}}
					},
				}

				// create callback channel
				callbackCh := make(chan *angoInMsg, 1)
				outMsg.CallbackID = c.callbackInc.Next()
				c.callbackChannels[outMsg.CallbackID] = callbackCh

				// write message
				response.Err = c.ws.WriteJSON(outMsg)
				if response.Err != nil {
					return {{/* when service is not oneway, this will return the error using named return values */}}
				}

				// wait for response message
				respMsg := <- callbackCh
				close(callbackCh)

				// check for error
				if(respMsg.Error != nil) {
					var errStr string
					response.Err = json.Unmarshal(respMsg.Error, &errStr)
					if response.Err != nil {
						return
					}
					response.Err = errors.New(errStr)
					return
				}
				retsData := &angoClientRetsData{{.CapitalizedName}}{}
				response.Err = json.Unmarshal(respMsg.Data, retsData)
				if response.Err != nil {
					return
				}
				{{range .Rets}}
					response.{{.CapitalizedName}} = retsData.{{.CapitalizedName}}{{end}}
			}()

			return
		}
	{{end}}
{{end}}