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
	router.HandleFunc("/api", api.GetOutages).Methods("GET")

	// Run server
	log.Fatal(http.ListenAndServe("localhost:8553", router))

	// Create cronjob to run retrieval & write every 45 minutes
	c := cron.New()
    c.AddFunc("@every 45m", api.UpdateOutages)
    c.Start()
    select {}
}