// Copied from https://github.com/gorilla/websocket/blob/master/examples/chat/hub.go

package sockets

import (
	"github.com/justinorringer/pal-pad-go/db"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	clients map[*Client]bool // Registered clients.

	broadcast chan []byte // Inbound messages from the clients.

	register chan *Client // Register requests from the clients.

	unregister chan *Client // Unregister requests from clients.
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run(rc *db.RedisClient) {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			// process the line and store it in the database
			err = ProcessMessage(rc)

			for client := range h.clients {
				// if the message is for a specific sketch, only send it to clients with sketchID
				// if client.sketchID != uuid.Nil && client.sketchID != message.sketchID {
				// 	continue
				// }
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
