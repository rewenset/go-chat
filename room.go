package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// SimpleRoom store only WebSocket connections
type SimpleRoom struct {
	users    map[*websocket.Conn]bool
	messages [][]byte
}

var rooms map[string]*SimpleRoom

func init() {
	rooms = make(map[string]*SimpleRoom)
}

// CreateRoom creates simpleRoom with no users and saves it
func CreateRoom(roomID string) error {
	if IsRoomExist(roomID) {
		return fmt.Errorf("room with id %s already exists", roomID)
	}
	rooms[roomID] = &SimpleRoom{
		users: make(map[*websocket.Conn]bool),
	}
	return nil
}

//IsRoomExist checks if room with is already created
func IsRoomExist(roomID string) bool {
	if _, ok := rooms[roomID]; ok {
		return true
	}
	return false
}

// JoinRoom add conection to the users of selected room
func JoinRoom(roomID string, conn *websocket.Conn) {
	rooms[roomID].users[conn] = true
}

// LeaveRoom removes connection from the users of selected room
func LeaveRoom(roomID string, conn *websocket.Conn) {
	delete(rooms[roomID].users, conn)
}

// BroadcastInRoom send message to all users in the room
func BroadcastInRoom(roomID string, msg []byte) {
	chatRoom := rooms[roomID]
	chatRoom.messages = append(chatRoom.messages, msg)
	log.Println(chatRoom.messages)
	for conn := range chatRoom.users {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			LeaveRoom(roomID, conn)
			log.Println(err)
		}
	}
}
