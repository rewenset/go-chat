package main

import (
	"log"

	"github.com/gorilla/websocket"
)

// SimpleRoom store only WebSocket connections and all messages in one slice
type SimpleRoom struct {
	users    map[*websocket.Conn]bool
	messages [][]byte
}

// join add conection to the users of selected room
func (r *SimpleRoom) join(conn *websocket.Conn) {
	r.users[conn] = true
}

// leave removes connection from the users of selected room
func (r *SimpleRoom) leave(conn *websocket.Conn) {
	delete(r.users, conn)
}

// broadcast send message to all users in the room
func (r *SimpleRoom) broadcast(msg []byte) {
	r.messages = append(r.messages, msg)
	for conn := range r.users {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			r.leave(conn)
			log.Println(err)
		}
	}
}
