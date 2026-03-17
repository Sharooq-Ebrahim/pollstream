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
	hub     *poll.Hub
	service *poll.PollService
}

func NewWSHandler(hub *poll.Hub, service *poll.PollService) *WSHandler {
	return &WSHandler{hub: hub, service: service}
}

func (h *WSHandler) HandleWS(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading connection: %v", err)
		return
	}

	client := h.hub.Register(conn)

	if id != "" {
		poll, err := h.service.GetPollByID(id)
		if err == nil {
			client.Send(poll)
		}
	}

	go client.ReadPump()
	go client.WritePump()
}
