// Package storage implements the Database structure.
// This file contains tests for parse.go.
package storage

import (
	"fmt"
	"reflect"
	"testing"
)

// TestSelectFromWhere tests the (*Connection).SelectFromWhere() method.
func TestSelectFromWhere(t *testing.T) {
	type Person struct {
		ID   int32  `sql:"id" typ:"SERIAL"`
		Name string `sql:"name"`
	}

	personA := Person{1, "A"}
	personB := Person{2, "B"}
	personC := Person{3, "C"}

	tests := []struct {
		cond       string
		args       []interface{}
		wantPeople []Person
	}{
		{
			"",
			[]interface{}{},
			[]Person{personA, personB, personC},
		}, {
			"name = 'C'",
			[]interface{}{},
			[]Person{personC},
		}, {
			"id <= 2",
			[]interface{}{},
			[]Person{personA, personB},
		}, {
			"id <= 2 AND name = 'A'",
			[]interface{}{},
			[]Person{personA},
		}, {
			"id = %d",
			[]interface{}{3},
			[]Person{personC},
		},
	}

	// Create a suitable table in the test database.
	conn := createTableUnsafe("People", Person{})
	defer conn.Close()
	defer conn.DropTable("People")

	for _, person := range []Person{personA, personB, personC} {
		if _, err := conn.InsertObject("People", person); err != nil {
			t.Fatalf("Failed to insert Person %v: %v.", person, err)
		}
	}

	for i, test := range tests {
		// Execute the SELECT FROM WHERE query.
		people, err := conn.SelectFromWhere(Person{}, "People", test.cond, test.args...)
		if err != nil {
			t.Errorf("TestSelectFromWhere()[%d] - failed to execute query: %v.", i, err)
			continue
		}

		// Cast the []interface{} slice into a []Person{} slice.
		havePeople := make([]Person, 0, len(people))
		for _, personI := range people {
			person := personI.(Person)
			havePeople = append(havePeople, person)
		}

		// Compare the retrieved and expected Person slices.
		if !reflect.DeepEqual(havePeople, test.wantPeople) {
			t.Errorf("TestSelectFromWhere()[%d] = %v, want people %v.", i, havePeople, test.wantPeople)
		}
	}
}

// TestInsertObject tests the (*Connection).InsertObject() method.
func TestInsertObject(t *testing.T) {
	type Person struct {
		ID   int16  `sql:"id" typ:"SMALLSERIAL"`
		Name string `sql:"name"`
		Age  int32  `sql:"age"`
		DNA  []byte `sql:"dna"`
	}

	tests := []struct {
		person Person
	}{
		{
			Person{
				ID:   1,
				Name: "",
				Age:  0,
				DNA:  []byte{},
			},
		}, {
			Person{
				ID:   2,
				Name: "John Cena",
				Age:  42,
				DNA:  []byte{1, 2, 3},
			},
		},
	}

	conn := createTableUnsafe("People", Person{})
	defer conn.Close()
	defer conn.DropTable("People")

	for i, test := range tests {
		// Insert the Person into the database.
		haveID, err := conn.InsertObject("People", test.person)
		if err != nil {
			t.Errorf("TestInsertObject()[%d] - failed to insert object: %v.", i, err)
			continue
		}

		// Verify that returned record ID scales with the test index.
		if wantID := i + 1; haveID != wantID {
			t.Errorf("TestInsertObject()[%d] = %d, want record ID %v.", i, haveID, wantID)
		}

		// Retrieve the Person from the database.
		query := fmt.Sprintf(`SELECT * FROM People WHERE name = '%s'`, test.person.Name)
		rows, err := conn.query(query)
		if err != nil {
			t.Errorf("TestInsertObject()[%d] - failed to execute query: %v.", i, err)
			continue
		}

		people, err := parseResponse(rows, Person{})
		if err != nil {
			t.Errorf("TestInsertObject()[%d] - failed to parse response: %v.", i, err)
			continue
		}

		if len(people) != 1 {
			t.Errorf("TestInsertObject()[%d] = %d, want 1 Person.", i, len(people))
			continue
		}

		// Verify that the retrieved Person is the same Person that was inserted.
		person := people[0]
		if !reflect.DeepEqual(test.person, person) {
			t.Errorf("TestInsertObject()[%d] = %v, want Person %v.", i, person, test.person)
		}
	}
}

