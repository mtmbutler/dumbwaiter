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

type Day struct {
	gorm.Model
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
		"host=%s port=%s user=%s dbname=%s password=%s",
		host, port, user, name, pass,
	)

	// Open a database connection
	fmt.Printf("Opening database connection: %s\n", dbUrl)
	var err error
	DB, err = gorm.Open("postgres", dbUrl)
	if err != nil {
		panic("Connection failed")
	} else {
		fmt.Println("Connection successful")
	}
	defer DB.Close()
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

	// Run
	fmt.Println("Listening")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

// getKey parses a uint primary key from a request body
func getKey(vars map[string]string) uint {
	key := vars["id"]
	id, _ := strconv.Atoi(key)
	return uint(id)
}

func returnAllDays(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllDays")
	days := make([]*Day, 0)
	DB.Find(&days)
	json.NewEncoder(w).Encode(days)
}

func returnSingleDay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := getKey(vars)
	fmt.Printf("Endpoint Hit: returnSingleDay <%d>\n", key)

	// Find the corresponding day
	var day Day
	DB.First(&day, key)
	if day.ID == key {
		json.NewEncoder(w).Encode(day)
	}
}

func createNewDay(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewDay")

	// Read the request, create an object and add it to the database
	reqBody, _ := ioutil.ReadAll(r.Body)
	var day Day
	json.Unmarshal(reqBody, &day)
	DB.Create(&day)

	// Respond with the new object
	json.NewEncoder(w).Encode(day)
}

func deleteDay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := getKey(vars)
	fmt.Printf("Endpoint Hit: deleteDay <%d>\n", key)

	// Find the corresponding day
	var day Day
	DB.First(&day, key)
	if day.ID == key {
		DB.Delete(&day)
	}
}

func updateDay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := getKey(vars)
	fmt.Printf("Endpoint Hit: updateDay <%d>\n", key)

	// Find the object
	var day Day
	DB.First(&day, key)

	if day.ID == key {
		// Read the request and update the object
		reqBody, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(reqBody, &day)
		DB.Save(&day)

		// Respond with the updated object
		json.NewEncoder(w).Encode(day)
	}
}

func main() {
	fmt.Println("Starting")
	handleRequests()
}
