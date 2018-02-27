package handler

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Clients struct {
	Conn      map[*websocket.Conn]bool
	Mutex     sync.Mutex
	Broadcast chan interface{}
}

func HandleConnection(upgrader *websocket.Upgrader, clients *Clients, w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer c.Close()

	clients.Conn[c] = true
	for {
		var res interface{}
		err := c.ReadJSON(&res)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("IsUnexpectedCloseError : %s\n", err)
			}
			log.Printf("ReadJSON : %s\n", err)
			delete(clients.Conn, c)
			break
		}

		clients.Broadcast <- res
	}
}

func HandleMessage(clients *Clients) {
	clients.Mutex.Lock()
	defer clients.Mutex.Unlock()

	for {
		select {
		case msg, ok := <-clients.Broadcast:
			for conn := range clients.Conn {
				defer conn.Close()

				if !ok {
					conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}

				err := conn.WriteJSON(msg)
				if err != nil {
					log.Printf("WriteJSON : %s\n", err)
					delete(clients.Conn, conn)
				}
			}
		}
	}
}
