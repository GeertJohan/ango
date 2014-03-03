package main

import (
	"fmt"
	"github.com/GeertJohan/go.incremental"
	"github.com/GeertJohan/go.rice"
	"net/http"
	"os"
)

// ++ This could move to implementation of ChatService.Server (type interface) see notes below.
var idInc = &incremental.Int{}

// ChatServiceSession implements ChatServiceHandler
type ChatServiceSession struct {
	name   string
	id     int
	client *ChatServiceClient
}

// NewChatServiceSession creates and returns a new ChatServiceHandler instance
func NewChatServiceSession(client *ChatServiceClient) ChatServiceSessionInterface {
	session := &ChatServiceSession{
		id:     idInc.Next(),
		client: client,
	}

	// ask users name in a goroutine to not block the session creation
	go func() {
		answer, err := client.AskQuestion("what's your name?")
		if err != nil {
			if err != ErrNotImplementedYet {
				panic("TODO: implement error handling." + err.Error())
			}
		}
		fmt.Printf("welcome %s\n", answer)
		session.name = answer
	}()

	// return session
	return session
}

func (cs *ChatServiceSession) Stop(err error) {
	fmt.Printf("Stopping session %d with error: %s\n", cs.id, err)
}

func (cs *ChatServiceSession) Add(a int, b int) (c int, err error) {
	c = a + b
	fmt.Printf("Call to Add(%d, %d) will return %d\n", a, b, c)
	cs.client.DisplayNotification("we have", fmt.Sprintf("return value %d", c))
	return c, nil
}

func (cs *ChatServiceSession) Add8(a int8, b int8) (c int16, err error) {
	c = int16(a) + int16(b)
	fmt.Printf("Call to Add8(%d, %d) will return %d\n", a, b, c)
	return c, nil
}

func (cs *ChatServiceSession) Notify(text string) {
	fmt.Printf("instance %d have notification: %s\n", cs.id, text)
}

//++ TODO: maybe create interface ChatService.Server with function NewSession.
//++ TODO: drop ErrorIncommingConnection and have global Debug implementation for debugging.
//++ TODO: Errors on incomming connections are not relevant during production runtime. (net/http doesn't expose those errors either..)
var server = &ChatServiceServer{
	NewSession: NewChatServiceSession,
	ErrorIncommingConnection: func(err error) {
		fmt.Printf("Error setting up connection: %s\n", err)
	},
}

func main() {
	httpFiles, err := rice.FindBox("http-files")
	if err != nil {
		fmt.Printf("Error opening http filex box: %s\n", err)
		os.Exit(1)
	}

	http.Handle("/", http.FileServer(httpFiles.HTTPBox()))
	http.Handle("/websocket-ango-chatService", server)

	err = http.ListenAndServe(":8123", nil)
	if err != nil {
		fmt.Printf("Error listenAndServe: %s\n", err)
		os.Exit(1)
	}
}
