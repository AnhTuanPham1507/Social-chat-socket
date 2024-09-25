package websocket

import (
	"log"
)

type Pool struct {
	Register chan *Client
	Unregister chan *Client
	Clients map[*Client]bool
	Broadcast chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register: make(chan *Client),
		Unregister: make(chan *Client),
		Clients: make(map[*Client]bool),
		Broadcast: make(chan Message),
	}
}

func (p *Pool) Start() {
	defer p.ReviveWebsocket()
	for {
		select {
			case client := <-p.Register:
				p.Clients[client] = true
				log.Println("Info:", "new client. Size of connection pool:", len(p.Clients))
				for c := range p.Clients {
					c.Connection.WriteJSON(Message{Message: "new user joined..."})
				}

			case client := <-p.Unregister:
				delete(p.Clients, client)
				log.Println("Info:", "disconnected a client. size of connection pool:", len(p.Clients))
				for c := range p.Clients {
					c.Connection.WriteJSON(Message{Message: "user disconnected..."})
				}

			case msg := <-p.Broadcast:
				log.Println("Info:", "broadcast message to clients in pool")
				for c := range p.Clients {
					c.Connection.WriteJSON(msg)
				}

		}
	}
}

func (p *Pool) ReviveWebsocket() {
	if err := recover(); err != nil {
		log.Println("Error in websocket pool:", err)
		go p.Start()
	}
}
