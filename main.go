package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var oneRoom []*websocket.Conn
var roomTmpl = template.Must(template.ParseFiles("room.html"))

var upgrader = websocket.Upgrader{
	// do not check origin for now
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	oneRoom = make([]*websocket.Conn, 0, 5)

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
	if err := roomTmpl.Execute(w, vars["roomID"]); err != nil {
		log.Printf("could not execute template: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func chat(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	oneRoom = append(oneRoom, conn)

	for {
		messageType, p, err := conn.ReadMessage()
		log.Printf("received message with type %d", messageType)
		if err != nil {
			log.Println(err)
			return
		}
		if messageType == websocket.TextMessage {
			for i, c := range oneRoom {
				if err := c.WriteMessage(websocket.TextMessage, p); err != nil {
					oneRoom = append(oneRoom[:i], oneRoom[i+1:]...)
					log.Println(err)
				}
			}
		}
	}
}
