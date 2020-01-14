# StructQL
This package abstracts the connection and management of a Postgres database server using native go structs.

## StructQL Logs
StructQL uses github.com/inflowml/logger for all logging. These will appear as standard out. In order to turn off StructQL logs export SQL_LOG=0.

## Storage Functions
Part of the `github.com/inflowml/inflow-micro/common/database/storage` package.
### Connect
Connects to the appropriate db based on the micro-service and returns the connection structure.
```go
func Connect(creds ConnectionConfig) (*Connection, error)
```
### Close
Closes the connection to the database, must be called when the microservice is finished using the db. Connections should only be closed when the micro-service terminates or is killed if possible.
```go
func (conn *Connection) Close() error
```
### CreateTableFromObject
Create table accepts the name of the table to be created and an interface representing the table columns.
```go
func (conn *Connection) CreateTableFromObject(table string, object interface{}) error {
```
A table is created out of a native go struct with StructQL tags. For example in order to create a `person` table with columns id, name, age, DNA the following code would be appropriate.
```go
type Person struct {
	ID   int32  `sql:"id" typ:"SERIAL" opt:"PRIMARY KEY"`
	Name string `sql:"name"`
	Age  int32  `sql:"age"`
	DNA  []byte `sql:"dna"`
}
...
err := conn.CreateTableFromObject("person", Person{})
if err != nil {
	// Handle Error
}
...
```
### InsertObject
InsertObject accepts a table name and an object interface and inserts it into the database
```go
func (conn *Connection) InsertObject(table string, object interface{}) error
```
The object interface must be tagged with the SQL Column names. Ref the following example of an acceptable struct which corresponds to a Person database table with columns id, name, age, and dna.
```go
type Person struct {
	ID   int32  `sql:"id" typ:"SERIAL" opt:"PRIMARY KEY"`
	Name string `sql:"name"`
	Age  int32  `sql:"age"`
	DNA  []byte `sql:"dna"`
}
```
In this case every time a new row is inserted a unique id will be assigned in the id column of the table. This will be automatically done by Postgres.
### SelectFrom
Accepts a struct type, and table name and returns the query as a slice of given struct. Note that the fields in the given struct are the columns that are listed in the `SELECT <Columns>` portion of the SQL query.
```go
func (conn *Connection) SelectFrom(object interface{}, table string) (interface{}, error) 
```
The following is an example of a struct type that would result in the query for Name and Age from a table
```go
type Person struct {
	Name string `sql:"name"`
	Age  int32  `sql:"age"`
}
```
### SelectFromWhere
Accepts a struct type, table name, and conditional and returns the query as a slice of given struct. Note that the fields in the given struct are the columns that are listed in the SELECT <Columns> portion of the SQL query. Additonally the conditional must be a string using standard SQL comparisons such as `age >= 50 AND name == John`
```go
func (conn *Connection) SelectFromWhere(object interface{}, table string, conditional string) (interface{}, error) 
```
The following is an example of a struct type that would result in the query for Name and Age from a table
```go
type Person struct {
	Name string `sql:"name"`
	Age  int32  `sql:"age"`
}
```

## Data Formats
The following are structs used to interact with InFlow storage that can be accessed via `dataforms.<Struct>` after importing `github.com/inflowml/inflow-micro/common/database/dataforms`
### SQLTypes
SQLTypes are the accepted storage types for InFlow databases. They can be accessed as dataforms.<SQLType>. The following are valid SQLTypes
- `IntSQL`
- `StringSQL`
- `TimeSQL`
- `JSONSQL`
- `BoolSQL`

### ColumnHeader
The ColumnHeader struct is used to define columns when creating a new table. It is comprised of a name of type string and SQLType of type SQLType
```go
type ColumnHeader struct {
	Name    string
	SQLType SQLType
}
```

## Testing Configurations
In order to run StructQL tests a local postgres server is required. One can be installed through by running `sudo install.sh` in the testutils directory. Once installed run test-srv.sh. Once you are finished with the server run `sudo service postgresql stop`

