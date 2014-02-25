package main

import (
	"fmt"
	"github.com/GeertJohan/go.incremental"
	"github.com/GeertJohan/go.rice"
	"net/http"
	"os"
)

var idInc = &incremental.Int{}

// Handler implements ChatServiceHandler
type ChatServiceSession struct {
	id int
}

// NewChatService creates and returns a new ChatServiceHandler instance
func NewChatServiceSession() ChatServiceSessionInterface {
	// Create new ChatService instance with next id
	return &ChatServiceSession{
		id: idInc.Next(),
	}
}

func (cs *ChatServiceSession) Add(a int, b int) (c int, err error) {
	return a + b, nil
}

func (cs *ChatServiceSession) Notify(text string) {
	fmt.Printf("instance %d have notification: %s\n", cs.id, text)
}

var server = &ChatServiceServer{
	NewSession: NewChatServiceSession,
}

func main() {
	httpFiles, err := rice.FindBox("http-files")
	if err != nil {
		fmt.Printf("Error opening http filex box: %s\n", err)
		os.Exit(1)
	}

	http.Handle("/", http.FileServer(httpFiles.HTTPBox()))
	http.Handle("/ango-websocket-chatService", server)

	err = http.ListenAndServe(":8123", nil)
	if err != nil {
		fmt.Printf("Error listenAndServe: %s\n", err)
		os.Exit(1)
	}
}
