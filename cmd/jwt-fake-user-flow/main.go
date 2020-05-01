package main

import (
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func main() {
	http.HandleFunc("/jwt", handleJWT)
	log.Printf("server stopped: %s", http.ListenAndServe(":8080", nil))
}

// client makes first call to server and gets new identity & identityToken (registration)
// getInfo(identityToken) return number of request with this token

func handleJWT(w http.ResponseWriter, r *http.Request) {
	token, err := generateToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "plain/text")
	w.Write([]byte(token))
}

var hmacSampleSecret = []byte("my_secret_key")

func generateToken() (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{

		// these 4 are recommendet to be always there
		// @see https://auth0.com/blog/a-look-at-the-latest-draft-for-jwt-bcp/
		// under "Validate All Possible Claims"
		"sub": "123456",
		"iss": "app.route.jwt",
		"aud": "",
		"exp": "",

		// ???
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),

		//custom stuff ???
		"foo":  "bar",
		"role": "admin",
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)

	return tokenString, err
}
