package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewSSEServer creates a new SSE (Server-Sent Events) server.
//
// It returns a pointer to an Event struct.
func NewSSEServer() (event *Event) {
	event = &Event{
		Message:       make(chan string),
		NewClients:    make(chan chan string),
		ClosedClients: make(chan chan string),
		TotalClients:  make(map[chan string]bool),
	}

	go event.listen()

	return
}

var stream *Event

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Create the websocket server
	go websocketMain()

	router := gin.Default()
	stream = NewSSEServer()

	router.Use(cors.Default())

	router.GET("/events", headersMiddlewareSSE(), stream.serveHTTP(),
		handleSSE)

	router.GET("/bl/version", readVersion)
	router.GET("/bl/status", getStatus)

	router.POST("/bl/read", flashRead)
	router.POST("/bl/flash", flashWrite)
	router.POST("/bl/erase", flashErase)
	router.POST("/bl/jump", jumpToApp)

	router.Run("localhost:3000")
}
