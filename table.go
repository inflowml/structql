package structql

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/inflowml/logger"
)

// Count mimics the PSQL format for Count responses in a Go Struct
type Count struct {
	Count int64 `sql:"count"`
}

// CreateTableFromObject creates an SQL table with the given name from the type
// of the provided object using the Connection receiver.  The provided object
// must be a structure where:
//  1. Each field is annotated with an "sql" tag.  Fields may also be annotated
//     with a "typ" or "opt" tag to override their default settings:
//     - The "sql" tag denotes the column name (e.g., "id").
//     - The "typ" tag denotes the column type (e.g., "SERIAL").
//     - The "opt" tag denotes column constraints (e.g., "PRIMARY KEY").
//  2. One field must be a 32-bit integer that corresponds to the "id" column.
func (conn *Connection) CreateTableFromObject(table string, object interface{}) error {
	template := reflect.TypeOf(object)

	// Verify that the object is a structure.
	if template.Kind() != reflect.Struct {
		return fmt.Errorf("type %s is not a structure", template.Name())
	}

	// match reports whether the field with the provided name is an ID column.
	match := func(name string) bool {
		field, _ := template.FieldByName(name)
		return field.Tag.Get("sql") == "id"
	}

	// Verify that the object contains an ID column.
	if _, ok := template.FieldByNameFunc(match); !ok {
		return fmt.Errorf("structure %s does not have a field for the ID column", template.Name())
	}

	// Construct a slice that holds the SQL table headers.
	headers := make([]string, 0, template.NumField())

	for i := 0; i < template.NumField(); i++ {
		field := template.Field(i)

		// Derive the name of the SQL column.
		sql, ok := field.Tag.Lookup("sql")
		if !ok {
			logger.Warning("Field %q in structure %s does not have an SQL column tag.", field.Name, template.Name())
			continue
		}

		// Derive the type of the SQL column.
		typ, err := getColumnType(field)
		if err != nil {
			logger.Warning("Field %q in structure %s does not have a PostgreSQL type: %v.", field.Name, template.Name(), err)
			continue
		}

		// Construct a column header from the column name, type, and constraints.
		header := fmt.Sprintf("%s %s %s", sql, typ, field.Tag.Get("opt"))
		headers = append(headers, header)
	}

	// Create the table (if it does not already exist).
	schema := strings.Join(headers, ", ")
	stmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", table, schema)
	logger.SQL(stmt)
	_, err := conn.exec(stmt)
	if err == nil {
		logger.SQL("Table created successfully ")
	}
	return err
}

// getColumnType derives the PostgreSQL type of the given structure field.
func getColumnType(field reflect.StructField) (string, error) {
	if typ, ok := field.Tag.Lookup("typ"); ok {
		return typ, nil
	}

	var typ string
	switch field.Type {
	case reflect.TypeOf(false):
		typ = "BOOL"
	case reflect.TypeOf(int16(0)):
		typ = "INT2"
	case reflect.TypeOf(int32(0)), reflect.TypeOf(int(0)):
		typ = "INT4"
	case reflect.TypeOf(int64(0)):
		typ = "INT8"
	case reflect.TypeOf(float32(0)):
		typ = "FLOAT4"
	case reflect.TypeOf(float64(0)):
		typ = "FLOAT8"
	case reflect.TypeOf(""):
		typ = "TEXT"
	case reflect.TypeOf(time.Time{}):
		typ = "TIMESTAMP"
	case reflect.TypeOf([]byte{}):
		typ = "BYTEA"
	default:
		return "", fmt.Errorf("type %s is not supported", field.Type)
	}
	return typ, nil
}

// DropTable drops the given table from the Connection receiver.
func (conn *Connection) DropTable(table string) error {
	stmt := fmt.Sprintf("DROP TABLE IF EXISTS %s;", table)
	_, err := conn.exec(stmt)
	return err
}

// CountRows accepts a table name and returns the number of rows in that table.
func (conn *Connection) CountRows(table string) (int64, error) {
	stmt := fmt.Sprintf("SELECT COUNT(*) FROM %s;", table)

	rows, err := conn.query(stmt)
	if err != nil {
		return 0, fmt.Errorf("failed to get row count for table %x: %v", table, err)
	}

	cnt, err := parseResponse(rows, Count{})
	if err != nil {
		return 0, fmt.Errorf("failed to parse count response: %v", err)
	}

	// Extract Count from parsed struct
	val := cnt[0].(Count).Count

	return val, nil
}

// CountRowsWhere accepts a table name and condition statement
// and returns the number of rows in that table that meet the condition
func (conn *Connection) CountRowsWhere(table string, cond string) (int64, error) {
	stmt := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s;", table, cond)

	rows, err := conn.query(stmt)
	if err != nil {
		return 0, fmt.Errorf("failed to get row count for table %x: %v", table, err)
	}

	cnt, err := parseResponse(rows, Count{})
	if err != nil {
		return 0, fmt.Errorf("failed to parse count response: %v", err)
	}

	// Extract Count from parsed struct
	val := cnt[0].(Count).Count

	return val, nil
}

// OldestEntry returns the oldest row in the given table
func (conn *Connection) OldestEntry(object interface{}, table string, timestampCol string) (interface{}, error) {

	stmt := fmt.Sprintf("SELECT * FROM %s ORDER BY %s LIMIT 1", table, timestampCol)

	rows, err := conn.query(stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to get oldest entry count for table %x sorting by %s: %v", table, timestampCol, err)
	}

	objects, err := parseResponse(rows, object)
	if err != nil {
		return nil, fmt.Errorf("failed to parse rows: %v", err)
	}

	if len(objects) < 1 {
		return nil, fmt.Errorf("no values in cache")
	}

	return objects[0], nil

}
