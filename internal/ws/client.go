package ws

import (
	"jeopardy/internal/peril"

	"github.com/fasthttp/websocket"
)

type Client struct {
	Conn        *websocket.Conn
	ID          int
	RoomID      string
	Username    string
	Message     chan *Message
	Type        clientType
	PlayerState *peril.PlayerState
}

type clientType int

const (
	screen clientType = iota
	player
)

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
		// _, m, err := c.Conn.ReadMessage()
		// if err != nil {
		// 	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		// 		log.Printf("error: %v\n", err)
		// 	}
		// 	break
		// }

		// change game state here

		// hub.Broadcast <- <new game state>

	}
}
