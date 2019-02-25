package handlers

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

//Notifier is an object that handles WebSocket notifications.
type Notifier struct {
	//eventQ is a go channel that
	//into which one goroutine can
	//write byte slices, and out of which
	//another goroutine can read those byte slices
	//eventQ chan []byte

	//protect it for concurrent use!
	//track all of the current WebSocket
	//connections
	clients map[int64]*websocket.Conn

	//slice will be used by multiple goroutines, need to
	//protect it for concurrent use
	mx sync.RWMutex
}

//NewNotifier constructs a new Notifier
func NewNotifier() *Notifier {
	n := make(map[int64]*websocket.Conn)
	return &Notifier{
		clients: n,
	}
}

//AddClient adds a new client to the Notifier
func (n *Notifier) AddClient(client *websocket.Conn, userID int64) {
	n.mx.Lock()
	defer n.mx.Unlock()

	n.clients[userID] = client

	//process incoming control messages from
	//the client
	go (func(userID int64) {
		for {
			if _, _, err := n.clients[userID].NextReader(); err != nil {
				n.mx.Lock()
				if n.clients[userID] != nil {
					n.clients[userID].Close()
					delete(n.clients, userID)
				}
				n.mx.Unlock()
				break
			}
		}
	})(userID)
}

//start starts the notification loop
func (n *Notifier) start(msg []byte) {

	prepMsg, err := websocket.NewPreparedMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Printf("error writing prep message", err)
		return
	}

	for userID, conn := range n.clients {
		_, exists := n.clients[userID]
		if exists {
			n.writePreppedMsg(prepMsg, userID, conn)
		}
	}
}

func (n *Notifier) writePreppedMsg(msg *websocket.PreparedMessage, userID int64, conn *websocket.Conn) {
	if err := conn.WritePreparedMessage(msg); err != nil {
		n.mx.Lock()
		delete(n.clients, userID)
		n.mx.Unlock()
	}
}
