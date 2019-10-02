// A simple RESTful API based on this tutorial:
//  - https://tutorialedge.net/golang/creating-restful-api-with-golang/
package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Day struct {
	Id        string     `json:"id"`
	Date      string     `json:"date"`
	AMWeight  [3]float32 `json:"amWeight"`
	PMWeight  [3]float32 `json:"pmWeight"`
	Snack     int        `json:"snack"`
	Breakfast int        `json:"breakfast"`
	Lunch     int        `json:"lunch"`
	Dinner    int        `json:"dinner"`
	Exercise  int        `json:"exercise"`
}

var Days []Day

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	// Instantiate a mux router
	myRouter := mux.NewRouter().StrictSlash(true)

	// Add handles
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/days", returnAllDays)
	myRouter.HandleFunc("/days/{id}", returnSingleDay)

	// Run
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func returnAllDays(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllDays")
	json.NewEncoder(w).Encode(Days)
}

func returnSingleDay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	fmt.Printf("Endpoint Hit: returnSingleDay <%s>\n", key)

	// Find the corresponding day
	for _, day := range Days {
		if day.Id == key {
			json.NewEncoder(w).Encode(day)
		}
	}
}

func main() {
	Days = []Day{
		Day{
			Id:        "1",
			Date:      "2019-09-30",
			AMWeight:  [3]float32{180, 180.2, 178.6},
			PMWeight:  [3]float32{179.8, 179.4, 180.0},
			Snack:     500,
			Breakfast: 400,
			Lunch:     500,
			Dinner:    600,
			Exercise:  100,
		},
		Day{
			Id:        "2",
			Date:      "2019-10-01",
			AMWeight:  [3]float32{179.8, 179.2, 179.6},
			PMWeight:  [3]float32{179.2, 179.8, 180.2},
			Snack:     400,
			Breakfast: 500,
			Lunch:     700,
			Dinner:    500,
			Exercise:  200,
		},
	}
	fmt.Println("Starting")
	handleRequests()
}
