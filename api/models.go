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
	OutageID int `json:"outage_id,omitempty"`
	Street string `json:"street,omitempty"`
	Suburb string `json:"suburb,omitempty"`
	Location string `json:"location,omitempty"`
	OutageType string `json:"outage_type,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate string `json:"end_date,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	TotalOutages int `json:"total_outages,omitempty"`
	TotalHours float64 `json:"total_hours,omitempty"`
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