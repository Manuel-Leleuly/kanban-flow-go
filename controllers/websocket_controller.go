package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{os.Getenv("BASE_URL")}
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				return true
			}
		}
		return false
	},
}

var clients = make(map[*websocket.Conn]bool)

func WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	clients[conn] = true

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			delete(clients, conn)
			return
		}
	}
}

// Broadcast message to all clients
func BroadcastMessage(message []byte) {
	for client := range clients {
		client.WriteMessage(websocket.TextMessage, message)
	}
}
