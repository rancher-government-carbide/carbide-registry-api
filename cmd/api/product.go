package api

import (
	"carbide-api/cmd/api/objects"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func serveProduct(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var product_name string
	product_name, r.URL.Path = ShiftPath(r.URL.Path)
	if r.URL.Path != "/" {
		var head string
		head, r.URL.Path = ShiftPath(r.URL.Path)
		switch head {
		case "release":
			serveRelease(w, r, db, product_name)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
		return
	}
	switch r.Method {
	case http.MethodPost:
		productPost(w, r, db)
		return
	case http.MethodOptions:
		return
	default:
		http.Error(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
		return
	}
}

func productPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var product objects.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("New Product: %v", product)

	return
}
