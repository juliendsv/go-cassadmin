package domain

var (
	DefaultStore NOSQLStore
)

type NOSQLStore interface {
	ShowKeyspaces() ([]Keyspace, error)
	ShowColumnFamily(ks, cf string) ([]map[string]string, error)
}

type Columnfamily struct {
	Name string
}

type Keyspace struct {
	Name           string
	Columnfamilies []Columnfamily
}
