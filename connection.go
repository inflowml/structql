package storage

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/inflowml/inflow-micro/common/logger"
	_ "github.com/lib/pq" // The PostgreSQL driver.
)

// Connection wraps the sql.DB type.
type Connection struct {
	db   *sql.DB
	name string
}

// Connect establishes a connection to the SQL database specified by the
// INFLOW_SERVICE, CLOUD_SQL_PORT, and CLOUD_SQL_PW environment variables.
func Connect() (*Connection, error) {
	// Set this in app.yaml when running in production.
	user := os.Getenv("INFLOW_SERVICE")
	database := os.Getenv("INFLOW_SERVICE")
	port := os.Getenv("CLOUD_SQL_PORT")
	pw := os.Getenv("CLOUD_SQL_PW")
	connectionInfo := fmt.Sprintf("database=%s user=%s password=%s port=%s host=localhost sslmode=disable", database, user, pw, port)

	// Attempt to open the database (this does NOT initiate a connection).
	sqlDB, err := sql.Open("postgres", connectionInfo)
	if err != nil {
		logger.SQL("Failed to open SQL database %q.", database)
		return nil, fmt.Errorf("failed to open SQL database using %q: %v", connectionInfo, err)
	}

	// Wrap the sql.DB object in the Database wrapper.
	conn := Connection{sqlDB, database}

	//Initiates connection to db.
	if err := conn.db.Ping(); err != nil {
		logger.SQL("Failed to connect to SQL database proxy.  Ensure the proxy is running and the appropriate environment variables are set.")
		return nil, fmt.Errorf("failed to connect to SQL database proxy using %q: %v", connectionInfo, err)
	}

	logger.SQL("Successfully connected to SQL database %q.", database)
	return &conn, nil
}

// Close closes the connection to the Database receiver.
func (conn *Connection) Close() error {
	if err := conn.db.Close(); err != nil {
		return fmt.Errorf("failed to close SQL database: %v", err)
	}
	logger.SQL("Successfully closed connection to SQL database %q.", conn.name)
	return nil
}

// exec executes the given SQL statement with the provided arguments on the Database receiver.
func (conn *Connection) exec(stmt string, args ...interface{}) (sql.Result, error) {
	result, err := conn.db.Exec(stmt, args...)
	if err != nil {
		return nil, fmt.Errorf("failed To execute SQL statement %q (result %v): %v", stmt, result, err)
	}
	return result, nil
}

// query queries the Database receiver with the given SQL query.
func (conn *Connection) query(stmt string) (*sql.Rows, error) {
	rows, err := conn.db.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("failed To execute SQL query: %v", err)
	}
	return rows, nil
}

// queryRow executes the given SQL statement with the provided arguments on the
// Database receiver and returns the row resulting from the query.
func (conn *Connection) queryRow(stmt string, args ...interface{}) *sql.Row {
	return conn.db.QueryRow(stmt, args...)
}

/*func setDB(service string) error {
	stmt := fmt.Sprintf("\\c %s", service)
	if _, err := db.Exec(stmt); err == nil {
		return nil
	}
	logger.Warning("Service Does Not Have A DB Creating DB: %s", service)
	return createDB(service)
}

func createDB(name string) error {
	stmt := fmt.Sprintf("\\CREATE TABLE %s", name)
	if resp, err := db.Exec(stmt); err != nil {
		logger.Error("%v", resp)
		return fmt.Errorf("Failed To Create New DB: %v", err)
	}
	return nil
}*/
