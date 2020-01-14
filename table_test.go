package storage

import (
	"reflect"
	"testing"
	"time"
)

//TestGetColumnType tests the getColumnType() method.
func TestGetColumnType(t *testing.T) {
	tests := []struct {
		typ      reflect.Type
		tag      reflect.StructTag
		wantType string
		wantErr  bool
	}{
		{
			reflect.TypeOf(new(int)),
			``,
			"",
			true,
		}, {
			reflect.TypeOf([]int{}),
			``,
			"",
			true,
		}, {
			reflect.TypeOf(map[int]int{}),
			``,
			"",
			true,
		}, {
			reflect.TypeOf(true),
			``,
			"BOOL",
			false,
		}, {
			reflect.TypeOf(int(0)),
			``,
			"INT4",
			false,
		}, {
			reflect.TypeOf(int16(0)),
			``,
			"INT2",
			false,
		}, {
			reflect.TypeOf(int32(0)),
			``,
			"INT4",
			false,
		}, {
			reflect.TypeOf(int64(0)),
			``,
			"INT8",
			false,
		}, {
			reflect.TypeOf(float32(0)),
			``,
			"FLOAT4",
			false,
		}, {
			reflect.TypeOf(float64(0)),
			``,
			"FLOAT8",
			false,
		}, {
			reflect.TypeOf(""),
			``,
			"TEXT",
			false,
		}, {
			reflect.TypeOf([]byte{}),
			``,
			"BYTEA",
			false,
		}, {
			reflect.TypeOf(time.Time{}),
			``,
			"TIMESTAMP",
			false,
		}, {
			reflect.TypeOf(time.Time{}),
			`typ:"FAKENEWS"`,
			"FAKENEWS",
			false,
		},
	}
	for i, test := range tests {
		field := reflect.StructField{Type: test.typ, Tag: test.tag}
		haveType, haveErr := getColumnType(field)
		if (haveErr != nil) != test.wantErr {
			t.Errorf("TestGetColumnType()[%d] = %v, want error %t.", i, haveErr, test.wantErr)
		}
		if haveType != test.wantType {
			t.Errorf("TestGetColumnType()[%d] = %q, want type %q.", i, haveType, test.wantType)
		}
	}
}

// TestCreateDropTable tests the (*Connection).CreateTableFromObject() and
// (*Connection).DropTable() methods.
func TestCreateDropTable(t *testing.T) {
	creds := GetTestCreds()

	tests := []struct {
		name    string
		object  interface{}
		insert  string
		wantErr bool
	}{
		{
			"empty",
			false,
			"",
			true,
		}, {
			"license",
			struct {
				DOB time.Time `sql:"dob"`
			}{},
			"",
			true,
		}, {
			"identifier",
			struct {
				ID int16 `sql:"id"`
			}{},
			"INSERT INTO identifier (id) VALUES (0)",
			false,
		}, {
			"material",
			struct {
				ID     int32   `sql:"id"`
				Name   string  `sql:"name"`
				Mass16 int16   `sql:"mass16"`
				Mass32 int32   `sql:"mass32"`
				Mass64 int64   `sql:"mass64"`
				Heat32 float32 `sql:"heat32"`
				Heat64 float64 `sql:"heat64"`
			}{},
			"INSERT INTO material (id, name, mass16, mass32, mass64, heat32, heat64) VALUES (0, '', 0, 0, 0, 0, 0)",
			false,
		}, {
			"Tree",
			struct {
				ID  int64  `sql:"id" typ:"BIGSERIAL" opt:"PRIMARY KEY"`
				Oak bool   `sql:"oak"`
				DNA []byte `sql:"dna"`
			}{
				ID:  1,
				Oak: true,
				DNA: []byte{1, 2, 3},
			},
			`INSERT INTO tree (oak, dna) VALUES (true, '\\001\\002\\003')`,
			false,
		},
	}

	// Connect to the test database.
	conn, err := Connect(creds)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v.", err)
	}
	defer conn.Close()

	for i, test := range tests {
		// Create the table.
		err := conn.CreateTableFromObject(test.name, test.object)
		if (err != nil) != test.wantErr {
			t.Errorf("TestCreateDropTable()[%d] = %v, want table error %t.", i, err, test.wantErr)
		}
		if err != nil {
			continue
		}

		// Insert the object into the table.
		if _, err := conn.exec(test.insert); err != nil {
			t.Errorf("TestCreateDropTable()[%d] - failed to insert object: %v.", i, err)
			continue
		}

		// Retrieve the object from the table.
		if rows, err := conn.query("SELECT * FROM " + test.name); err != nil {
			t.Errorf("TestCreateDropTable()[%d] - failed to execute query: %v.", i, err)
		} else if !rows.Next() {
			t.Errorf("TestCreateDropTable()[%d] - no rows were inserted.", i)
		}

		// Drop the table.
		if err := conn.DropTable(test.name); err != nil {
			t.Errorf("TestCreateDropTable()[%d] - failed to drop table: %v.", i, err)
		}
	}
}
