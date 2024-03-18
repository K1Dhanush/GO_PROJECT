package main

import (
	auth "GORM/BasicMiddleware"
	logging "GORM/MiddleWare"

	cache "GORM/REDIS"
	jwt "GORM/Token"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Gobal Variable

// model
type Product struct {
	gorm.Model           //To Create additional Fields && Model is a --struct
	Item         string  `json:"item"`
	Price        float32 `json:"price"`
	ReturnPolicy int     `json:"return_policy"`
}

var db *gorm.DB

// var sessionStore = sessions.NewCookieStore([]byte("your-secret-key"))

var Client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379", //"redis:6379",
	Password: "",
	DB:       0, //default database
})

func main() {
	//driver to connect with MySQL
	//parseTime is used to send the current --- the created , updated timing to the GORM
	// automatically parse time values
	dsn := "root:Sathyabama*40110529@tcp(host.docker.internal:3306)/product?parseTime=true"
	// dsn := "root:Sathyabama*40110529@tcp(localhost:3306)/product?parseTime=true"

	//error is interface
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Print(err.Error())
		panic("Unable to Connect")
	}

	//Creating the Table in Schema
	errr := db.AutoMigrate(&Product{})
	if errr != nil {
		fmt.Println("Error migrating database:", errr)
		panic("Unable to migrate database")
	}

	//Setting the client
	cache.SetRedisClient(Client)

	//For testing the redis Connection.
	_, err2 := Client.Ping().Result()
	if err2 != nil {
		panic(err2)
	}

	r := mux.NewRouter()
	//for logging MiddleWare
	r.Use(logging.Logging)

	//for Basic MiddleWare
	r.Use(auth.NewAuth)

	//for redisCaching MiddleWare requests for only get--method.
	r.Use(cache.RedisCaching)

	//Creating RESTAPI's
	// r.HandleFunc("/login", loginSession)
	r.HandleFunc("/addProduct", addProduct).Methods("POST")
	r.HandleFunc("/getAllProducts", getAllProducts).Methods("GET")
	r.HandleFunc("/getProduct/{id}", getProduct).Methods("GET")
	r.HandleFunc("/updateProduct/{id}", updateProduct).Methods("PUT")
	r.HandleFunc("/deleteProduct/{id}", deleteProduct).Methods("DELETE")
	// r.HandleFunc("/logout", logout)
	r.HandleFunc("/jwt-token", jwt.Token)

	//Activating the Server
	// http.ListenAndServe(":8081", r)
	http.ListenAndServe("0.0.0.0:8082", r)
}

// // Session Management.
// func loginSession(w http.ResponseWriter, r *http.Request) {
// 	username := r.URL.Query().Get("username")
// 	if username == "" {
// 		fmt.Fprint(w, "Username is Required")
// 		return
// 	}
// 	//retrieve a session associated with a given HTTP request (r)
// 	session, err := sessionStore.Get(r, "session-name")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	//storing the data into the sessionStore
// 	session.Values["username"] = username
// 	session.Values["authentication"] = true

// 	errr := session.Save(r, w)
// 	if errr != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	fmt.Fprint(w, "Login Successfully")
// }
// func logout(w http.ResponseWriter, r *http.Request) {
// 	session, err := sessionStore.Get(r, "session-name")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	session.Values["authentication"] = false
// 	session.Options.MaxAge = -1 //for clearing the session
// 	errr := session.Save(r, w)
// 	if errr != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	fmt.Fprint(w, "Logout Successfully")
// }

func addProduct(w http.ResponseWriter, r *http.Request) {
	//Create a Var
	var prod Product

	//requesting to Give response in JSON format
	w.Header().Add("Content-Type", "application/json")

	//JSON-to-Struct
	decode := json.NewDecoder(r.Body) //reads the data which is present in JSON
	err := decode.Decode(&prod)       //Storing in a variable
	if err != nil {
		http.Error(w, "Enter the data in Correct Format", http.StatusNotAcceptable)
		return
	}
	if db == nil {
		http.Error(w, "Database connection is not initialized", http.StatusInternalServerError)
		return
	}

	// Save the product
	result := db.Save(&prod)
	if result.Error != nil {
		http.Error(w, "Failed to save product", http.StatusInternalServerError)
		return
	}

	// Respond with the saved product
	json.NewEncoder(w).Encode(prod)

}

func getAllProducts(w http.ResponseWriter, r *http.Request) {
	//Slice
	var prod []Product
	// retrieve records from the database that match certain criteria and Store in &prod --- instance of a struct
	_ = db.Find(&prod)

	//requesting to Give response in JSON format
	w.Header().Add("Content-Type", "application/json")

	//Encode
	json.NewEncoder(w).Encode(prod) //Encodes -- Slice data str to JSON Format
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	var prod Product

	vars := mux.Vars(r)
	id := vars["id"]

	db.Where("id=?", id).First(&prod) //Here "i=?"d should be same as in Table because it is PK

	//requesting to Give response in JSON format
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prod)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	//variables from the request
	vars := mux.Vars(r)
	id := vars["id"]

	var existingProd Product
	//to check Whether the product is present are not and Storing in a new varible
	result := db.Where("id=?", id).First(&existingProd)
	if result.RowsAffected == 0 { //DB is a struct which contains ERROR
		http.Error(w, "Product is not Present", http.StatusNotFound)
		return
	}

	//var
	var update Product
	//Decoding and storing in a new Variable
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, "Not in JSON FORMAT", http.StatusBadRequest)
		return
	}

	log.Println(existingProd)

	//Updating
	// existingProd.Model = update.Model
	existingProd.Item = update.Item
	existingProd.Price = update.Price
	existingProd.ReturnPolicy = update.ReturnPolicy
	log.Println(update)

	res := db.Save(&existingProd)
	if res.Error != nil { //DB is a struct which contains ERROR
		http.Error(w, "Unable to update the Product", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	//For Output
	json.NewEncoder(w).Encode(existingProd)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var prod Product
	result := db.Where("id=?", id).First(&prod)
	if result.RowsAffected == 0 {
		http.Error(w, "Product is not present", http.StatusNotFound)
		return
	}
	res := db.Delete(&prod)
	if res.Error != nil {
		http.Error(w, "Unable to delete Product", http.StatusInternalServerError)
		return
	}

	//Response
	json.NewEncoder(w).Encode(&prod)
}
