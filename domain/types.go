package domain

var (
	DefaultStore NOSQLStore
)

type NOSQLStore interface {
	ListKeyspaces() ([]Keyspace, error)
	// ListColumnFamilyResults(string) error
}

type Columnfamily struct {
	Name string
}

type CfRow struct {
	Name string
}

type Keyspace struct {
	Name           string
	Columnfamilies []Columnfamily
}
