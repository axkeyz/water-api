package main

import (
    "github.com/axkeyz/water-down-again/api"
)

func main() {
	outages := api.GetAPIData()
	api.WriteOutage(outages)
}