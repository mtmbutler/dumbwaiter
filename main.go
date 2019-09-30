// A simple RESTful API based on this tutorial:
//  - https://tutorialedge.net/golang/creating-restful-api-with-golang/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Day struct {
	// date string `json:"date"`
	// amWeight  [3]float32 `json:"am_weight"`
	// pmWeight  [3]float32 `json:"pm_weight"`
	snack     int `json:"snack"`
	breakfast int `json:"breakfast"`
	lunch     int `json:"lunch"`
	dinner    int `json:"dinner"`
	exercise  int `json:"exercise"`
}

var Days []Day

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/days", returnAllDays)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func returnAllDays(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllDays")
	json.NewEncoder(w).Encode(Days)
}

func main() {
	Days = []Day{
		Day{
			// "2019-09-30",
			// [3]float32{180, 180.2, 178.6},
			// [3]float32{179.8, 179.4, 180.0},
			500,
			400,
			500,
			600,
			100,
		},
		Day{
			// "2019-10-01",
			// [3]float32{180, 180.2, 178.6},
			// [3]float32{179.8, 179.4, 180.0},
			600,
			300,
			400,
			500,
			300,
		},
	}
	fmt.Println("Starting")
	handleRequests()
}
