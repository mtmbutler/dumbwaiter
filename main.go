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
	Email   string `json:"email"`
	ApiKey  string `json:"apiKey";gorm:"unique;not null"`
	IsAdmin bool   `json:"isAdmin"`
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

// createAdmin creates a superuser "admin@dumbwaiter" with an API key defined by the
// env variable ADMIN_KEY. If "admin@dumbwaiter" already exists, its API key is set
// to ADMIN_KEY.
func createAdmin() {
	fmt.Println("Checking admin")
	var admin User
	DB.FirstOrInit(&admin, User{email: "admin@dumbwaiter"})
	admin.apiKey = os.Getenv("ADMIN_KEY")
	DB.Save(&admin)
}

// getUser hits the DB for the correct user based on the apiKey in the request body.
// If there isn't a match, returns nil.
func getUser(r *http.Request) *User {
	apiKey := r.URL.Query().Get("apiKey")
	fmt.Println(apiKey)
	var user User
	DB.Where("api_key = ?", apiKey).First(&user)
	if user.ApiKey == apiKey {
		return &user
	} else {
		return nil
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
	user := getUser(r)
	if user != nil {
		days := make([]*Day, 0)
		DB.Where("user_id = ?", user.ID).Find(&days)
		json.NewEncoder(w).Encode(days)
	}
}

func returnSingleDay(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	key := getKey(r)
	fmt.Printf("Endpoint Hit: returnSingleDay <%d>\n", key)

	// Find the corresponding day
	if user != nil {
		var day Day
		DB.Where("user_id = ?", user.ID).First(&day, key)
		if day.ID == key {
			json.NewEncoder(w).Encode(day)
		}
	}
}

func createNewDay(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	fmt.Println("Endpoint Hit: createNewDay")

	if user != nil {
		// Read the request, create an object and add it to the database
		reqBody, _ := ioutil.ReadAll(r.Body)
		var day Day
		json.Unmarshal(reqBody, &day)
		day.UserID = int(user.ID)
		DB.Create(&day)

		// Respond with the new object
		json.NewEncoder(w).Encode(day)
	}
}

func deleteDay(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	key := getKey(r)
	fmt.Printf("Endpoint Hit: deleteDay <%d>\n", key)

	if user != nil {
		// Find the corresponding day
		var day Day
		DB.Where("user_id = ?", user.ID).First(&day, key)
		if day.ID == key {
			DB.Delete(&day)
		}
	}
}

func updateDay(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	key := getKey(r)
	fmt.Printf("Endpoint Hit: updateDay <%d>\n", key)

	if user != nil {
		// Find the object
		var day Day
		DB.Where("user_id = ?", user.ID).First(&day, key)

		if day.ID == key {
			// Read the request and update the object
			reqBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(reqBody, &day)
			day.UserID = int(user.ID)
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
