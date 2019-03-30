package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// Hub allows to manage rooms
type Hub struct {
	rooms map[string]*SimpleRoom
}

// NewHub  creates a hub
func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]*SimpleRoom),
	}
}

// createRoom creates simpleRoom with no users and saves it
func (h *Hub) createRoom(roomID string) (*SimpleRoom, error) {
	if h.isRoomExist(roomID) {
		return nil, fmt.Errorf("the room with id %s already exists", roomID)
	}

	r := &SimpleRoom{
		users: make(map[*websocket.Conn]bool),
	}
	h.rooms[roomID] = r
	return r, nil
}

// getRoom creates simpleRoom with no users and saves it
func (h *Hub) getRoom(roomID string) (*SimpleRoom, error) {
	if r, ok := h.rooms[roomID]; ok {
		return r, nil
	}
	return nil, fmt.Errorf("the room with id %s does not exist", roomID)
}

// isRoomExist checks if room with is already created
func (h *Hub) isRoomExist(roomID string) bool {
	if _, ok := h.rooms[roomID]; ok {
		return true
	}
	return false
}
