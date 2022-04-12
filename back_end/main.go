package main

import (
	"log"
	"net/http"

	"github.com/YangzhenZhao/shengji/back_end/server"
)

func main() {
	hub := server.NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ServerWS(hub, w, r)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListerAndServer: ", err)
	}
}
