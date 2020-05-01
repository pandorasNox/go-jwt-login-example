package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func main() {
	http.HandleFunc("/jwt", handleJWT)
	http.HandleFunc("/someresource", handleSomeResource)
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

func handleSomeResource(w http.ResponseWriter, r *http.Request) {

	headerToken := r.Header.Get("Authorization")
	splitToken := strings.Split(headerToken, "Bearer")
	maybeJWT := strings.TrimSpace(splitToken[1])

	token, err := validateToken(maybeJWT)
	_ = token
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "plain/text")
	// w.Write([]byte(token))
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		fmt.Println(token.Method.Alg())

		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("parsing the jwt didn't worked: %s", err)
	}

	//this is a type assertion @see https://tour.golang.org/methods/15
	//type MapClaims map[string]interface{}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// fmt.Println(claims["foo"], claims["nbf"])
		return claims, nil
	} else {
		// fmt.Println(err)
		return nil, fmt.Errorf("claims not ok or token not valid: %s", err)
	}
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
