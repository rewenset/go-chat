package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var indexTmpl = template.Must(template.ParseFiles("templates/index.html"))
var roomTmpl = template.Must(template.ParseFiles("templates/room.html"))

var upgrader = websocket.Upgrader{
	// do not check origin for now
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", index)
	r.HandleFunc("/room", newRoom).Methods("POST")
	r.HandleFunc("/room/{roomID}", room)
	r.HandleFunc("/room/{roomID}/chat", chat)

	log.Fatal(http.ListenAndServe(":8000", r))
}

func index(w http.ResponseWriter, r *http.Request) {
	if err := indexTmpl.Execute(w, nil); err != nil {
		log.Printf("could not execute template: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func newRoom(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("could not parse form: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	} else {
		roomID := r.PostFormValue("roomID")
		err = CreateRoom(roomID)
		if err != nil {
			log.Printf("could not create room: %v", err)
			http.Error(w, "room already exists", http.StatusBadRequest)
		} else {
			http.Redirect(w, r, fmt.Sprintf("http://%s/room/%s", r.Host, roomID), http.StatusFound)
		}
	}
}

func room(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["roomID"]

	if !IsRoomExist(roomID) {
		http.Error(w, "chat room does not exist", http.StatusBadRequest)
		return
	}

	if err := roomTmpl.Execute(w, roomID); err != nil {
		log.Printf("could not execute template: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func chat(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["roomID"]

	if !IsRoomExist(roomID) {
		http.Error(w, "chat room does not exist", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	JoinRoom(roomID, conn)

	for {
		messageType, p, err := conn.ReadMessage()
		log.Printf("received message with type %d", messageType)
		if err != nil {
			log.Println(err)
			return
		}
		if messageType == websocket.TextMessage {
			BroadcastInRoom(roomID, p)
		}
	}
}
