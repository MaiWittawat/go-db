package realtime

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	clients  map[string]*websocket.Conn
	lock     sync.Mutex
	upgrader websocket.Upgrader
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		clients: make(map[string]*websocket.Conn),
		lock:    sync.Mutex{},
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true }, // สำหรับ dev
		},
	}
}

func (ws *WebSocketServer) Online(userID string, conn *websocket.Conn) error {
	ws.lock.Lock()
	defer ws.lock.Unlock()
	ws.clients[userID] = conn
	return nil
}

func (ws *WebSocketServer) Offline(userID string, conn *websocket.Conn) error {
	ws.lock.Lock()
	defer ws.lock.Unlock()
	delete(ws.clients, userID)
	return conn.Close()
}

func (ws *WebSocketServer) SendTo(userID string, data any) error {
	ws.lock.Lock()
	conn, ok := ws.clients[userID]
	ws.lock.Unlock()
	if !ok {
		return errors.New("user not connected")
	}
	return conn.WriteJSON(data)
}