# Database
This library abstracts the connection and management of the cloud sql proxy required in order to store structured data on GCP.

## Logger
For logs related to your SQL export SQL_LOG=1 (This is the default using development configurations)

## Setup
### Installation
1. cd `./common/database/proxy`
2. Execute `./install.sh`
3. Execute `./start_proxy.sh`
4. Open a new terminal
5. From proj root: `source ./configure.sh`
6. Use storage functions to interact with database

## Storage Functions
Part of the `github.com/inflowml/inflow-micro/common/database/storage` package.
### Connect
Connects to the appropriate db based on the micro-service and returns the connection structure.
```go
func Connect() (*Connection, error)
```
### Close
Closes the connection to the database, must be called when the microservice is finished using the db. Connections should only be closed when the micro-service terminates or is killed if possible.
```go
func (conn *Connection) Close() error
```
### CreateTable
Create table accepts the name of the table to be created and a splice of ColumnHeaders representing the data to be stored in each column.
```go
func (conn *Connection) CreateTable(table string, headers []dataforms.ColumnHeader) error
```
### InsertObject
InsertObject accepts a table name and an object interface and inserts it into the database
```go
func (conn *Connection) InsertObject(table string, object interface{}) error
```
The object interface must be tagged with the SQL Column names. Ref the following example of an acceptable struct which corresponds to a Person database table with columns name, age, and dna.
```go
type Person struct {
	Name string `sql:"name"`
	Age  int32  `sql:"age"`
	DNA  []byte `sql:"dna"`
}
``` 
### SelectFrom
Accepts a struct type, and table name and returns the query as a slice of given struct. Note that the fields in the given struct are the columns that are listed in the `SELECT <Columns>` portion of the SQL query.
```go
func (conn *Connection) SelectFrom(prototype reflect.Type, table string) (interface{}, error) 
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
func (conn *Connection) SelectFrom(prototype reflect.Type, table string, conditional string) (interface{}, error) 
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

## Connect Via Terminal
After succesfully installing the proxy you can connect to the DB via PostgreSQL by executing `psql -U test -h localhost -p 3306`

## Notes
- In production all microservices have their own database, this means a microservice can only access it's own data.
- In development all microservices share a test database, this database will not be available in production.  Ensure that you can programatically set up the tables you need or contact a DB admin to set up production tables in advance.
- All communication between the proxy and the DB is SSL encrypted, local connection between an application and the proxy is not.

## Connect Via Terminal
After succesfully installing the proxy you can connect to the DB via PostgreSQL by executing `psql -U test -h localhost -p 3306`

## Notes
- In production all microservices have their own database, this means a microservice can only access it's own data.
- In development all microservices share a test database, this database will not be available in production.  Ensure that you can programatically set up the tables you need or contact a DB admin to set up production tables in advance.
- All communication between the proxy and the DB is SSL encrypted, local connection between an application and the proxy is not.

