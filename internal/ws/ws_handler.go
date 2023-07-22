package ws

import (
	"encoding/json"
	"fmt"
	"jeopardy/internal/peril"
	"jeopardy/util"
	"log"
	"net/http"

	"github.com/fasthttp/websocket"
)

type Handler struct {
	hub *Hub
}

func NewHandler(h *Hub) *Handler {
	return &Handler{hub: h}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	// /ws/CreateRoom

	// upgrade this to a websocket connection as well -- this is the screen!
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade screen websocket connection: %v\n", err.Error())
	}

	// generate a unique 4 letter game id
	roomId := util.GenerateRoomId()
	_, ok := h.hub.Rooms[roomId]
	for ok {
		roomId = util.GenerateRoomId()
		_, ok = h.hub.Rooms[roomId]
	}

	h.hub.Rooms[roomId] = NewRoom(roomId, conn)
	// h.hub.Rooms[roomId] = &Room{
	// 	ID:        roomId,
	// 	GameState: peril.NewGame(),
	// 	Players:   make(map[int]*Client),
	// 	Screen: &Client{
	// 		Conn: conn,
	// 	},
	// }

	// send the room name
	conn.WriteMessage(websocket.TextMessage, []byte(h.hub.Rooms[roomId].ID))

}

func (h *Handler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade player websocket connection: %v\n", err.Error())
	}

	// /ws/JoinRoom?roomId={roomID}
	// get query string!
	roomId := r.URL.Query().Get("roomId")

	// check that room with that roomId exists!!!
	if _, ok := h.hub.Rooms[roomId]; !ok {
		conn.WriteMessage(websocket.TextMessage, []byte("Game not found. Websocket connection terminated"))
		conn.Close()
		return
	}

	room := h.hub.Rooms[roomId]
	room.Mutex.Lock()
	// and that there is space for more participants!!
	if len(room.Players) >= 5 {
		conn.WriteMessage(websocket.TextMessage, []byte("Game full. Websocket connection terminated"))
		conn.Close()
		return
	}

	id := len(room.Players) + 1
	cl := &Client{
		Conn:   conn,
		ID:     id,
		RoomID: roomId,
		PlayerState: &peril.PlayerState{
			Score: 0,
			Role:  peril.Waiter,
			ID:    id,
		},
	}

	h.hub.Register <- cl
	room.Mutex.Unlock()

	conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Connected to %s", h.hub.Rooms[roomId].ID)))

	// this will be in a message! there will already be a conn!
	// Let user pick their username
	// if _, ok := h.hub.Rooms[roomId].Players[playerId]; !ok {
	// 	w.Write([]byte("Username taken"))
	// }
}

type RoomRes struct {
	ID string `json:"id"`
}

func (h *Handler) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms := make([]RoomRes, 0)

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{ID: r.ID})
	}

	roomJson, _ := json.Marshal(rooms)
	w.Write([]byte(roomJson))
}
