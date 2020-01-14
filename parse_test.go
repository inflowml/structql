// Package storage implements the Database structure.
// This file contains tests for parse.go.
package storage

import (
	"fmt"
	"reflect"
	"testing"
)

// TestParseResponse tests the parseResponse() method.
func TestParseResponse(t *testing.T) {
	type Person struct {
		Name string  `sql:"name"`
		Age  int32   `sql:"age"`
		Mass float32 `sql:"mass"`
	}

	adam := Person{"Adam", 10, 242.0}
	brad := Person{"Brad", 20, 199.9}
	chad := Person{"Chad", 30, 206.9}

	tests := []struct {
		query      string
		wantPeople []Person
	}{
		{
			`SELECT * FROM People WHERE name = 'Duke'`,
			[]Person{},
		}, {
			`SELECT * FROM People WHERE name = 'Adam'`,
			[]Person{adam},
		}, {
			`SELECT * FROM People WHERE age >= 20`,
			[]Person{brad, chad},
		},
	}

	// Create a suitable table in the test database.
	conn, err := Connect()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v.", err)
	}
	if _, err := conn.exec(`CREATE TABLE People (name TEXT, age INT, mass FLOAT4);`); err != nil {
		t.Fatalf("Failed to create table: %v.", err)
	}
	defer func() {
		conn.exec(`DROP TABLE People;`)
		conn.Close()
	}()

	// Add Adam, Brad, and Chad to the database.
	for _, person := range []Person{adam, brad, chad} {
		cmd := fmt.Sprintf("INSERT INTO People (name, age, mass) VALUES ('%s', %d, %f);", person.Name, person.Age, person.Mass)
		if _, err := conn.exec(cmd); err != nil {
			t.Fatalf("Failed to insert Person %q: %v.", person.Name, err)
		}
	}

	for i, test := range tests {
		rows, err := conn.query(test.query)
		if err != nil {
			t.Errorf("TestParseResponse()[%d] - failed to execute query: %v.", i, err)
			continue
		}

		havePeople, err := parseResponse(rows, Person{})
		if err != nil {
			t.Errorf("TestParseResponse()[%d] - failed to parse response: %v.", i, err)
			continue
		}

		if len(havePeople) != len(test.wantPeople) {
			t.Errorf("TestParseResponse()[%d] = %d, want %d people.", i, len(havePeople), len(test.wantPeople))
			continue
		}
		for j, havePerson := range havePeople {
			wantPerson := test.wantPeople[j]
			if !reflect.DeepEqual(havePerson, wantPerson) {
				t.Errorf("TestParseResponse()[%d][%d] = %v, want Person %v.", i, j, havePerson, wantPerson)
			}
		}
	}
}
