// filtermodels.go contains the SQL query components.
package api

import (
	"fmt"
	"net/url"
	"strings"
)

type Query struct {
	Selects  []string
	Table    string
	Wheres   []string
	Orderbys []string
	Limit    string
	Offset   string
	GroupBy  []string
}

// SetSearchWhere adds a search by location name (street/suburb)
// or outage id depending on the value of user input.
func (query *Query) SetSearchWhere(value []string) {
	for _, i := range value {
		if isInt(i) {
			query.SetSearchIDWhere(i)
		} else {
			query.SetLocationWhere(i)
		}
	}
}

// SetSearchLocationWhere adds a search by location name
// (street/suburb) to the array of wheres.
func (query *Query) SetLocationWhere(location string) {
	query.Wheres = append(
		query.Wheres, fmt.Sprintf(
			`(lower(suburb) LIKE lower('%%%s%%')
			OR lower(street) LIKE lower('%%%s%%') 
			OR lower(suburb) LIKE lower('%%%s%%')
			OR lower(street) LIKE lower('%%%s%%'))`,
			location, location,
			CleanAddressName(location, "suburb"),
			CleanAddressName(location, "street"),
		),
	)
}

// SetOneLocationWhere adds a search by location name of
// one specific address type (street or suburb).
func (query *Query) SetOneLocationWhere(
	addressName, addressType string) {
	query.Wheres = append(
		query.Wheres, fmt.Sprintf(
			`(lower(%s) LIKE lower('%%%s%%')
			OR lower(%s) LIKE lower('%%%s%%'))`,
			addressType, addressName, addressType,
			CleanAddressName(addressName, addressType),
		),
	)
}

// SetSearchIDWhere adds a search by outage id to the
// array of wheres.
func (query *Query) SetSearchIDWhere(id string) {
	query.Wheres = append(query.Wheres,
		fmt.Sprintf(
			`lower(cast(outage_id as text)) 
			LIKE lower('%%%s%%')`, id,
		),
	)
}

// MakeWhereString combines the Wheres into a single SQL
// WHERE statement, joined by an " AND " or " OR " SQL
// condition.
func (query *Query) MakeWhereString(condition string) string {
	return " WHERE " + strings.Join(query.Wheres, condition)
}

// MakeOrderByString returns an orderby string after adding
// the orderby strings to query.Orderbys.
// The orderby strings are in the format "column_name asc/desc"
func (query *Query) MakeOrderbyString(
	orderbys []string,
) string {
	query.SetOrderbysField(orderbys)
	return query.MakeOrderbyStringFromFields()
}

// SetOrderbyFields adds strings in the format "column_name asc/desc"
// to the orderbys field.
func (query *Query) SetOrderbysField(orderbys []string) {
	query.Orderbys = append(query.Orderbys, orderbys...)
}

// MakeOrderString makes a single SQL ORDER BY statement by
// combining the Orderbys array.
func (query *Query) MakeOrderbyStringFromFields() string {
	return " ORDER BY " + strings.Join(query.Orderbys, ", ")
}

// MakePaginationString returns a limit-offset string after setting the
// query to corresponding limit and offset values.
func (query *Query) MakePaginationString(
	limit string, offset string,
) string {
	query.SetPaginationFields(limit, offset)
	return query.MakePaginationStringFromFields()
}

// Set pagination fields sets the limit and offset pagination fields.
func (query *Query) SetPaginationFields(
	limit string, offset string,
) {
	query.Limit = limit
	query.Offset = offset
}

// MakePaginationString makes an SQL limit-offset pagination string
// by using the query limit and offset.
func (query *Query) MakePaginationStringFromFields() (
	pagination string) {
	// Get sorting order (ascending / descending)
	if query.Limit != "" && query.Offset != "" {
		// Pagination string
		pagination = fmt.Sprintf(
			"LIMIT %s OFFSET %s", query.Limit, query.Offset,
		)
	}

	return
}

// SetSignedWhere sets a column with its operational sign and
// the value to be assigned.
// For example:
//		signedColumn = "outage_type ="
//		value = "Planned"
// Returns:
//		"outage_type = 'Planned'"
func (query *Query) SetSignedWhere(signedColumn, value string) {
	if value != "" {
		query.Wheres = append(
			query.Wheres,
			fmt.Sprintf("%s '%s'", signedColumn, value),
		)
	}
}

// SetLocationRadiusWhere adds a search by longitude, latitude and radius to
// the array of Wheres.
func (query *Query) SetLocationRadiusWhere(
	longitude, latitude, radius string,
) {
	if longitude != "" && latitude != "" && radius != "" {
		query.Wheres = append(
			query.Wheres,
			fmt.Sprintf(
				`ST_DWithin(location, 
				ST_SetSRID(ST_Point(%s, %s), 4326), %s)`,
				longitude, latitude, radius,
			),
		)
	}
}

func (query *Query) SetDateWheres(params url.Values) {
	for _, param := range DateColumns {
		if value := params.Get(param); len(value) > 0 {
			_, column := IsDateParam(param)
			query.SetSignedWhere(column, value)
		}
	}
}

func (query *Query) SetAllAddressTypeWheres(
	params url.Values, addressType string) {
	if values := params[addressType]; len(values) > 0 {
		for _, i := range values {
			query.SetOneLocationWhere(i, addressType)
		}
	}
}

func (query *Query) SetAllAddressWheres(params url.Values) {
	query.SetAllAddressTypeWheres(params, "street")
	query.SetAllAddressTypeWheres(params, "suburb")
}

func (query *Query) SetWheres(params url.Values) {
	query.SetSearchWhere(params["search"])

	query.SetDateWheres(params)
	query.SetAllAddressWheres(params)

	query.SetLocationRadiusWhere(
		params.Get("longitude"), params.Get("latitude"),
		params.Get("radius"),
	)

	query.SetSignedWhere(
		"outage_type = ", params.Get("outage_type"),
	)
}
