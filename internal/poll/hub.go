package poll

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients map[*Client]bool
	broadcast chan interface{}
	register chan *Client
	unregister chan *Client

	mu sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan interface{}, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) Broadcast(v interface{}) {
	h.broadcast <- v
}

func (h *Hub) Register(conn *websocket.Conn) *Client {
	client := &Client{hub: h, conn: conn, send: make(chan interface{}, 256)}
	h.register <- client
	return client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}
