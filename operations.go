package main

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/inflowml/logger"
	"github.com/inflowml/structql/dataforms"
)

// SelectFrom executes a SELECT FROM query on the Connection receiver over the
// given object type and table.
func (conn *Connection) SelectFrom(object interface{}, table string) ([]interface{}, error) {
	return conn.executeSelect(object, table, "")
}

// SelectFromWhere executes a SELECT FROM WHERE query on the Connection receiver
// over the given object type, table, and conditional.  Additional arguments are
// substituted into the conditional in a style similar to printf().
func (conn *Connection) SelectFromWhere(object interface{}, table string, cond string, args ...interface{}) ([]interface{}, error) {
	return conn.executeSelect(object, table, cond, args...)
}

// executeSelect executes a SELECT FROM WHERE query on the Connection receiver
// over the given object, table, and conditional.  Setting the conditional
// to "" indicates that no conditional is desired.  Additional arguments are
// substituted into the conditional in a style similar to printf().
func (conn *Connection) executeSelect(object interface{}, table string, cond string, args ...interface{}) ([]interface{}, error) {
	// TODO: SQL Sanitization
	template := reflect.TypeOf(object)

	// Construct a slice that holds the SQL column name of each object field.
	cols := make([]string, 0, template.NumField())
	for i := 0; i < template.NumField(); i++ {
		field := template.Field(i)
		if col, ok := field.Tag.Lookup("sql"); ok {
			cols = append(cols, col)
		}
	}

	// Format the columns into a comma-separated list.
	colJoin := strings.Join(cols, ", ")

	// Translate the columns, table, and conditional into an SQL statement.
	stmt := fmt.Sprintf("SELECT %s FROM %s;", colJoin, table)
	if cond != "" {
		stmt = fmt.Sprintf("SELECT %s FROM %s WHERE %s;", colJoin, table, cond)
		stmt = fmt.Sprintf(stmt, args...)
	}

	// Execute the query on the SQL database.
	rows, err := conn.query(stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %v", stmt, err)
	}

	// Parse the rows from the query into a slice of Go objects based on the prototype.
	return parseResponse(rows, object)
}

//SelectForUpdate is to be used in conjuction with Lock and Unlock to facilitate row locking
//TODO add this to standard select but use arguments instead of new function
func (conn *Connection) SelectForUpdate(object interface{}, table string, cond string, args ...interface{}) ([]interface{}, error) {
	// TODO: SQL Sanitization
	template := reflect.TypeOf(object)

	// Construct a slice that holds the SQL column name of each object field.
	cols := make([]string, 0, template.NumField())
	for i := 0; i < template.NumField(); i++ {
		field := template.Field(i)
		if col, ok := field.Tag.Lookup("sql"); ok {
			cols = append(cols, col)
		}
	}

	// Format the columns into a comma-separated list.
	colJoin := strings.Join(cols, ", ")

	// Translate the columns, table, and conditional into an SQL statement.
	stmt := fmt.Sprintf("SELECT %s FROM %s;", colJoin, table)
	if cond != "" {
		stmt = fmt.Sprintf("SELECT %s FROM %s WHERE %s FOR UPDATE;", colJoin, table, cond)
		stmt = fmt.Sprintf(stmt, args...)
	}

	// Execute the query on the SQL database.
	rows, err := conn.query(stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %v", stmt, err)
	}

	// Parse the rows from the query into a slice of Go objects based on the prototype.
	return parseResponse(rows, object)
}

// Insert inserts the given Entries into the provided table using the Connection
// receiver.
//
// WARNING: This function is deprecated; use InsertObject() instead.
func (conn *Connection) Insert(table string, entries []dataforms.Entry) error {
	// Construct two slices to hold the values and columns of the Entries.
	numEntries := len(entries)
	vals := make([]string, 0, numEntries)
	cols := make([]string, 0, numEntries)

	// Append each Entry to the slices using the default Go representation.
	for _, entry := range entries {
		col := fmt.Sprintf("%v", entry.ColumnName)
		val := fmt.Sprintf("%v", entry.Value)
		vals = append(vals, val)
		cols = append(cols, col)
	}

	// Format the columns and values into comma-separated lists.
	colList := strings.Join(cols, ", ")
	valList := strings.Join(cols, ", ")

	// Execute the insertion on the SQL database.
	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON CONFLICT DO NOTHING/UPDATE;", table, colList, valList)
	_, err := conn.exec(stmt)
	return err
}

