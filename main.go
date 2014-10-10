package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/juliendsv/go-cassadmin/dao"
	"github.com/juliendsv/go-cassadmin/domain"
	"github.com/juliendsv/go-cassadmin/handler"
)

// http://www.alexedwards.net/blog/a-mux-showdown

func main() {

	// We init our cassandra
	domain.DefaultStore, _ = dao.NewCassandraStore()

	r := mux.NewRouter()

	// Web
	// r.HandleFunc("/", handler.Home)

	// Api
	r.HandleFunc("/api/keyspaces", handler.APIKeyspaces)
	// r.HandleFunc("/api/keyspace/{ks}/{cf}", handler.APIShowCf)
	// r.HandleFunc("/api/{ks}/{cf}", handler.APIShowCf)

	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}
