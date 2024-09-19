package main

import (
	"fmt"
	"log"
	"time"

	ws "github.com/anicoll/evtwebsocket"
)

func main() {
	c := ws.New(
		ws.OnConnected(func(w ws.Connection) {
			log.Println("Connected")
		}),
		ws.OnMessage(func(msg []byte, w ws.Connection) {
			log.Printf("OnMessage: %s\n", msg)
		}),
		// When the client disconnects for any reason
		ws.OnError(func(err error) {
			log.Printf("** ERROR **\n%s\n", err.Error())
		}),
		// This is used to match the request and response messagesP>termina
		ws.WithMatchMsg(func(req, resp []byte) bool {
			return string(req) == string(resp)
		}),
		// 	// Auto reconnect on error
		ws.WithReconnect(true),
		ws.WithPingIntervalSec(5),
		ws.WithPingMsg([]byte("PING")),
	)

	// Connect
	if err := c.Dial("ws://echo.websocket.org", ""); err != nil {
		log.Fatal(err)
	}

	for i := 1; i <= 100; i++ {

		// Create the message with a callback
		msg := ws.Msg{
			Body: []byte(fmt.Sprintf("Hello %d", i)),
			Callback: func(resp []byte, w ws.Connection) {
				log.Printf("[%d] Callback: %s\n", i, resp)
			},
		}

		log.Printf("[%d] Sending message: %s\n", i, msg.Body)

		// Send the message to the server
		if err := c.Send(msg); err != nil {
			log.Println("Unable to send: ", err.Error())
		}

		// Take a break
		time.Sleep(time.Second * 2)
	}

}
