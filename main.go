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

var hub = NewHub()

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
	r.Use(logger)

	log.Println("Starting go-chat server on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.RequestURI)
		h.ServeHTTP(w, r)
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	if err := indexTmpl.Execute(w, nil); err != nil {
		log.Printf("could not execute template: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func newRoom(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("could not parse form: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	roomID := r.PostFormValue("roomID")

	_, err := hub.createRoom(roomID)
	if err != nil {
		log.Printf("could not create room: %v", err)
		http.Error(w, "room already exists", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("http://%s/room/%s", r.Host, roomID), http.StatusFound)
}

func room(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["roomID"]

	if !hub.isRoomExist(roomID) {
		http.Error(w, "chat room does not exist", http.StatusBadRequest)
		return
	}

	if err := roomTmpl.Execute(w, roomID); err != nil {
		log.Printf("could not execute template: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func chat(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["roomID"]
	room, err := hub.getRoom(roomID)
	if err != nil {
		log.Printf("could not get the room: %v", err)
		http.Error(w, "chat room does not exist", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	room.join(conn)

	for {
		messageType, p, err := conn.ReadMessage()
		log.Printf("message received (type %d)", messageType)
		if err != nil {
			log.Println(err)
			room.leave(conn)
			break
		}
		if messageType == websocket.TextMessage {
			room.broadcast(p)
		}
	}
}
