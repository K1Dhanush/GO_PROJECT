package basicauth

import (
	"net/http"
)

var (
	credentials = map[string]string{
		"user": "Dhaneshwar"}
)

type MiddlewareFunc func(http.Handler) http.Handler

func NewAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//for explicting Certain requests
		if r.URL.Path == "/jwt-token" {
			// Call the next handler directly
			next.ServeHTTP(w, r)
			return
		}

		username, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "Authentication credentials not provided", http.StatusUnauthorized)
			return
		}
		if credPasswd, found := credentials[username]; found && password == credPasswd {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		}
	})
}

/* Custom -- middleWare Auth*/

// type Basic struct {
// 	handler http.Handler
// }

// // For Authentication
// func GetAuthenticated(r *http.Request) error {
// 	username, password, ok := r.BasicAuth()
// 	if !ok {
// 		//errors is pkg
// 		return errors.New("authentication credentials not provided")
// 	}
// 	if credPasswd, found := credentails[username]; found && password == credPasswd {
// 		return nil
// 	} else {
// 		return errors.New("InValid Credentails")
// 	}
// }

// // Implementing the http.Hander method ServeHttp
// func (b *Basic) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	err := GetAuthenticated(r)
// 	if err != nil {
// 		http.Error(w, "Please... Enter the Valid Details to Login.", http.StatusUnauthorized)
// 		return
// 	}

// 	//calling the original Handler or Request
// 	b.handler.ServeHTTP(w, r)
// }

// func NewLogger(handlerToWrap http.Handler) *Basic { // step --1
// 	return &Basic{handlerToWrap}
// }
