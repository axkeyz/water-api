// filtermodels.go contains the SQL WHERE query components.
// SQL WHERE statements filter results during SQL queries.
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

// SetSearchWhere adds a SQL WHERE that filters database
// records depending on user input. If the user input is a
// string or by outage id if the user input is an integer.
// The SQL WHERE statement is added to *Query.Wheres.
func (query *Query) SetSearchWhere(value []string) {
	for _, i := range value {
		if isInt(i) {
			query.SetOutageIDWhere(i)
		} else {
			query.SetAddressWhere(i)
		}
	}
}

// SetAddressWhere adds a SQL WHERE statement that filters
// database records of the given address in both the street
// and suburb columns. The SQL WHERE statement is added to
// *Query.Wheres.
func (query *Query) SetAddressWhere(address string) {
	// Search address in the street and suburb columns and
	// attempt to unabbreviate shorthands if applicable
	query.Wheres = append(
		query.Wheres, fmt.Sprintf(
			`(lower(suburb) LIKE lower('%%%s%%')
			OR lower(street) LIKE lower('%%%s%%') 
			OR lower(suburb) LIKE lower('%%%s%%')
			OR lower(street) LIKE lower('%%%s%%'))`,
			address, address,
			CleanAddressName(address, "suburb"),
			CleanAddressName(address, "street"),
		),
	)
}

// SetAddressWhere adds a SQL WHERE statement that filters
// database records of the given address in either the street
// or the suburb column. The SQL WHERE statement is added to
// *Query.Wheres.
func (query *Query) SetAddressOfTypeWhere(
	addressName, addressType string) {
	query.Wheres = append(
		// Search address in the given addressType column and
		// attempt to unabbreviate shorthands if applicable
		query.Wheres, fmt.Sprintf(
			`(lower(%s) LIKE lower('%%%s%%')
			OR lower(%s) LIKE lower('%%%s%%'))`,
			addressType, addressName, addressType,
			CleanAddressName(addressName, addressType),
		),
	)
}

// SetLocationRadiusWhere adds a SQL WHERE statement that
// filters database records of a radius circle (in m) around
// a longitude and latitude. The SQL WHERE statement is added
// to *Query.Wheres.
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

// SetOutageTypeWhere adds a SQL WHERE that filters database
// records by the outage_type column. The SQL WHERE statement
// is added to *Query.Wheres.
func (query *Query) SetOutageTypeWhere(outageType string) {
	query.SetSignedWhere(
		"outage_type = ", outageType,
	)
}

// SetOutageIDWhere adds a SQL WHERE that filters database
// records by the outage_id column. The SQL WHERE statement
// is added to *Query.Wheres.
func (query *Query) SetOutageIDWhere(value string) {
	query.SetSignedWhere(
		"outage_id = ", value,
	)
}

// SetAddressWhere sets a SQL WHERE statement that filters
// database records of the column with an SQL operation sign
// (such as >=, =, LIKE) with the value to be assigned. The
// SQL WHERE statement is added to *Query.Wheres.
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

// MakeOrderbyString returns an orderby string after adding
// the orderby strings to *Query.Orderbys. The orderby strings
// are in the format "column_name asc/desc"
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

// MakeOrderbyStringFromFields makes a single SQL ORDER BY
// statement by combining the *Query.Orderbys.
func (query *Query) MakeOrderbyStringFromFields() string {
	return " ORDER BY " + strings.Join(query.Orderbys, ", ")
}

// MakePaginationString makes an SQL limit-offset pagination
// string by using the *Query.Limit and *Query.Offset.
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
