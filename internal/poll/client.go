package poll

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan interface{}
}

func (c *Client) Send(v interface{}) {
	c.send <- v
}

func (c *Client) WritePump() {

	for message := range c.send {
		err := c.conn.WriteJSON(message)
		if err != nil {
			log.Println("write error:", err)
			return
		}
	}

}

func (c *Client) ReadPump() {

	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

	}

}
