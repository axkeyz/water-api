// models.go contains the database models from the official API 
// & the database of this app. It also contains some mapping functions.

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
	OutageID int `json:"Outage ID,omitempty"`
	Street string `json:"Street,omitempty"`
	Suburb string `json:"Suburb,omitempty"`
	Location string `json:"Location,omitempty"`
	StartDate string `json:"Start Date,omitempty"`
	EndDate string `json:"End Date,omitempty"`
	OutageType string `json:"Outage Type,omitempty"`
	CreatedAt string `json:"Created At,omitempty"`
	UpdatedAt string `json:"Updated At,omitempty"`
	TotalOutages int `json:"Total Outages,omitempty"`
	TotalHours float64 `json:"Total Hours,omitempty"`
}

// DBWaterOutageCol returns a reference for a column of a DBWaterOutage
func DBWaterOutageCol(colname string, outage *DBWaterOutage) interface{} {
    switch colname {
		case "outage_id":
			return &outage.OutageID
		case "street":
			return &outage.Street
		case "suburb":
			return &outage.Suburb
		case "location":
			return &outage.Location
		case "start_date":
			return &outage.StartDate
		case "end_date":
			return &outage.EndDate
		case "outage_type":
			return &outage.OutageType
		case "created_at":
			return &outage.CreatedAt
		case "updated_at":
			return &outage.UpdatedAt
		case "total_outages":
			return &outage.TotalOutages
		case "total_hours":
			return &outage.TotalHours
		default:
			panic("unknown column " + colname)
    }
}

// An AppError struct maps an error for this app.
type AppError struct {
	ErrorCode int64 `json:"Error Code"`
	Message string `json:"Message"`
	Details string `json:"Details"`
}