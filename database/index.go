package database

// PostgresIndex is a record in PG describing a table index
type PostgresIndex struct {
	TableName string `gorm:"tablename"`
	IndexName string `gorm:"indexname"`
	IndexDef  string `gorm:"indexdef"`
}

// PostgresIndexes is a slice of indexes
type PostgresIndexes []*PostgresIndex

// Get the index with the specified name
func (indexes PostgresIndexes) Get(name string) *PostgresIndex {
	for _, index := range indexes {
		if index.IndexName == name {
			return index
		}
	}
	return nil
}
