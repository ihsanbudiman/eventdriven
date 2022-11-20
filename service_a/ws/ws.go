package ws

import "github.com/gorilla/websocket"

type Hub struct {
	Public []*websocket.Conn
	Room   map[string][]*websocket.Conn
}

func NewHub() *Hub {
	return &Hub{}
}

// join public
func (h *Hub) Join(conn *websocket.Conn) {
	h.Public = append(h.Public, conn)
}

func (h *Hub) LeavePublic(conn *websocket.Conn) {
	for i, c := range h.Public {
		if c == conn {
			h.Public = append(h.Public[:i], h.Public[i+1:]...)
			break
		}
	}
}

// leave public and room
func (h *Hub) Leave(conn *websocket.Conn) {
	h.LeavePublic(conn)
}

// publish message to all connections in a public conn
func (h *Hub) Publish(msg []byte) error {
	for _, conn := range h.Public {
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			h.LeavePublic(conn)
			return err
		}
	}

	return nil
}
