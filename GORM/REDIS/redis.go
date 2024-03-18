package Caching

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis"
)

// var ctx = context.Background()
var Client *redis.Client // Declare a global variable for Redis client

// SetRedisClient sets the Redis client
func SetRedisClient(redisClient *redis.Client) {
	Client = redisClient
}

// Struct
type Product struct {
	Item  string `json:"item"`
	Price int    `json:"price"`
}

// for variable
type JSONResponse struct {
	Users []Product `json:"products"`
}

type MiddlewareFunc func(http.Handler) http.Handler

// Logging is a middleware function that logs information about incoming requests.
func RedisCaching(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//  Only cache GET requests
		if r.Method != "GET" {
			next.ServeHTTP(w, r)
			return
		}
		key := r.URL.Path
		cachedData, err := Client.Get(key).Result()
		if err == nil {
			// Data found in cache, serve it
			w.Write([]byte(cachedData))
			return
		}
		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		if err == redis.Nil {
			var response JSONResponse
			jsonData, err1 := json.Marshal(response)
			if err1 != nil {
				http.Error(w, "Data is not Json- Format", http.StatusInternalServerError)
				return
			}
			key := "DummyJson"
			err2 := Client.Set(key, jsonData, 1*time.Minute)
			if err2 != nil {
				fmt.Fprint(w, err2)
				return
			}
		}
	})
}
