package main

import (
	"code.google.com/p/go.net/websocket"
	"net/http"
)

// CalculatorHandlerIface should be implemented by the implementation part of the protocol
type CalculatorHandlerIface interface {
	// OnConnect is called when the connection was successfully set up
	// default method for any ango serviceHandler
	OnConnect()

	// OnDisconnect is called when the connection was closed.
	// When err is nil, the connection was closed by the client or service on purpose (call to .Close())
	// default method for any ango serviceHandler
	OnDisconnect(err error)

	// Add should add a and b and return the result
	// This is a synchronous service
	Add(a int, b int) (int, error)

	// Subtract should subtract a from b and return it
	// This is a synchronous service
	Subtract(a int, b int) (int, error)

	// Clear history
	// This is an asynchronous service
	Clear()
}

type CalculatorClient struct{}

func (cc *CalculatorClient) DisplayNotification(subject string, text string) error {
	//++ do itt!!
}

// NewCalculatorHandlerFunc is a factory function that should return a calculatorHandlerIface implementation
type NewCalculatorHandlerFunc func(*CalculatorClient, *ClientInfo) *CalculatorHandlerIface

// CalculatorServer handles new requests for websockets and creates new sessions
type CalculatorServer struct {
	newHandler NewCalculatorHandlerFunc
}

// NewCaculatorService creates and returns a new CalculatorServer instance
func NewCalculatorServer(handler NewCalculatorHandlerFunc) *CalculatorServer {
	cs := &CalculatorServer{
		newHandler: handler,
	}
	return cs
}

func (cs *CalculatorServer) ServiceHTTP(w http.ResponseWriter, r *http.Request) {
	websocket.Handler(cs.wsHandler).ServeHTTP(w, r)
}

func (cs *CalculatorServer) wsHandler(ws *websocket.Conn) {
	si := &SessionInfo{
		WsConn: ws,
	}
	handler := newHandler(si)
}
