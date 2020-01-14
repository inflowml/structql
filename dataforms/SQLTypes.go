package dataforms

// SQLType is implemented by any type with a Format function which returns the
// equivalent PostgreSQL representation of the type.
type SQLType interface {
	Format() string
}

type (
	// IntSQLType represents a 32-bit integer SQL type.
	IntSQLType struct{}
	// Int64SQLType represents a 64-bit integer SQL type.
	Int64SQLType struct{}
	// FloatSQLType represents a 32-bit floating-point SQL type.
	FloatSQLType struct{}
	// DoubleSQLType represents a 64-bit floating-point SQL type.
	DoubleSQLType struct{}
	// StringSQLType represents a string SQL type.
	StringSQLType struct{}
	// TimeSQLType represents a timestamp SQL type.
	TimeSQLType struct{}
	// JSONSQLType represents a JSON SQL type.
	JSONSQLType struct{}
	// BoolSQLType represents a boolean SQL type.
	BoolSQLType struct{}
	// SerialSQLType represents an auto-incrementing integer SQLtype.
	SerialSQLType struct{}
)

// Created a singleton for each SQL type for the sake of convenience.
var (
	IntSQL    IntSQLType
	Int64SQL  Int64SQLType
	FloatSQL  FloatSQLType
	DoubleSQL DoubleSQLType
	StringSQL StringSQLType
	TimeSQL   TimeSQLType
	JSONSQL   JSONSQLType
	BoolSQL   BoolSQLType
	SerialSQL SerialSQLType
)

// Format returns the name of the 32-bit integer PostgreSQL type.
func (IntSQLType) Format() string {
	return "INTEGER"
}

// Format returns the name of the 64-bit integer PostgreSQL type.
func (Int64SQLType) Format() string {
	return "BIGINT"
}

// Format returns the name of the 32-bit floating-point PostgreSQL type.
func (FloatSQLType) Format() string {
	return "REAL"
}

// Format returns the name of the 64-bit floating-point PostgreSQL type.
func (DoubleSQLType) Format() string {
	return "DOUBLE PRECISION"
}

// Format returns the name of the string PostgreSQL type.
func (StringSQLType) Format() string {
	return "TEXT"
}

// Format returns the name of the timestamp PostgreSQL type.
func (TimeSQLType) Format() string {
	return "TIMESTAMP"
}

// Format returns the name of the JSON PostgreSQL type.
func (JSONSQLType) Format() string {
	return "JSON"
}

// Format returns the name of the boolean PostgreSQL type.
func (BoolSQLType) Format() string {
	return "BOOLEAN"
}

// Format returns the name of the serial PostgreSQL type.
func (SerialSQLType) Format() string {
	return "SERIAL"
}
