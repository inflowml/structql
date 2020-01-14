package dataforms

// ColumnHeader represents a column header in a SQL table.
type ColumnHeader struct {
	Name    string
	SQLType SQLType
}

// Entry represents an entry in a SQL table.
type Entry struct {
	ColumnName string
	Value      string
}
