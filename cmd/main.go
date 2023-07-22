package main

import (
	"jeopardy/internal/ws"
	"jeopardy/router"
	"log"
)

func main() {
	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)

	go hub.Run()

	r := router.NewRouter(wsHandler)
	log.Fatal(router.Start(":8080", r))
}
