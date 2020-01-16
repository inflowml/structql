package structql

// Driver is a custom type for all supported SQL drivers
type Driver string

const (
	// Postgres driver value is to be used for postgres databases
	Postgres Driver = "postgres"
	// MySQL driver value is to be used for MySql databases
	MySQL Driver = "mysql"
)
