package api

import (
	"carbide-api/cmd/api/objects"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func serveRelease(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string) {
	var release_name string
	release_name, r.URL.Path = ShiftPath(r.URL.Path)
	if r.URL.Path != "/" {
		var head string
		head, r.URL.Path = ShiftPath(r.URL.Path)
		switch head {
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
		return
	}
	if release_name == "" {
		switch r.Method {
		case http.MethodGet:
			releaseGet(w, r, db, product_name)
			return
		case http.MethodPost:
			releasePost(w, r, db, product_name)
			return
		case http.MethodOptions:
			return
		default:
			http.Error(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	} else {
		switch r.Method {
		case http.MethodGet:
			releaseGet1(w, r, db, product_name, release_name)
			return
		case http.MethodPut:
			releasePut1(w, r, db, product_name, release_name)
			return
		case http.MethodDelete:
			releaseDelete1(w, r, db, product_name, release_name)
			return
		case http.MethodOptions:
			return
		default:
			http.Error(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}

// Responds with a JSON array of all releases in the database
//
// Success Code: 200 OK
func releaseGet(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string) {
	var all_releases []objects.Release
	all_releases, err := objects.GetAllReleasesforProduct(db, product_name)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	all_releases_json, err := json.Marshal(all_releases)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(all_releases_json)
	if err != nil {
		log.Print(err)
	}
	return
}

// Accepts a JSON payload of a new release and responds with the new JSON object after it's been successfully created in the database
//
// Success Code: 201 Created
func releasePost(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string) {
	var new_release objects.Release
	err := json.NewDecoder(r.Body).Decode(&new_release)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusBadRequest)
		log.Print(err)
		return
	}
	parent_product, err := objects.GetProduct(db, product_name)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	new_release.ProductId = &parent_product.Id
	err = objects.AddRelease(db, new_release)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	created_release, err := objects.GetRelease(db, new_release)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	log.Printf("New release %s has been successfully created", *created_release.Name)
	created_release_json, err := json.Marshal(created_release)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(created_release_json)
	if err != nil {
		log.Print(err)
	}
	return
}

// Responds with the JSON representation of a release
//
// Success Code: 200 OK
func releaseGet1(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string, release_name string) {
	var retrieved_release objects.Release
	retrieved_release.Name = &release_name
	parent_product, err := objects.GetProduct(db, product_name)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
	}
	retrieved_release.ProductId = &parent_product.Id
	retrieved_release, err = objects.GetRelease(db, retrieved_release)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	retrieved_release_json, err := json.Marshal(retrieved_release)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(retrieved_release_json)
	if err != nil {
		log.Print(err)
	}
	return
}

// Accepts a JSON payload of the updated release and responds with the new JSON object after it's been successfully updated in the database
//
// Success Code: 200 OK
func releasePut1(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string, release_name string) {
	var recieved_release objects.Release
	err := json.NewDecoder(r.Body).Decode(&recieved_release)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusBadRequest)
		log.Print(err)
		return
	}
	recieved_release.Name = &release_name
	parent_product, err := objects.GetProduct(db, product_name)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
	}
	recieved_release.ProductId = &parent_product.Id
	err = objects.UpdateRelease(db, recieved_release)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	updated_release, err := objects.GetRelease(db, recieved_release)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	log.Printf("Release %s has been successfully updated", *updated_release.Name)
	updated_release_json, err := json.Marshal(updated_release)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(updated_release_json)
	if err != nil {
		log.Print(err)
	}
	return
}

// Deletes the release and responds with an empty payload
//
// Success Code: 204 No Content
func releaseDelete1(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string, release_name string) {
	var release_to_delete objects.Release
	release_to_delete.Name = &release_name
	parent_product, err := objects.GetProduct(db, product_name)
	if err != nil {
		log.Print(err)
	}
	release_to_delete.ProductId = &parent_product.Id
	err = objects.DeleteRelease(db, release_to_delete)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("Release %s has been successfully deleted", release_name)
	w.WriteHeader(http.StatusNoContent)
	return
}
