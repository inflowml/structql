// Package structql implements the Database structure.
// This file contains tests for manager.go.
package structql

import (
	"os"
	"testing"
)

// TestConnect tests the Connect() and Close() methods.
func TestConnectClose(t *testing.T) {
	creds := GetTestCreds()

	conn, err := Connect(creds)
	if err != nil {
		t.Fatalf("TestConnectClose() - failed to connect: %v.", err)
	}
	if err = conn.db.Close(); err != nil {
		t.Errorf("TestConnectClose() - failed to close: %v.", err)
	}
}

// setup asserts the creation of a Database connection.
func setup(t *testing.T) *Connection {
	creds := GetTestCreds()

	conn, err := Connect(creds)
	if err != nil {
		t.Fatalf("Failed to connect to SQL database %q during setup: %v.", conn.name, err)
	}
	return conn
}

func GetTestCreds() ConnectionConfig {
	driverEnv := os.Getenv("SQL_DRIVER")
	driver := Postgres
	if driverEnv == "MY_SQL" {
		driver = MySQL
	}
	return ConnectionConfig{
		User:     "StructqlUser",
		Password: "StructqlPW",
		Database: "testdb",
		Host:     "localhost",
		Port:     "5432",
		Driver:   driver,
	}
}
