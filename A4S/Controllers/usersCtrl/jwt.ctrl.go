package users

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type MyCustomClaims struct {
	Sub string `json:"name`
	Exp int64  `json:"exp"`
	Iat int64  `json:"iat"`
	jwt.StandardClaims
}

var (
	identityURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	provider    = New()
	signingKey  = genRandBytes()
)

// New creates a new oauth2 config
func New() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_KEY"),
		ClientSecret: os.Getenv("GOOGLE_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:3000/auth/gplus/callback",
		Scopes:       []string{"email", "profile"},
	}
}

// AuthURLHandler just redirects the user to correct Oauth sign in page for our
// provider.
func AuthURLHandler(w http.ResponseWriter, r *http.Request) {
	authURL := provider.AuthCodeURL("", oauth2.AccessTypeOffline)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// CallbackURLHandler handles all of the Oauth flow.
func CallbackURLHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := provider.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	client := provider.Client(oauth2.NoContext, token)
	resp, err := client.Get(identityURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	defer resp.Body.Close()

	user := make(map[string]string)
	// This decode will error, since it can't decode every value returned
	// from the API. However, we don't need to worry about this, since it can
	// correctly decode our user's email address.
	json.NewDecoder(resp.Body).Decode(&user)

	email := user["email"]
	genToken(w, email)
}

func genToken(w http.ResponseWriter, user string) {
	sub := user
	exp := time.Now().Add(time.Hour * 72).Unix()
	iat := time.Now().Unix()

	claims := MyCustomClaims{
		sub,
		exp,
		iat,
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// World's laziest way to issue a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\n\ttoken: " + tokenString + "\n}"))
}

// VerifyToken gets the token from an HTTP request, and ensures that it's
// valid. It'll return the user's username as a string.
func VerifyToken(r *http.Request) (string, error) {
	token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return signingKey, nil
	})
	if err != nil {
		return "", err
	}
	if token.Valid == false {
		return "", jwt.ErrInvalidKey
	}
	return token.Claims["sub"].(string), nil
}
