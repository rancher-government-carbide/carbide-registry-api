package api

import (
	"carbide-api/cmd/api/objects"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func serveImage(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string, release_name string) {
	switch r.Method {
	case http.MethodPost:
		imagePost(w, r, db)
		return
	case http.MethodOptions:
		return
	default:
		http.Error(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
		return
	}
}

func imagePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var image objects.Image
	err := json.NewDecoder(r.Body).Decode(&image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	return
}
