package handler

import (
	"fmt"
	"net/http"
)

func Ui(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Keyspace</h1>")
}
