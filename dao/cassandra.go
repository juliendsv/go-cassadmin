package dao

import (
	"strconv"

	log "github.com/cihub/seelog"
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

func (c CassandraStore) ShowKeyspaces() ([]domain.Keyspace, error) {
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

func (c CassandraStore) ShowColumnFamily(ks, cf string) error {
	rows, err := c.Session.Query("SELECT * FROM  " + ks + "." + cf + ";").Iter().SliceMap()
	if err != nil {
		return err
	}

	map_cf := make(map[string]string)

	for _, r := range rows {
		for k, result := range r {
			// TODO manage all types, save it
			// and maybe we should save the value as bytes[] instead of string
			switch v := result.(type) {
			case int:
				map_cf[k] = strconv.Itoa(v)
			case string:
				map_cf[k] = v
			}
		}
		log.Infof("map: %v", map_cf)
	}
	return nil
}

func makeColumnfamilies(names []string) []domain.Columnfamily {
	columnfamilies := make([]domain.Columnfamily, len(names))
	for i, name := range names {
		columnfamilies[i] = domain.Columnfamily{Name: name}
	}
	return columnfamilies
}
