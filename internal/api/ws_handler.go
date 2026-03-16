package api

import (
	"log"
	"net/http"
	"pollstream/internal/poll"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	hub *poll.Hub
}

func NewWSHandler(hub *poll.Hub) *WSHandler {
	return &WSHandler{hub: hub}
}

func (h *WSHandler) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading connection: %v", err)
		return
	}

	h.hub.Register(conn)

	go func() {
		defer h.hub.Unregister(conn)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Printf("websocket read error: %v", err)
				break
			}
		}
	}()
}
