package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/skyquest/server/internal/models"
)

// Client represents a WebSocket client connection
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	sessionID string
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's event loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()
			log.Printf("Client connected. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mutex.Unlock()
			log.Printf("Client disconnected. Total clients: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

// Register adds a client to the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister removes a client from the hub
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// BroadcastFlights sends flight updates to all connected clients
func (h *Hub) BroadcastFlights(flights []models.Flight) {
	msg := models.WSMessage{
		Type: "flight:update",
		Payload: models.WSFlightUpdate{
			Flights: flights,
		},
	}
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling flight update: %v", err)
		return
	}
	h.broadcast <- data
}

// SendToClient sends a message to a specific client by session ID
func (h *Hub) SendToClient(sessionID string, msg models.WSMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.clients {
		if client.sessionID == sessionID {
			select {
			case client.send <- data:
			default:
				close(client.send)
				delete(h.clients, client)
			}
			break
		}
	}
}

// SendRoundStart notifies a client that a new round has started
func (h *Hub) SendRoundStart(sessionID string, roundNumber int, flight *models.Flight) {
	h.SendToClient(sessionID, models.WSMessage{
		Type: "round:start",
		Payload: models.WSRoundStart{
			SessionID:   sessionID,
			RoundNumber: roundNumber,
			Flight:      flight,
		},
	})
}

// SendGuessResult sends the result of a guess to the client
func (h *Hub) SendGuessResult(sessionID string, roundNumber int, score models.ScoreResult, totalScore int) {
	h.SendToClient(sessionID, models.WSMessage{
		Type: "guess:result",
		Payload: models.WSGuessResult{
			SessionID:   sessionID,
			RoundNumber: roundNumber,
			Score:       score,
			TotalScore:  totalScore,
		},
	})
}

// SendGameEnd notifies a client that their game has ended
func (h *Hub) SendGameEnd(sessionID string, totalScore int, rank int) {
	h.SendToClient(sessionID, models.WSMessage{
		Type: "game:end",
		Payload: models.WSGameEnd{
			SessionID:  sessionID,
			TotalScore: totalScore,
			Rank:       rank,
		},
	})
}

// NewClient creates a new client
func NewClient(hub *Hub, conn *websocket.Conn, sessionID string) *Client {
	return &Client{
		hub:       hub,
		conn:      conn,
		send:      make(chan []byte, 256),
		sessionID: sessionID,
	}
}

// GetConn returns the WebSocket connection
func (c *Client) GetConn() *websocket.Conn {
	return c.conn
}

// GetSendChannel returns the send channel
func (c *Client) GetSendChannel() chan []byte {
	return c.send
}

// SetSessionID sets the session ID for the client
func (c *Client) SetSessionID(sessionID string) {
	c.sessionID = sessionID
}

// ReadPump pumps messages from the WebSocket connection to the hub
func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle incoming messages (e.g., session registration)
		var msg models.WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case "register":
			if payload, ok := msg.Payload.(map[string]interface{}); ok {
				if sessionID, ok := payload["sessionId"].(string); ok {
					c.sessionID = sessionID
				}
			}
		}
	}
}

// WritePump pumps messages from the hub to the WebSocket connection
func (c *Client) WritePump() {
	defer func() {
		c.conn.Close()
	}()

	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

