package domain

var (
	DefaultStore NOSQLStore
)

type NOSQLStore interface {
	CreateKeyspace(keyspace_name string) error
	DropKeyspace(keyspace_name string) error
	ShowKeyspaces() ([]Keyspace, error)
	ShowColumnFamily(ks, cf string) ([]map[string]string, error)
}

type Columnfamily struct {
	Name string
}

type Row struct {
}

type Column struct {
	Name      string
	Value     string // should be something more generic
	Timestamp int
}

type Keyspace struct {
	Name           string
	Columnfamilies []Columnfamily
}
