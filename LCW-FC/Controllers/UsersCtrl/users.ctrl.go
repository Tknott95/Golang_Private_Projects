package userCtrl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	b64 "encoding/base64"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	databaseCtrl "github.com/tknott95/LCW-FC/Controllers/DatabaseCtrl"
	userModel "github.com/tknott95/LCW-FC/Models/user"
)

type Token struct {
	Raw       string                 // The raw token.  Populated when you Parse a token
	Method    SigningMethod          // The signing method used or to be used
	Header    map[string]interface{} // The first segment of the token
	Claims    Claims                 // The second segment of the token
	Signature string                 // The third segment of the token.  Populated when you Parse a token
	Valid     bool                   // Is the token valid?  Populated when you Parse/Verify a token
}

type Claims interface {
	Valid() error
}

type SigningMethod interface {
	Verify(signingString, signature string, key interface{}) error // Returns nil if signature is valid
	Sign(signingString string, key interface{}) (string, error)    // Returns encoded signature or error
	Alg() string                                                   // returns the alg identifier for this method (example: 'HS256')
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {

}

func Logout(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	os.Setenv("admin", "false")
	println("\nAdmin Logged Out\n")
}

func Login(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var loginRes userModel.JsonLogin
	decoder := json.NewDecoder(req.Body)
	decoder.Decode(&loginRes)

	b64.StdEncoding.DecodeString(loginRes.Password)

	w.Header().Set("loginres", "Logged in correctly as "+loginRes.Email)

	/*  sqlModel.SQLStore */
	var _userModel = userModel.User{}

	rows, err := databaseCtrl.Store.DB.Query(`SELECT * FROM users;`)
	fmt.Println(w, "Established users db connection", nil)
	if err != nil {
		println("User fetch error: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&_userModel.UID, &_userModel.Username, &_userModel.Fullname, &_userModel.Notifications, &_userModel.Created, &_userModel.Password, &_userModel.Email, &_userModel.IsAdmin)
	}

	b64.StdEncoding.DecodeString(loginRes.Password)

	fmt.Println("Given: ", loginRes.Email, " Act: ", _userModel.Email, " Given Pass: ", loginRes.Password, "Act Pass: ", _userModel.Password)

	if loginRes.Email == _userModel.Email && loginRes.Password == _userModel.Password {

		// token := jwt.New(jwt.SigningMethodRS256)
		// claims := make(jwt.MapClaims)
		// claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
		// claims["iat"] = time.Now().Unix()
		// token.Claims = claims

		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	fmt.Fprintln(w, "Error extracting the key")
		// 	fatal(err)
		// }

		// fmt.Println(token)

		// tokenString, err := token.SignedString(signKey)

		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	fmt.Fprintln(w, "Error while signing the token")
		// 	fatal(err)
		// }

		// response := Token{tokenString}
		// // fmt.Println(response)
		// JsonResponse(response, w)

		mySigningKey := []byte("AllYourBase")

		type MyCustomClaims struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Exp      int64  `json:"exp"`
			Iat      int64  `json:"iat"`
			jwt.StandardClaims
		}

		exp := time.Now().Add(time.Hour * time.Duration(1)).Unix()
		iat := time.Now().Unix()
		// Create the Claims
		claims := MyCustomClaims{
			loginRes.Email,
			loginRes.Password,
			exp,
			iat,
			jwt.StandardClaims{
				ExpiresAt: 15000,
				Issuer:    "test",
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, err := token.SignedString(mySigningKey)
		fmt.Printf("%v %v", ss, err)

		println("User LOGGED IN CORRECTLY")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"loggedin": true,"admin": false}`)
		fmt.Println("INFO: ", _userModel.Email, _userModel.Password)
		return
	}

	io.WriteString(w, `{"loggedin": false}`)
	w.WriteHeader(http.StatusUnauthorized)
}

func IsAdmin() bool {
	if os.Getenv("admin") == "true" {
		return true
	}
	return false
}
