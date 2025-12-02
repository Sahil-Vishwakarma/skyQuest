package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	ws "github.com/skyquest/server/internal/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from our frontend
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:5173" ||
			origin == "http://localhost:3000" ||
			origin == ""
	},
}

type WebSocketHandler struct {
	hub *ws.Hub
}

func NewWebSocketHandler(hub *ws.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

// HandleWebSocket handles WebSocket upgrade and connection
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	sessionID := c.Query("sessionId")
	client := ws.NewClient(h.hub, conn, sessionID)
	h.hub.Register(client)

	// Start goroutines for reading and writing
	go client.WritePump()
	go client.ReadPump()
}

