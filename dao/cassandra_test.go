package dao

import (
	"testing"
)

func TestNewCassandraStore(t *testing.T) {
	c, err_init := NewCassandraStore()
	if err_init != nil {
		t.Fatalf("We were not expecting error initializing cassandra but got %v", err_init)
	}
	ks, err := c.ListKeyspaces()
	if err != nil {
		t.Fatalf("We were not expecting error listing keyspaces but got %v", err)
	}
	t.Fatalf("%v", ks)
}
