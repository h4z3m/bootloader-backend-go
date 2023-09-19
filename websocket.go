package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var connected bool = false

// websocketCloseHandler handles the closing of a websocket connection.
//
// It takes two parameters:
//   - code: an integer representing the reason for the closure.
//   - text: a string containing additional information about the closure.
//
// It does not return any value.
func websocketCloseHandler(code int, text string) error {
	log.Printf("Connection with client IP: %s closed\n", text)
	// Pass event to SSE that the connection has closed
	stream.Message <- "Disconnected"
	connected = false
	return nil
}

// websocketConnectHandler handles the connection with a client IP.
//
// It takes a string parameter clientIP, representing the IP address of the client.
// The function does the following:
// - Logs a message indicating that the connection with the client IP was opened.
// - Passes an event to the SSE (Server-Sent Events) that the connection was opened.
// - Sets the connected variable to true.
func websocketConnectHandler(clientIP string) {
	log.Printf("Connection with client IP: %s opened\n", clientIP)
	// Pass event to SSE that the connection was opened
	stream.Message <- "Connected"
	connected = true
}

// serveWebsocket handles the WebSocket connection upgrade and communication.
//
// It takes a http.ResponseWriter and a http.Request as parameters.
// It does not return anything.
func serveWebsocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Pass event to SSE that the connection was opened
	websocketConnectHandler(r.RemoteAddr)

	// Set the close handler if connection is lost
	ws.SetCloseHandler(websocketCloseHandler)

	defer ws.Close()
	for {

		log.Println("Waiting for command")
		// Attempt to receive a message from the channel
		var command BL_Command = <-command_channel

		commandId := command.Type()

		log.Println("Received command:", command)

		// Convert JSON to string to be sent to websocket

		data, _ := json.Marshal(command)
		log.Println(string(data))
		// Send a command to websocket
		ws.WriteMessage(websocket.TextMessage, data)

		// Wait for a response
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("Received message:", string(p))

		// Handle the response
		var response interface{}

		switch messageType {
		case websocket.TextMessage:
			response, err = responseMap[commandId](p)
			if err != nil {
				log.Println(err)
			}
		default:
			log.Println("Unknown message type:", messageType)
		}

		// Send a response to the channel
		select {
		case response_channel <- response.(BL_Response):
		default:
			log.Println("Channel full")
		}

	}

}

// websocketMain is the main function for handling websocket connections.
//
// It sets up the websocket endpoint by calling the serveWebsocket function.
// Then it starts the HTTP server to listen on port 8080 using http.ListenAndServe.
func websocketMain() {

	http.HandleFunc("/", serveWebsocket)

	http.ListenAndServe(":8080", nil)

}
