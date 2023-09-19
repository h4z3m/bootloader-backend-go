package main

import (
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

type Event struct {
	// Events are pushed to this channel by the main events-gathering routine
	Message chan string

	// New client connections
	NewClients chan chan string

	// Closed client connections
	ClosedClients chan chan string

	// Total client connections
	TotalClients map[chan string]bool
}

type ClientChan chan string

// listen listens for events on the stream and performs the following actions:
// 1. Adds new available client to the stream.
// 2. Removes closed client from the stream.
// 3. Broadcasts message to all clients in the stream.
//
// No parameters.
// No return type.
func (stream *Event) listen() {
	for {
		select {
		// Add new available client
		case client := <-stream.NewClients:
			stream.TotalClients[client] = true
			log.Printf("Client added. %d registered clients", len(stream.TotalClients))

			// Remove closed client
		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client)
			close(client)
			log.Printf("Removed client. %d registered clients", len(stream.TotalClients))

			// Broadcast message to client
		case eventMsg := <-stream.Message:
			for clientMessageChan := range stream.TotalClients {
				clientMessageChan <- eventMsg
			}
		}
	}
}

// headersMiddlewareSSE is a middleware function that sets the necessary headers for Server-Sent Events (SSE) response.
//
// It takes in a gin Context object as a parameter.
// It does not return any value.
func headersMiddlewareSSE() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}

// serveHTTP is a function that handles the HTTP request for the Event struct.
//
// It takes a gin.Context as a parameter and returns a gin.HandlerFunc.
func (stream *Event) serveHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize client channel
		clientChan := make(ClientChan, 1)
		// Send new connection to event server
		stream.NewClients <- clientChan
		defer func() {
			// Send closed connection to event server
			stream.ClosedClients <- clientChan
		}()

		c.Set("clientChan", clientChan)

		c.Next()
	}
}

// handleSSE handles the Server-Sent Events (SSE) for the given gin context.
//
// It expects a "clientChan" key to be set in the gin context, which should contain
// a channel of type ClientChan. If the "clientChan" key is not found or the value
// associated with it is not of type ClientChan, the function returns early.
//
// It sends the status of the client event to the clientChan channel.
//
// It then streams messages from the clientChan channel to the client using the SSE protocol.
// If a message is received from the clientChan channel, it sends the message to the client
// using the "message" event type of the SSE protocol. If the clientChan channel is closed,
// the function returns false to stop the streaming.
func handleSSE(c *gin.Context) {

	v, ok := c.Get("clientChan")
	if !ok {
		return
	}
	clientChan, ok := v.(ClientChan)
	if !ok {
		return
	}

	clientChan <- getClientEventStatus()

	c.Stream(func(w io.Writer) bool {
		// Stream message to client from message channel
		if msg, ok := <-clientChan; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})
}
