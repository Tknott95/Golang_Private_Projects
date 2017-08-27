package router

import (
	"html/template"
	"log"
	"net/http"
	"time"

	analytics "github.com/tknott95/Private_Go_Projects/A4S/Controllers/analyticsCtrl"
	auth "github.com/tknott95/Private_Go_Projects/A4S/Controllers/authCtrl"
	hub "github.com/tknott95/Private_Go_Projects/A4S/Controllers/websocketCtrl"
)

var index = template.Must(template.ParseFiles("./Views/index.html"))
var t = time.Now().Format(time.RFC3339)

func InitServer() {
	logger := analytics.CreateLogger("8080-http", "req hit")
	analytics.MonitorTracer()
	go hub.DefaultHub.Start()

	auth.InitKeys()
	log.Println(auth.GenerateJWT("tk", "tk"))

	log.Println(`Server taken off and running on port 8080 🚀`)

	/*  Chat Serverw/ sockets  */
	http.Handle("/", analytics.Time(logger, home))
	http.HandleFunc("/ws", hub.WSHandler)

	/*  Panic Simulator Use for errors where goroutines crash for recovery 99.99% of the time */
	http.Handle("/panic", analytics.Recover(panicSimulator))

	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, req *http.Request) {
	index.Execute(w, nil)
}

func panicSimulator(w http.ResponseWriter, req *http.Request) {
	panic(analytics.ErrInvalidEmail)
}
