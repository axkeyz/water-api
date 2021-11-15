package main

import (
	"github.com/axkeyz/water-down-again/api"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("Server is running")
	
	// Create a cronjob for every hour
	go func() {
		for {
			api.UpdateOutages()
			<-time.After(1 * time.Hour)
		}
	}()
	
	// Retrieve & write from Watercare API to this app's database for the first time
	api.UpdateOutages()

	// Init the mux router
	router := mux.NewRouter()

	// Setup routes
	router.HandleFunc("/", api.GetOutages).Methods("GET")
	router.HandleFunc("/count", api.CountOutages).Methods("GET")

	// Run server
	log.Println(http.ListenAndServe(":8080", router))
}