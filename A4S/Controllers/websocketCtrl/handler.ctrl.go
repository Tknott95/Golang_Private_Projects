package hub

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Conn struct {
	WS   *websocket.Conn
	Send chan string
}

func (conn *Conn) SendToHub() {
	defer conn.WS.Close()
	for {
		_, msg, err := conn.WS.ReadMessage()
		if err != nil {
			/*  user disconnected, most likely refreshed browser jsut return */
			return
		}
		DefaultHub.Echo <- string(msg)
	}
}

func (conn *Conn) Write(msg string) error {
	return conn.WS.WriteMessage(websocket.TextMessage, []byte(msg))
}

/*  sends msg from hub to websocket connection */
func (conn *Conn) ReceiveFromHub() {
	defer conn.WS.Close()
	for {
		conn.Write(<-conn.Send)
	}
}

/*  handles the http request */
func WSHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := WSUpgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	conn := &Conn{
		Send: make(chan string),
		WS:   ws,
	}
	DefaultHub.Join <- conn

	go conn.SendToHub()
	conn.ReceiveFromHub()
}

/*  used to upgrade the protocol to allow websockets */
var WSUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
