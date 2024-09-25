package websocket

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	ID uint
	Connection *websocket.Conn
	Pool *Pool
	Name string
}

type Message struct {
	RoomID int32 `json:"RoomID,omitempty"`
	Message string `json:"Message,omitempty"`
	Owner string `json:"Owner,omitempty"`
}

func (c *Client) InitConnection(requestBody chan []byte) {
	defer func() {
		c.Pool.Unregister <- c
		c.Connection.Close()
	}()
	defer c.Pool.ReviveWebsocket()

	c.Pool.Register <- c

	for {		
		_, p, _ := c.Connection.ReadMessage()
		
		requestBody <- p
	}
}