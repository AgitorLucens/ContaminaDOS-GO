package ws

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"sync"
	"github.com/gin-gonic/gin"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Active WebSocket connections
var (
	connections = make(map[*websocket.Conn]bool)
	mutex       = sync.Mutex{}
)

// @title My API
// @version 1.0
// @description API Documentation including WebSocket
// @host localhost:8080
// @BasePath /

// @Summary WebSocket Connection
// @Description Establish a WebSocket connection to send and receive messages.
// @Tags WebSocket
// @Produce json
// @Success 101 {string} string "Switching Protocols"
// @Router /ws [get]
func HandleWebSocket(c *gin.Context) {
	// Upgrade HTTP request to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade:", err)
		return
	}
	defer conn.Close()

	// Add connection to active list
	mutex.Lock()
	connections[conn] = true
	mutex.Unlock()
	log.Println("New client connected")

	// Read messages from client
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Client disconnected:", err)
			mutex.Lock()
			delete(connections, conn)
			mutex.Unlock()
			break
		}

		log.Println("Received:", string(msg))

		// Broadcast message to all clients
		broadcast(msg)
	}
}

// Broadcast message to all active connections
func broadcast(message []byte) {
	mutex.Lock()
	defer mutex.Unlock()

	for conn := range connections {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Failed to send message:", err)
			conn.Close()
			delete(connections, conn)
		}
	}
}