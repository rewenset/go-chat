package main

import (
	"log"

	"github.com/gorilla/websocket"
)

// Room store only WebSocket connections, broadcast channel and all messages
type Room struct {
	users    map[*websocket.Conn]bool
	bc       chan []byte
	messages [][]byte
}

// NewRoom creates the room and starts goroutine to listen on broadcasted messages
func NewRoom() *Room {
	r := Room{
		users: make(map[*websocket.Conn]bool),
		bc:    make(chan []byte),
	}
	go func() {
		for m := range r.bc {
			go r.broadcast(m)
		}
	}()
	return &r
}

// join add conection to the users of selected room
func (r *Room) join(conn *websocket.Conn) {
	r.users[conn] = true
}

// leave removes connection from the users of selected room
func (r *Room) leave(conn *websocket.Conn) {
	delete(r.users, conn)
}

// broadcast send message to all users in the room
func (r *Room) broadcast(msg []byte) {
	r.messages = append(r.messages, msg)
	for conn := range r.users {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			r.leave(conn)
			log.Println(err)
		}
	}
}
