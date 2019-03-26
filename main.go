package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type simpleRoom struct {
	users []*websocket.Conn
}

var rooms map[string]*simpleRoom
var roomTmpl = template.Must(template.ParseFiles("room.html"))

var upgrader = websocket.Upgrader{
	// do not check origin for now
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	rooms = make(map[string]*simpleRoom)

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", index)
	r.HandleFunc("/room/{roomID}", room)
	r.HandleFunc("/room/{roomID}/chat", chat)

	log.Fatal(http.ListenAndServe(":8000", r))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "go-chat")
}

func room(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if _, ok := rooms[vars["roomID"]]; !ok {
		log.Printf("creating new room: %s", vars["roomID"])
		rooms[vars["roomID"]] = &simpleRoom{}
	}

	if err := roomTmpl.Execute(w, vars["roomID"]); err != nil {
		log.Printf("could not execute template: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func chat(w http.ResponseWriter, r *http.Request) {
	var chatRoom *simpleRoom
	if cr, ok := rooms[mux.Vars(r)["roomID"]]; ok {
		chatRoom = cr
	} else {
		http.Error(w, "chat room does not exist", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	chatRoom.users = append(chatRoom.users, conn)

	for {
		messageType, p, err := conn.ReadMessage()
		log.Printf("received message with type %d", messageType)
		if err != nil {
			log.Println(err)
			return
		}
		if messageType == websocket.TextMessage {
			for i, c := range chatRoom.users {
				if err := c.WriteMessage(websocket.TextMessage, p); err != nil {
					chatRoom.users = append(chatRoom.users[:i], chatRoom.users[i+1:]...)
					log.Println(err)
				}
			}
		}
	}
}
