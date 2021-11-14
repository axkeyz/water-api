package main

import (
	"github.com/axkeyz/water-down-again/api"
	"github.com/robfig/cron"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	// Retrieve & write from Watercare API to this app's database for the first time
	api.UpdateOutages()

	// Init the mux router
	router := mux.NewRouter()

	// Setup routes
	router.HandleFunc("/", api.GetOutages).Methods("GET")
	router.HandleFunc("/count", api.CountOutages).Methods("GET")

	// Run server
	log.Println(http.ListenAndServe(":8080", router))
	log.Println("Server is running")

	// Create cronjobs
	c := cron.New()
	// Retrieve Watercare API & update app database
	c.AddFunc("@hourly", api.UpdateOutages)

	c.Start()
 	select {}
}