// Package structql implements the Database structure.
// This file contains tests for manager.go.
package structql

import (
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
	return ConnectionConfig{
		User:     "postgres",
		Password: "postgres",
		Database: "testdb",
		Host:     "localhost",
		Port:     "5432",
	}
}
