package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// do not check origin for now
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", index)
	r.HandleFunc("/room/{roomID}", room)

	log.Fatal(http.ListenAndServe(":8000", r))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "go-chat")
}

func room(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}
