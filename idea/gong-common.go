package main

import (
	"code.google.com/p/go.net/websocket"
)

// ClientInfo holds generic information about a client connection
type ClientInfo struct {
	WsConn *websocket.Conn
}
