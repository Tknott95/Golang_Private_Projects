package hub

import (
	"encoding/json"
	"net/http"
)

func Notify(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		if req.Header.Get("Content-Type") == "application/json" {
			jq := make(map[string]string)
			json.NewDecoder(req.Body).Decode(&jq)
			msg := jq["message"]

			if msg == "" {
				http.Error(w, `Messages must have the key "messages"`, http.StatusBadRequest)
				return
			}
			DefaultHub.Echo <- msg
			w.Write([]byte("Sen Message successfully"))
			return
		}

		msg := req.FormValue("message")
		req.ParseForm()

		if msg == "" {
			http.Error(w, "No message found in request", http.StatusBadRequest)
			return
		}

		DefaultHub.Echo <- msg
		w.Write([]byte("Sent message successfully"))
		return
	default:
		http.Error(w, "Only POST method supported", http.StatusMethodNotAllowed)
		return
	}
}
