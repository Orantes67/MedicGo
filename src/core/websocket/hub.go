package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Client representa una conexión WebSocket activa.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// Hub gestiona el conjunto de clientes conectados y difunde mensajes.
// Implementa events.EventPublisher.
type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]bool
}

// NewHub crea e inicializa un Hub listo para usar.
func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
	}
}

// Register agrega un cliente al Hub y arranca su goroutine de escritura.
func (h *Hub) Register(conn *websocket.Conn) {
	client := &Client{
		hub:  h,
		conn: conn,
		send: make(chan []byte, 256),
	}

	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	// Goroutine de escritura: envía mensajes al cliente.
	go client.writePump()
	// Goroutine de lectura: detecta cierre de conexión.
	go client.readPump()
}

// Publish serializa el evento y lo difunde a todos los clientes conectados.
// Implementa events.EventPublisher.
func (h *Hub) Publish(event interface{}) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		select {
		case client.send <- payload:
		default:
			// Canal lleno: descartamos este cliente para no bloquear.
			log.Printf("[WS] canal lleno para cliente %v, desconectando", client.conn.RemoteAddr())
			close(client.send)
			delete(h.clients, client)
		}
	}
	return nil
}

// unregister elimina un cliente del Hub.
func (h *Hub) unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[c]; ok {
		delete(h.clients, c)
		close(c.send)
	}
}

// ────────────────────────────────────────────────────────────────────────────────
// Pumps del cliente
// ────────────────────────────────────────────────────────────────────────────────

// writePump reenvía al cliente todo lo que llega por el canal send.
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("[WS] error escribiendo a %v: %v", c.conn.RemoteAddr(), err)
			break
		}
	}
}

// readPump consume los mensajes entrantes; cuando la conexión se cierra
// desregistra el cliente.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister(c)
		c.conn.Close()
	}()

	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			break
		}
	}
}
