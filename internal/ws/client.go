package ws

import (
	"jeopardy/internal/peril"
	"log"

	"github.com/fasthttp/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	ID       int
	RoomID   string
	Username string
	Message  chan *peril.Game
	// Type        clientType
}

// type clientType int

// const (
// 	screen clientType = iota
// 	player
// )

type Message interface {
}

// write a message to client
func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		// you get a certain Game State
		// decide here what to actually write, based on game state and client player state / screen status

		c.Conn.WriteJSON(message)
	}
}

// read message coming from a client
func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v\n", err)
			}
			break
		}

		hub.Rooms[c.RoomID].GameState.UpdateGame(c.ID, m)

		hub.Broadcast <- hub.Rooms[c.RoomID]
	}
}
