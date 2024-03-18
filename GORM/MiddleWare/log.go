package logging

import (
	"log"
	"net/http"
	"time"
)

type MiddlewareFunc func(http.Handler) http.Handler

// Logging is a middleware function that logs information about incoming requests.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Println(r.URL.Path)
		log.Println(start)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
