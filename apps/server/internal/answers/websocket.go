package answers

import (
	"database/sql"
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Hub struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]struct{}
}

type event struct {
	Type    string   `json:"type"`
	Answers []answer `json:"answers,omitempty"`
	Answer  *answer  `json:"answer,omitempty"`
}

func NewHub() *Hub {
	return &Hub{clients: make(map[*websocket.Conn]struct{})}
}

func (h *Hub) Handler(db *sql.DB) echo.HandlerFunc {
	upgrader := websocket.Upgrader{}

	return func(c echo.Context) error {
		connection, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer connection.Close()

		h.mu.Lock()
		h.clients[connection] = struct{}{}
		items, err := list(c.Request().Context(), db)
		if err == nil {
			err = connection.WriteJSON(event{Type: "snapshot", Answers: items})
		}
		h.mu.Unlock()
		if err != nil {
			return err
		}
		defer h.remove(connection)

		for {
			if _, _, err := connection.ReadMessage(); err != nil {
				return nil
			}
		}
	}
}

func (h *Hub) Broadcast(item answer) {
	payload, err := json.Marshal(event{Type: "answer", Answer: &item})
	if err != nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	for connection := range h.clients {
		_ = connection.SetWriteDeadline(time.Now().Add(5 * time.Second))
		if err := connection.WriteMessage(websocket.TextMessage, payload); err != nil {
			delete(h.clients, connection)
			_ = connection.Close()
		}
	}
}

func (h *Hub) remove(connection *websocket.Conn) {
	h.mu.Lock()
	delete(h.clients, connection)
	h.mu.Unlock()
}
