// A simple RESTful API based on this tutorial:
//  - https://tutorialedge.net/golang/creating-restful-api-with-golang/
package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type User struct {
	gorm.Model
	Email  string `json:"email"`
	ApiKey string `json:"apiKey";gorm:"unique;not null"`
}

type Day struct {
	gorm.Model
	UserID    int
	User      User
	Date      string  `json:"date"`
	AMWeight  float32 `json:"amWeight"`
	PMWeight  float32 `json:"pmWeight"`
	Snack     uint    `json:"snack"`
	Breakfast uint    `json:"breakfast"`
	Lunch     uint    `json:"lunch"`
	Dinner    uint    `json:"dinner"`
	Exercise  uint    `json:"exercise"`
}

var DB *gorm.DB

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	// Parse environment variables for DB auth
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	name := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PASS")
	dbUrl := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, name, pass,
	)

	// Open a database connection
	fmt.Printf("Opening database connection to %s\n", name)
	var err error
	DB, err = gorm.Open("postgres", dbUrl)
	if err != nil {
		fmt.Println(err)
		panic("Connection failed")
	} else {
		fmt.Println("Connection successful")
	}
	defer DB.Close()
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Day{})

	// Instantiate a mux router
	myRouter := mux.NewRouter().StrictSlash(true)

	// Add handles
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/days", returnAllDays).Methods("GET")
	myRouter.HandleFunc("/days", createNewDay).Methods("POST")
	myRouter.HandleFunc("/days/{id}", updateDay).Methods("PUT")
	myRouter.HandleFunc("/days/{id}", deleteDay).Methods("DELETE")
	myRouter.HandleFunc("/days/{id}", returnSingleDay).Methods("GET")
	myRouter.HandleFunc("/users", returnAllUsers).Methods("GET")
	myRouter.HandleFunc("/users", createNewUser).Methods("POST")
	myRouter.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	myRouter.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	myRouter.HandleFunc("/users/{id}", returnSingleUser).Methods("GET")

	// Run
	fmt.Println("Listening")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

// getKey parses a uint primary key from a request body
func getKey(r *http.Request) uint {
	vars := mux.Vars(r)
	fmt.Println(vars)
	key := vars["id"]
	id, _ := strconv.Atoi(key)
	return uint(id)
}

// getUserID hits the DB for the correct user ID based on the apiKey in the request
// body. If there isn't a match, returns 0.
func getUserID(r *http.Request) uint {
	apiKey := r.URL.Query().Get("apiKey")
	fmt.Println(apiKey)
	var user User
	DB.Where("api_key = ?", apiKey).First(&user)
	if user.ApiKey == apiKey {
		return user.ID
	} else {
		return 0
	}
}

func returnAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllUsers")
	users := make([]*User, 0)
	DB.Find(&users)
	json.NewEncoder(w).Encode(users)
}

func returnSingleUser(w http.ResponseWriter, r *http.Request) {
	key := getKey(r)
	fmt.Printf("Endpoint Hit: returnSingleUser <%d>\n", key)

	// Find the corresponding user
	var user User
	DB.First(&user, key)
	if user.ID == key {
		json.NewEncoder(w).Encode(user)
	}
}

func createNewUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewUser")

	// Read the request, create an object and add it to the database
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user User
	json.Unmarshal(reqBody, &user)
	DB.Create(&user)

	// Respond with the new object
	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	key := getKey(r)
	fmt.Printf("Endpoint Hit: deleteUser <%d>\n", key)

	// Find the corresponding user
	var user User
	DB.First(&user, key)
	if user.ID == key {
		DB.Delete(&user)
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	key := getKey(r)
	fmt.Printf("Endpoint Hit: updateUser <%d>\n", key)

	// Find the object
	var user User
	DB.First(&user, key)

	if user.ID == key {
		// Read the request and update the object
		reqBody, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(reqBody, &user)
		DB.Save(&user)

		// Respond with the updated object
		json.NewEncoder(w).Encode(user)
	}
}

func returnAllDays(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllDays")
	userID := getUserID(r)
	if userID != 0 {
		days := make([]*Day, 0)
		DB.Where("user_id = ?", userID).Find(&days)
		json.NewEncoder(w).Encode(days)
	}
}

func returnSingleDay(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	key := getKey(r)
	fmt.Printf("Endpoint Hit: returnSingleDay <%d>\n", key)

	// Find the corresponding day
	if userID != 0 {
		var day Day
		DB.Where("user_id = ?", userID).First(&day, key)
		if day.ID == key {
			json.NewEncoder(w).Encode(day)
		}
	}
}

func createNewDay(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	fmt.Println("Endpoint Hit: createNewDay")

	if userID != 0 {
		// Read the request, create an object and add it to the database
		reqBody, _ := ioutil.ReadAll(r.Body)
		var day Day
		json.Unmarshal(reqBody, &day)
		day.UserID = int(userID)
		DB.Create(&day)

		// Respond with the new object
		json.NewEncoder(w).Encode(day)
	}
}

func deleteDay(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	key := getKey(r)
	fmt.Printf("Endpoint Hit: deleteDay <%d>\n", key)

	if userID != 0 {
		// Find the corresponding day
		var day Day
		DB.Where("user_id = ?", userID).First(&day, key)
		if day.ID == key {
			DB.Delete(&day)
		}
	}
}

func updateDay(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	key := getKey(r)
	fmt.Printf("Endpoint Hit: updateDay <%d>\n", key)

	if userID != 0 {
		// Find the object
		var day Day
		DB.Where("user_id = ?", userID).First(&day, key)

		if day.ID == key {
			// Read the request and update the object
			reqBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(reqBody, &day)
			day.UserID = int(userID)
			DB.Save(&day)

			// Respond with the updated object
			json.NewEncoder(w).Encode(day)
		}
	}
}

func main() {
	fmt.Println("Starting")
	handleRequests()
}