// InsertObject inserts the given object into the specified table and returns
// the record ID of the inserted row.
func (conn *Connection) InsertObject(table string, object interface{}) (int, error) {
	// Extract the underlying type and value of the object.
	objType := reflect.TypeOf(object)
	objValue := reflect.ValueOf(object)

	// Ensure the given object is a structure.
	if objType.Kind() != reflect.Struct {
		return 0, fmt.Errorf("type %T is not a structure", object)
	}

	// Cache the number of fields in the object; this value is used a few times.
	numFields := objType.NumField()

	// Construct a slice that holds the SQL column names of object fields.
	cols := make([]string, 0, numFields)
	// Construct a slice that holds the PostreSQL backreferences of object fields.
	refs := make([]string, 0, numFields)
	// Construct a slice that holds the values of object fields.
	vals := make([]interface{}, 0, numFields)

	// Append an element to each slice for every SQL field in the object.
	for i := 0; i < numFields; i++ {
		// Extract the type and value of the current field.
		fieldType := objType.Field(i)
		fieldValue := objValue.Field(i)

		// Derive the SQL column name corresponding to the current field.
		col, ok := fieldType.Tag.Lookup("sql")
		if !ok {
			logger.Warning("Field %q in structure %T does not have an SQL column tag.", fieldType.Name, object)
			continue
		}

		// Skip the current field if it has a SERIAL type.
		typ := fieldType.Tag.Get("typ")
		if strings.Contains(strings.ToUpper(typ), "SERIAL") {
			continue
		}

		// Let the PostgreSQL driver handle the formatting of the value.
		val := fieldValue.Interface()

		// The PostgreSQL backreference format is the same as the regex format.
		ref := fmt.Sprintf("$%d", len(refs)+1)

		// Update the column, backreference, and value slices.
		cols = append(cols, col)
		refs = append(refs, ref)
		vals = append(vals, val)
	}

	// Format the columns and backreferences into comma-separated lists.
	colList := strings.Join(cols, ", ")
	refList := strings.Join(refs, ", ")

	// Declare an integer to hold the record ID returned by the INSERT statement.
	var id int

	// Insert the object into the specified table and return the record ID.
	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON CONFLICT DO NOTHING RETURNING id;", table, colList, refList)
	row := conn.queryRow(stmt, vals...)
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return id, err
}

// UpdateObject updates the given object in the specified table.
func (conn *Connection) UpdateObject(table string, object interface{}) error {
	// Extract the underlying type and value of the object.
	objTyp := reflect.TypeOf(object)
	objVal := reflect.ValueOf(object)

	// Cache the number of fields in the object; this value is used a few times.
	numFields := objTyp.NumField()

	// Construct a slice that holds the SET clause entries of the UPDATE command.
	sets := make([]string, 0, numFields)
	// Construct a slice that holds the values of object fields.
	vals := make([]interface{}, 0, numFields)

	// Declare an integer to hold the backreference index of the ID column.
	var id int

	// Append an element to each slice for every SQL field in the object.
	for i := 0; i < numFields; i++ {
		// Extract the type and value of the current field.
		fieldTyp := objTyp.Field(i)
		fieldVal := objVal.Field(i)

		// Derive the SQL column name corresponding to the current field.
		col, ok := fieldTyp.Tag.Lookup("sql")
		if !ok {
			logger.Warning("Field %q in structure %T does not have an SQL column tag.", fieldTyp.Name, object)
			continue
		}

		// Let the PostgreSQL driver handle the formatting of the value.
		val := fieldVal.Interface()

		// Create a PostgreSQL SET clause entry with a backreference to the field value.
		ref := len(vals) + 1
		set := fmt.Sprintf("%s = $%d", col, ref)

		// Update the SET clause and value slices.
		sets = append(sets, set)
		vals = append(vals, val)

		// Set the backreference index of the ID column.
		if col == "id" {
			id = ref
		}
	}

	// Format the SET clause as a comma-separated list of SET clause entries.
	setList := strings.Join(sets, ", ")

	// Update the object in the specified table.  For more information, see
	// https://www.postgresql.org/docs/current/sql-update.html.
	stmt := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d;", table, setList, id)
	_, err := conn.exec(stmt, vals...)
	return err
}

// DeleteObject deletes the given object from the specified table.
func (conn *Connection) DeleteObject(table string, object interface{}) error {
	// Extract the underlying type and value of the object.
	objTyp := reflect.TypeOf(object)
	objVal := reflect.ValueOf(object)

	// isID reports whether the field with the provided name is an ID column.
	isID := func(name string) bool {
		field, _ := objTyp.FieldByName(name)
		return field.Tag.Get("sql") == "id"
	}

	// Retrieve the object ID.
	id := objVal.FieldByNameFunc(isID)
	if reflect.DeepEqual(id, reflect.Value{}) {
		return fmt.Errorf("structure %T does not have an ID field", object)
	}

	// Delete the object from the specified table.  For more information, see
	// https://www.postgresql.org/docs/current/sql-delete.html.
	stmt := fmt.Sprintf("DELETE FROM %s WHERE id = $1;", table)
	_, err := conn.exec(stmt, id.Interface())
	return err
}

//Lock will execute the SQL BEGIN command which aids in concurrent operations
//Unlock must be called once the transaction is complete.
func (conn *Connection) Lock() error {
	_, err := conn.db.Exec("BEGIN;")
	if err != nil {
		return fmt.Errorf("Failed to execute BEGIN statement: %v", err)
	}
	return nil
}

//Unlock will execute the SQL END command
func (conn *Connection) Unlock() error {
	_, err := conn.db.Exec("END;")
	if err != nil {
		return fmt.Errorf("Failed to execute BEGIN statement: %v", err)
	}
	return nil
}
