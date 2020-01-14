// Package storage implements the Database structure.
// This file contains tests for manager.go.
package storage

import (
	"testing"
)

// TestConnect tests the Connect() and Close() methods.
func TestConnectClose(t *testing.T) {
	conn, err := Connect()
	if err != nil {
		t.Fatalf("TestConnectClose() - failed to connect: %v.", err)
	}
	if err = conn.db.Close(); err != nil {
		t.Errorf("TestConnectClose() - failed to close: %v.", err)
	}
}

// setup asserts the creation of a Database connection.
func setup(t *testing.T) *Connection {
	conn, err := Connect()
	if err != nil {
		t.Fatalf("Failed to connect to SQL database %q during setup: %v.", conn.name, err)
	}
	return conn
}
