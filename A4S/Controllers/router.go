package router

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	analytics "github.com/tknott95/Private_Go_Projects/Concurrency4Go/Controllers/analyticsCtrl"
	users "github.com/tknott95/Private_Go_Projects/Concurrency4Go/Controllers/usersCtrl"
	hub "github.com/tknott95/Private_Go_Projects/Concurrency4Go/Controllers/websocketCtrl"
)

const loginTemplate = `
<h1>Enter your username and password</h1>
<form action="/" method="POST">
	<input type="text" name="user" required>

	<label for="password">Password</label>
	<input type="password" name="password" required>

	<input type="submit" value="Submit">
</form>
`

func authHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		t, _ := template.New("login").Parse(loginTemplate)
		t.Execute(w, nil)
	case "POST":
		user := r.FormValue("user")
		pass := r.FormValue("password")
		err := users.AuthenticateUser(user, pass)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users.SetSession(w, user)
		w.Write([]byte("Signed in successfully"))
	}
}

func restrictedHandler(w http.ResponseWriter, r *http.Request) {
	user := users.GetSession(w, r)
	w.Write([]byte(user))
}

func oauthRestrictedHandler(w http.ResponseWriter, r *http.Request) {
	user, err := users.VerifyToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Write([]byte(user))
}

func sanitizeInputExample(str string) {
	fmt.Println("JS: ", template.JSEscapeString(str))
	fmt.Println("HTML: ", template.HTMLEscapeString(str))
}

var index = template.Must(template.ParseFiles("./Views/index.html"))
var t = time.Now().Format(time.RFC3339)

func InitServer() {
	sanitizeInputExample("<script>alert(\"Hi!\");</sciprt>")

	username, password := "ben", "qwerty123"

	err := users.NewUser(username, password)
	if err != nil {
		fmt.Printf("User already exists: %s\n", err.Error())
	} else {
		fmt.Printf("Succesfully created and authenticated user \033[32m%s\033[0m\n", username)
	}

	logger := analytics.CreateLogger("8080-http", "req hit")
	analytics.MonitorTracer()
	go hub.DefaultHub.Start()

	log.Println(`Server taken off and running on port 8080 🚀`)

	/*  Chat Serverw/ sockets  */
	http.Handle("/", analytics.Time(logger, home))
	http.HandleFunc("/ws", hub.WSHandler)

	/* Oauth for User Functionality */
	http.HandleFunc("/auth/gplus/authorize", users.AuthURLHandler)
	http.HandleFunc("/auth/gplus/callback", users.CallbackURLHandler)
	http.HandleFunc("/oauth", oauthRestrictedHandler)
	http.HandleFunc("/restricted", restrictedHandler)

	/*  Panic Simulator Use for errors where goroutines crash for recovery 99.99% of the time */
	http.Handle("/panic", analytics.Recover(panicSimulator))

	http.ListenAndServe(":8080", nil)
}

func oauthRestrictedHandler(w http.ResponseWriter, r *http.Request) {
	user, err := users.VerifyToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Write([]byte(user))
}

func home(w http.ResponseWriter, req *http.Request) {
	index.Execute(w, nil)
}

func panicSimulator(w http.ResponseWriter, req *http.Request) {
	panic(analytics.ErrInvalidEmail)
}

func restrictedHandler(w http.ResponseWriter, r *http.Request) {
	user := users.GetSession(w, r)
	w.Write([]byte(user))
}
