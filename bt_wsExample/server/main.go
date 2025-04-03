// server.go
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	hub := newHub()
	go hub.run()
	router := mux.NewRouter()
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
