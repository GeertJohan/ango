package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/GeertJohan/ango/example/chatservice"

	"github.com/GeertJohan/go.rice"
)

var idInc = &incremental.Int{}

// ChatServiceSession implements chatservice.SessionHandler
type ChatServiceSession struct {
	name   string
	id     int
	client *chatservice.Client
}

// NewChatServiceSession creates and returns a new ServiceHandler instance
func NewChatServiceSession(client *chatservice.Client) chatservice.Session {
	session := &ChatServiceSession{
		id:     idInc.Next(),
		client: client,
	}

	// ask users name in a goroutine to not block the session creation
	go func() {
		result := <-client.AskQuestion("what's your name?")
		if result.Err != nil {
			if result.Err != chatservice.ErrNotImplementedYet {
				panic("TODO: implement error handling." + result.Err.Error())
			}
		}
		fmt.Printf("welcome %s\n", result.Answer)
		session.name = result.Answer
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
	cs.client.DisplayNotification("The server did some calculations..", fmt.Sprintf("We would like to inform you that %d+%d equals %d", a, b, c))
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

//++ TODO: maybe drop ErrorIncommingConnection and have global Debug implementation for debugging. (thoughts.md > hooks),
//			Errors on incomming connections are not relevant during production runtime. (net/http doesn't expose those errors either..)
var server = &chatservice.Server{
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
	http.Handle("/websocket-ango-chatservice", server)

	err = http.ListenAndServe(":8123", nil)
	if err != nil {
		fmt.Printf("Error listenAndServe: %s\n", err)
		os.Exit(1)
	}
}
