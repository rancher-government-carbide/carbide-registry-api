package api

import (
	"carbide-api/cmd/api/objects"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func serveRelease(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string) {
	var release_name string
	release_name, r.URL.Path = ShiftPath(r.URL.Path)
	if r.URL.Path != "/" {
		var head string
		head, r.URL.Path = ShiftPath(r.URL.Path)
		switch head {
		case "image":
			serveImage(w, r, db, product_name, release_name)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
		return
	}
	switch r.Method {
	case http.MethodPost:
		releasePost(w, r, db)
		return
	case http.MethodOptions:
		return
	default:
		http.Error(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
		return
	}
}

func releasePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var release objects.Release
	err := json.NewDecoder(r.Body).Decode(&release)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	return
}
