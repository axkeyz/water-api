// models.go contains the database models from the official API 
// & the database of this app.

package api

// A WaterOutage struct maps a water outage from the Watercare API instance.
type WaterOutage struct {
	OutageID int64 `json:"outageId"`
	Location string `json:"location"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	StartDate string `json:"startDate"`
	EndDate string `json:"endDate"`
	OutageType string `json:"outageType"`
}

// A DBWaterOutage struct maps a water outage from the database of this app.
type DBWaterOutage struct {
	OutageID int `json:"outage_id"`
	Street string `json:"street"`
	Suburb string `json:"suburb"`
	Location string `json:"location"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	OutageType string `json:"outage_type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// A DBCountOutage struct maps count-based queries from the database of this app.