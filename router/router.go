package router

import (
	"jeopardy/internal/ws"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(wsHandler *ws.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	r.Get("/ws/createRoom", wsHandler.CreateRoom)
	r.Get("/ws/joinRoom", wsHandler.JoinRoom)
	r.Get("/getRooms", wsHandler.GetRooms)
	return r
}

func Start(addr string, mux *chi.Mux) error {
	log.Println("Starting server on port 8080")
	return http.ListenAndServe(addr, mux)
}
