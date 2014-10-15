package dao

import (
	"testing"
)

const (
	TEST_KEYSPACE = "this_is_a_test"
)

func createTestKeyspace() error {
	c, err_init := NewCassandraStore()
	if err_init != nil {
		return err_init
	}

	_ = c.DropKeyspace(TEST_KEYSPACE)
	err_create := c.CreateKeyspace(TEST_KEYSPACE)
	if err_create != nil {
		return err_create
	}
	return nil
}

func TestNewCassandraStore(t *testing.T) {
	err := createTestKeyspace()
	if err != nil {
		t.Fatalf("We were not expecting creating keyspace but got %v", err)
	}
	c, err_init := NewCassandraStore()
	if err_init != nil {
		t.Fatalf("We were not expecting error initializing cassandra but got %v", err_init)
	}
	ks, err := c.ShowKeyspaces()
	if err != nil {
		t.Fatalf("We were not expecting error listing keyspaces but got %v", err)
	}
	t.Fatalf("%v", ks)
}

// func TestShowCf(t *testing.T) {
// 	c, err_init := NewCassandraStore()
// 	if err_init != nil {
// 		t.Fatalf("We were not expecting error initializing cassandra but got %v", err_init)
// 	}
// 	ks, err := c.ShowColumnFamily("TEST_KEYSPACE", "users")
// 	if err != nil {
// 		t.Fatalf("We were not expecting error listing keyspaces but got %v", err)
// 	}
// 	t.Fatalf("%v", ks)
// }
