package poll

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients    map[string]map[*Client]bool
	broadcast  chan *Poll
	register   chan *Client
	unregister chan *Client

	mu sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Poll, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.pollID] == nil {
				h.clients[client.pollID] = make(map[*Client]bool)
			}
			h.clients[client.pollID][client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.pollID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.clients, client.pollID)
					}
				}
			}
			h.mu.Unlock()
		case poll := <-h.broadcast:
			h.mu.Lock()
			if clients, ok := h.clients[poll.ID]; ok {
				for client := range clients {
					select {
					case client.send <- poll:
					default:
						close(client.send)
						delete(clients, client)
						if len(clients) == 0 {
							delete(h.clients, poll.ID)
						}
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) Broadcast(poll *Poll) {
	h.broadcast <- poll
}

func (h *Hub) Register(conn *websocket.Conn, pollID string) *Client {
	client := &Client{hub: h, conn: conn, send: make(chan interface{}, 256), pollID: pollID}
	h.register <- client
	return client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}
