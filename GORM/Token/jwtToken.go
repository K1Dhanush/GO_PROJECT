package token

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Global var map
var (
	credentails = map[string]string{
		"user": "abc@123"}

	secretKey = []byte("secret-key")
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Token(w http.ResponseWriter, r *http.Request) {
	//Taking credentails from the Request Body
	var user User
	//Decoding JSON to struct
	res := json.NewDecoder(r.Body).Decode(&user)
	if res != nil {
		http.Error(w, "DATA is not in JSON Format", http.StatusInternalServerError)
		return
	}

	if credPasswd, found := credentails[user.Username]; found && credPasswd == user.Password {
		tokenString, err := CreateToken(user.Username)
		if err != nil {
			http.Error(w, "Token is not Created", http.StatusInternalServerError)
			return
		}
		//we are using in Fprintln in HandlerFunc So, that we have to use w--io.Writer
		fmt.Fprintln(w, tokenString)
	} else {
		fmt.Fprintln(w, "Invalid Credentails")
	}
}

// func extractToken(r *http.Request) (string, error) {
// 	tokenString := r.Header.Get("Authorization")
// 	if tokenString == "" {
// 		return "", fmt.Errorf("Please Enter the Credentails")
// 	}
// 	//while passing the tokenString we have to ---remove the Bearer inFront of Token
// 	tokenString = strings.Split(tokenString, " ")[1] //no of times to split the token
// 	return tokenString, nil
// }

// func verifyToken(tokenString string) error {
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return secretKey, nil
// 	})
// 	log.Println(token)
// 	if err != nil {
// 		return err
// 	}
// 	if !token.Valid {
// 		return fmt.Errorf("in-valid Token")
// 	}
// 	return nil
// }

func CreateToken(username string) (string, error) {
	claims := &jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Second * 30).Unix(),
	}

	//Creates a JWT Token with signing method and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// log.Println(token)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
