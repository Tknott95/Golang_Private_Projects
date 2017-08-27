package userModel

var UsersDB = map[string]User{}
var SessionDB = map[string]string{}

type User struct {
	UID           int
	Username      string
	Fullname      string
	Notifications int
	Created       string
	Password      string
	Email         string
	IsAdmin       int
}

type JsonSignup struct {
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type JsonLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
