package dao

import (
	"fmt"
	"strconv"

	log "github.com/cihub/seelog"
	"github.com/gocql/gocql"

	"github.com/juliendsv/go-cassadmin/domain"
)

const (
	clusterNodes = "127.0.0.1"
	port         = 19043
	consistency  = gocql.Quorum

	defaultLimit = 100
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

func (c CassandraStore) CreateKeyspace(keyspace_name string) error {
	err := c.exec(fmt.Sprintf(`CREATE KEYSPACE %s
	WITH replication = {
		'class' : 'SimpleStrategy',
		'replication_factor' : %d
	}`, keyspace_name, 1))
	return err
}

func (c CassandraStore) DropKeyspace(keyspace_name string) error {
	err := c.exec(fmt.Sprintf(`DROP KEYSPACE %s`, keyspace_name))
	return err
}

func (c CassandraStore) exec(query string) error {
	if err := c.Session.Query(query).Consistency(consistency).Exec(); err != nil {
		return fmt.Errorf("error executing query %s: %v", query, err)
	}
	return nil
}

func (c CassandraStore) ShowKeyspaces() ([]domain.Keyspace, error) {
	ks, err_ks := c.Session.Query("SELECT keyspace_name FROM system.schema_keyspaces;").Iter().SliceMap()
	if err_ks != nil {
		return nil, err_ks
	}
	map_ks := make(map[string][]string, 0)
	for _, r := range ks {
		if r["keyspace_name"] != "system" && r["keyspace_name"] != "system_traces" {
			map_ks[r["keyspace_name"].(string)] = []string{}
		}
	}

	// TODO split this in two function kss ans cfs

	cfs, err_cf := c.Session.Query("SELECT keyspace_name, columnfamily_name FROM system.schema_columnfamilies;").Iter().SliceMap()
	if err_cf != nil {
		return nil, err_cf
	}
	for _, r := range cfs {
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

func (c CassandraStore) ShowColumnFamily(ks, cf string) ([]map[string]string, error) {
	rows, err := c.Session.Query("SELECT * FROM  " + ks + "." + cf + ";").Iter().SliceMap()
	if err != nil {
		return nil, err
	}
	map_cf := make(map[string]string)
	res := make([]map[string]string, len(rows))
	for i, r := range rows {
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
		res[i] = map_cf
	}

	log.Infof("res: %v", res)
	return res, nil
}

func makeColumnfamilies(names []string) []domain.Columnfamily {
	columnfamilies := make([]domain.Columnfamily, len(names))
	for i, name := range names {
		columnfamilies[i] = domain.Columnfamily{Name: name}
	}
	return columnfamilies
}
