// filters_maps.go contains variables/constants that are used to
// store data and key-value pairs.
package api

var FilterableParams = []string{
	"suburb", "street", "outage_type", "search",
	"before_start_date", "before_end_date", "after_end_date",
	"after_start_date", "location", "outage_id",
}

var FilterableCountParams = []string{
	"total_hours", "total_outages",
}

var SQLSigns = map[string]string{
	"after":  ">=",
	"before": "<=",
	"like":   "LIKE",
}

var DateColumns = []string{
	"before_end_date", "after_end_date",
	"before_start_date", "after_start_date",
}
