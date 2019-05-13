package session

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type ClientManager interface {
	AddClient(conn *websocket.Conn)
	Broadcast(p []byte)
}

func NewClientManager() ClientManager {
	return &clientManagerImpl{}
}

type clientManagerImpl struct {
	clients      []Client
	clientsMutex sync.RWMutex
}

func (cm *clientManagerImpl) AddClient(conn *websocket.Conn) {
	cm.clientsMutex.Lock()
	defer cm.clientsMutex.Unlock()
	c := NewClient(conn)
	cm.clients = append(cm.clients, c)
	log.Printf("new client connected, serving %d clients", len(cm.clients))
}

func (cm *clientManagerImpl) Broadcast(p []byte) {
	cm.clientsMutex.RLock()
	defer cm.clientsMutex.RUnlock()

	disconnectedClients := 0
	for _, c := range cm.clients {
		_, err := c.Send(p)
		if err != nil {
			c.SetConnected(false)
			disconnectedClients++
		}
	}

	if disconnectedClients > 0 {
		newLen := len(cm.clients) - disconnectedClients
		newClients := make([]Client, newLen)
		i := 0
		if newLen > 0 {
			for _, c := range cm.clients {
				if c.IsConnected() {
					newClients[i] = c
					i++
				}
			}
		}

		cm.clients = newClients
		log.Printf("client(s) disconnected, serving %d clients", len(cm.clients))
	}
}

type Client interface {
	Send([]byte) (int, error)
	IsConnected() bool
	SetConnected(bool)
}

func NewClient(conn *websocket.Conn) Client {
	client := clientImpl{conn, true}
	return &client
}

type clientImpl struct {
	conn        *websocket.Conn
	isConnected bool
}

func (c *clientImpl) Send(b []byte) (int, error) {
	err := c.conn.WriteMessage(websocket.TextMessage, b)
	return len(b), err
}

func (c *clientImpl) IsConnected() bool {
	return c.isConnected
}

func (c *clientImpl) SetConnected(connected bool) {
	c.isConnected = connected
}
