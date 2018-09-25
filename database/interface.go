package database

// Query used to list records
type Query struct {
	Offset              int
	Limit               int
	Parent              string
	Prefix              string
	OrderBy             string
	OrderByDesc         bool
	OrderByProperty     string
	OrderByPropertyDesc bool
	PropertyFilter      map[string]string
}

// Key used to identify and look up a record
type Key struct {
	ID   string
	Path string
}

// Index on a specified Record property used to ensure efficient queries.
// Optionally a Parent path may be specified to create the Index only on
// a subset of Records with that Parent path.
type Index struct {
	Name     string
	Parent   string
	Property string
}

// Flame is the database interface
type Flame interface {

	// Get a Record using the key
	Get(Key) (*Record, error)

	// Save the Record
	Save(*Record) error

	// Delete the Record
	Delete(*Record) error

	// List Records matching the query
	List(Query) ([]*Record, error)

	// GetIndexes returns a list of existing indexes
	GetIndexes() ([]Index, error)

	// CreateIndex creates a new index on the specified field or property
	CreateIndex(Index) error

	// DeleteIndex deletes the index
	DeleteIndex(Index) error

	// HasIndex returns true if the Index already exists
	HasIndex(Index) (bool, error)

	// DropTable removes the table from the database, deleting all records
	DropTable() error

	// Migrate the database schema
	Migrate() error
}
