package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/juliendsv/go-cassadmin/domain"
)

func APIKeyspaces(w http.ResponseWriter, r *http.Request) {
	keyspaces, err := domain.DefaultStore.ShowKeyspaces()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(keyspaces)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func APIShowCf(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_ = domain.DefaultStore.ShowColumnFamily(vars["ks"], vars["cf"])

	w.Header().Set("Content-Type", "application/json")
	// w.Write()
}
