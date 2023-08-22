package ws

import (
	"jeopardy/internal/peril"
	"sync"

	"github.com/fasthttp/websocket"
)

type Room struct {
	ID        string
	GameState *peril.Game
	Players   map[int]*Client
	Screen    *Client
	Mutex     sync.Mutex
}

func NewRoom(id string, conn *websocket.Conn) *Room {
	return &Room{
		ID:        id,
		GameState: peril.NewGame(),
		Players:   make(map[int]*Client),
		Screen: &Client{
			Conn: conn,
		},
	}
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Room
	Close      chan *Room
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Room),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			if _, ok := h.Rooms[cl.RoomID]; ok {
				r := h.Rooms[cl.RoomID]
				r.Mutex.Lock()

				if _, ok := r.Players[cl.ID]; !ok {
					r.Players[cl.ID] = cl
				}
				r.Mutex.Unlock()
			}
		case cl := <-h.Unregister:
			if _, ok := h.Rooms[cl.RoomID]; ok {
				if _, ok := h.Rooms[cl.RoomID].Players[cl.ID]; ok {
					delete(h.Rooms[cl.RoomID].Players, cl.ID)
					close(cl.Message)
				}
			}
		// case r := <-h.Close:
		// this will come from the screen, after game is done
		// r.ID
		// close this room: delete from hub

		case m := <-h.Broadcast:
			for _, cl := range m.Players {
				cl.Message <- m.GameState
			}

			// send to screen as well!
			m.Screen.Message <- m.GameState
		}
	}
}
