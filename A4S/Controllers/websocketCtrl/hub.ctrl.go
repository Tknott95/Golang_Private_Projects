package hub

import (
	analytics "github.com/tknott95/Private_Go_Projects/A4S/Controllers/analyticsCtrl"
)

var socketMsgLogger = analytics.CreateLogger("socketmsg", "New Socket Msg Made")
var socketNewConnlogger = analytics.CreateLogger("newsocketconn", "New Socket Connection")

var DefaultHub = NewHub()

type Hub struct {
	Join  chan *Conn
	Conns map[*Conn]bool
	Echo  chan string
}

func NewHub() *Hub {
	return &Hub{
		Join:  make(chan *Conn),
		Conns: make(map[*Conn]bool),
		Echo:  make(chan string),
	}
}

func (hub *Hub) Start() {
	for {
		select {
		case conn := <-hub.Join:
			DefaultHub.Conns[conn] = true

		case msg := <-hub.Echo:
			for conn := range hub.Conns {

				conn.Send <- msg
			}
		}
	}
}
