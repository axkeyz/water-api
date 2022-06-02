package main

import (
	"log"
	"net/http"
	"time"

	"github.com/axkeyz/water-down-again/api"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("Server is running")

	// Create a cronjob for every hour to retrieve & write from Watercare API to this
	// app's database
	go func() {
		for {
			api.UpdateOutages()
			<-time.After(1 * time.Hour)
		}
	}()

	// Reformat street and suburb of old outages from previous builds
	api.CleanupOutages()

	// Init the mux router
	router := mux.NewRouter()
	router.Use(mux.CORSMethodMiddleware(router))

	// Setup routes
	router.HandleFunc("/", api.GetOutages).Methods("GET")
	router.HandleFunc("/count", api.CountOutages).Methods("GET")

	// Run server
	log.Println(http.ListenAndServe(":8080", router))
}
