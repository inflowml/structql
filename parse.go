package structql

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/inflowml/logger"
)

// parseResponse parses the given SQL rows into a slice of structures with the
// provided object type.  For example, consider the Person structure below:
//
//  type Person struct {
//    ID   int32  `sql:"id"`
//    Name string `sql:"name"`
//    Age  int32  `sql:"age"`
//  }
//
// Each field of the structure is annotated with a tag which maps that field to
// a column in the database table.  This enables an *sql.Rows object to be
// converted into a slice of Person structures with the following code:
//
//  people, err := parseResponse(rows, Person{})
//
// Given that people has type []interface{}, it is necessary to cast an entry of
// people into a Person object before accessing a member of that Person object.
func parseResponse(rows *sql.Rows, object interface{}) ([]interface{}, error) {
	template := reflect.TypeOf(object)

	// Verify that the object is a structure.
	if template.Kind() != reflect.Struct {
		return []interface{}{}, fmt.Errorf("type %T is not a structure", object)
	}

	// Construct a map that associates the name of a column with the name of a field.
	ctfMap := map[string]string{}

	// Populate the map column-to-field map using the template.
	for i := 0; i < template.NumField(); i++ {
		field := template.Field(i)
		if col, ok := field.Tag.Lookup("sql"); ok {
			ctfMap[col] = field.Name
		}
	}

	// Get the names and types of the columns.
	colNames, err := rows.Columns()
	if err != nil {
		return []interface{}{}, fmt.Errorf("failed to get column names: %v", err)
	}
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return []interface{}{}, fmt.Errorf("failed to get column types: %v", err)
	}

	// Construct a slice of suitable arguments to (*sql.Rows).Scan().
	entries := []interface{}{}
	for _, colType := range colTypes {
		var scanType reflect.Type

		// The PostgreSQL library does not support floating-point column types.
		// See: https://github.com/lib/pq/issues/761
		switch colType.DatabaseTypeName() {
		case "FLOAT4":
			scanType = reflect.TypeOf(float32(0))
		case "FLOAT8":
			scanType = reflect.TypeOf(float64(0))
		default:
			scanType = colType.ScanType()
		}

		entry := reflect.New(scanType).Interface()
		entries = append(entries, entry)
	}

	// Construct a slice to hold the converted entries of each row.
	vessels := []interface{}{}

	// Loop over the rows.
	for rows.Next() {
		// Scan the current row into the entry slice.
		rows.Scan(entries...)

		// Construct a vessel to hold the entries.
		vessel := reflect.New(template).Elem()

		// Loop over the entries.
		for i := range entries {
			// Find the column and field name associated with the current entry.
			colName := colNames[i]
			fieldName, ok := ctfMap[colName]
			if !ok {
				logger.Warning("No field in structure %T is tagged with SQL column %q.", template, colName)
				continue
			}

			// Populate a field from the vessel with the contents of the entry.
			field := vessel.FieldByName(fieldName)
			entry := reflect.ValueOf(entries[i]).Elem()
			field.Set(entry)
		}
		vessels = append(vessels, vessel.Interface())
	}
	return vessels, nil
}
