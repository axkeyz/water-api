// filtermodels.go contains the SQL query components.
package api

import (
	"fmt"
	"strings"
)

type Query struct {
	Selects  []string
	Wheres   []string
	Orderbys []string
	GroupBy  []string
}

// MakeWhereString combines the Wheres into a single SQL
// WHERE statement, joined by an " AND " or " OR " SQL
// condition.
func (query *Query) MakeWhereString(condition string) string {
	return " WHERE " + strings.Join(query.Wheres, condition)
}

// SetSearchWhere adds a search by location name (street/suburb)
// or outage id depending on the value of user input.
func (query *Query) SetSearchWhere(value []string) {
	for _, i := range value {
		if isInt(i) {
			query.SetOutageIDWhere(i)
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

// SetLocationRadiusWhere adds a SQL WHERE condition for longitude,
// the latitude and radius parameters.
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

// SetOutageIDWhere adds a SQL WHERE condition by outage_type.
func (query *Query) SetOutageTypeWhere(value string) {
	query.SetSignedWhere(
		"outage_type = ", value,
	)
}

// SetOutageIDWhere adds a SQL WHERE condition by outage_id.
func (query *Query) SetOutageIDWhere(value string) {
	query.SetSignedWhere(
		"outage_id = ", value,
	)
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

// MakePaginationString makes an SQL limit-offset pagination string
// by using the query limit and offset.
func (query *Query) MakePaginationString(
	limit, offset string) (pagination string) {
	if limit != "" && offset != "" {
		// Pagination string
		pagination = fmt.Sprintf(
			"LIMIT %s OFFSET %s", limit, offset,
		)
	}

	return
}