// TestUpdateObject tests the (*Connection).UpdateObject() method.
func TestUpdateObject(t *testing.T) {
	type Person struct {
		ID   int32  `sql:"id" opt:"PRIMARY KEY"`
		Name string `sql:"name"`
	}

	tests := []struct {
		person Person
	}{
		{
			Person{
				ID:   1,
				Name: "Joseph",
			},
		}, {
			Person{
				ID:   1,
				Name: "Faith",
			},
		},
	}

	conn := createTableUnsafe("People", Person{})
	defer conn.Close()
	defer conn.DropTable("People")

	base := Person{ID: 1, Name: "Rook"}
	if _, err := conn.InsertObject("People", base); err != nil {
		t.Fatalf("Failed to insert %#v into table: %v.", base, err)
	}

	for i, test := range tests {
		if err := conn.UpdateObject("People", test.person); err != nil {
			t.Errorf("TestUpdateObject()[%d] - failed to update object: %v.", i, err)
			continue
		}

		people, err := conn.SelectFrom(Person{}, "People")
		if err != nil {
			t.Errorf("TestUpdateObject()[%d] - failed to query database: %v.", i, err)
		} else if len(people) != 1 {
			t.Errorf("TestUpdateObject()[%d] = %d, want 1 Person.", i, len(people))
		} else if !reflect.DeepEqual(test.person, people[0]) {
			t.Errorf("TestUpdateObject()[%d] = %v, want Person %v.", i, people[0], test.person)
		}
	}
}

// TestDeleteObject tests the (*Connection).DeleteObject() method.
func TestDeleteObject(t *testing.T) {
	type Pizza struct {
		ID      int32  `sql:"id" opt:"PRIMARY KEY"`
		Topping string `sql:"topping"`
	}

	cheese := Pizza{1, "Cheese"}
	deluxe := Pizza{2, "Deluxe"}

	tests := []struct {
		pizza      Pizza
		wantPizzas []Pizza
	}{
		{
			Pizza{3, "Pepperoni"},
			[]Pizza{cheese, deluxe},
		}, {
			cheese,
			[]Pizza{deluxe},
		}, {
			Pizza{3, "Pepperoni"},
			[]Pizza{deluxe},
		}, {
			deluxe,
			[]Pizza{},
		},
	}

	conn := createTableUnsafe("Pizza", Pizza{})
	defer conn.Close()
	defer conn.DropTable("Pizza")

	for _, pizza := range []Pizza{cheese, deluxe} {
		if _, err := conn.InsertObject("Pizza", pizza); err != nil {
			t.Fatalf("Failed to insert Pizza %v: %v.", pizza, err)
		}
	}

	for i, test := range tests {
		if err := conn.DeleteObject("Pizza", test.pizza); err != nil {
			t.Errorf("TestDeleteObject()[%d] - failed to delete Pizza: %v.", i, err)
			continue
		}

		rows, err := conn.SelectFrom(Pizza{}, "Pizza")
		if err != nil {
			t.Errorf("TestDeleteObject()[%d] - failed to select Pizza: %v.", i, err)
			continue
		}

		havePizzas := make([]Pizza, len(rows))
		for i, row := range rows {
			havePizzas[i] = row.(Pizza)
		}
		if !reflect.DeepEqual(havePizzas, test.wantPizzas) {
			t.Errorf("TestDeleteObject()[%d] = %v, want Pizza %v.", i, havePizzas, test.wantPizzas)
		}
	}
}

// createTableUnsafe constructs a database Connection and creates a table with
// the given name from the provided object.  Failure to do so results in a panic.
func createTableUnsafe(table string, object interface{}) *Connection {
	conn, err := Connect()
	if err != nil {
		panic(fmt.Sprintf("Failed to construct Connection: %v.", err))
	}
	if err := conn.CreateTableFromObject(table, object); err != nil {
		panic(fmt.Sprintf("Failed to create table %q: %v.", table, err))
	}
	return conn
}
