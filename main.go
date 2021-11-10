package main

import (
    "github.com/axkeyz/water-down-again/api"
	"github.com/robfig/cron"
	"fmt"
)

func UpdateOutages() {
	outages := api.GetAPIData()
	api.WriteOutage(outages)
	fmt.Printf("Done.\n")
}

func main() {
	// Retrieve & write from Watercare API to this app's database for the first time
	UpdateOutages()

	// Create cronjob to run retrieval & write every 45 minutes
	c := cron.New()
    c.AddFunc("@every 45m", UpdateOutages)
    c.Start()
    select {}
}