package dao

import (
	// log "github.com/cihub/seelog"
	"github.com/gocql/gocql"

	"github.com/juliendsv/go-cassadmin/domain"
)

const (
	clusterNodes = "127.0.0.1"
	port         = 19043
	consistency  = gocql.Quorum
)

type CassandraStore struct {
	cluster gocql.ClusterConfig
	Session gocql.Session
}

func NewCassandraStore() (domain.NOSQLStore, error) {
	c := gocql.NewCluster(clusterNodes)
	c.Consistency = consistency
	c.Port = port
	session, err := c.CreateSession()
	if err != nil {
		return nil, err
	}
	return CassandraStore{
		cluster: *c,
		Session: *session,
	}, nil
}

func (c CassandraStore) ListKeyspaces() ([]domain.Keyspace, error) {
	rows, err := c.Session.Query("SELECT keyspace_name, columnfamily_name FROM system.schema_columnfamilies;").Iter().SliceMap()
	if err != nil {
		return nil, err
	}

	map_ks := make(map[string][]string, 0)
	// log.Infof("result select: %v", rows)
	for _, r := range rows {
		if r["keyspace_name"] != "system" && r["keyspace_name"] != "system_traces" {
			map_ks[r["keyspace_name"].(string)] = append(map_ks[r["keyspace_name"].(string)], r["columnfamily_name"].(string))
		}
	}

	keyspaces := make([]domain.Keyspace, 0)
	for ks, cf := range map_ks {
		keyspaces = append(keyspaces, domain.Keyspace{
			Name:           ks,
			Columnfamilies: makeColumnfamilies(cf),
		})
	}

	return keyspaces, nil
}

// func (c CassandraStore) ListColumnfamilies(string keyspace) ([]domain.Columnfamily, error) {
// 	rows, err := c.Session.Query("SELECT keyspace_name, columnfamily_name FROM system.schema_columnfamilies;").Iter().SliceMap()
// 	if err != nil {
// 		return nil, err
// 	}

// 	map_ks := make(map[string][]string, 0)
// 	log.Infof("result select: %v", rows)
// 	for _, r := range rows {
// 		if r["keyspace_name"] != "system" && r["keyspace_name"] != "system_traces" {
// 			map_ks[r["keyspace_name"].(string)] = append(map_ks[r["keyspace_name"].(string)], r["columnfamily_name"].(string))
// 		}
// 	}

// 	keyspaces := make([]domain.Keyspace, 0)
// 	for ks, cf := range map_ks {
// 		keyspaces = append(keyspaces, domain.Keyspace{
// 			Name:           ks,
// 			Columnfamilies: makeColumnfamilies(cf),
// 		})
// 	}

// 	return keyspaces, nil
// }

func makeColumnfamilies(names []string) []domain.Columnfamily {
	columnfamilies := make([]domain.Columnfamily, len(names))
	for i, name := range names {
		columnfamilies[i] = domain.Columnfamily{Name: name}
	}
	return columnfamilies
}