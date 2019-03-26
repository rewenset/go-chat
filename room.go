package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// SimpleRoom store only WebSocket connections
type SimpleRoom struct {
	users []*websocket.Conn
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
	rooms[roomID] = &SimpleRoom{}
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
	rooms[roomID].users = append(rooms[roomID].users, conn)
}

// BroadcastInRoom send message to all users
func BroadcastInRoom(roomID string, msg []byte) {
	chatRoom := rooms[roomID]
	for i, c := range chatRoom.users {
		if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
			chatRoom.users = append(chatRoom.users[:i], chatRoom.users[i+1:]...)
			log.Println(err)
		}
	}
}
