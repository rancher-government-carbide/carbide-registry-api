package api

import (
	DB "carbide-api/pkg/database"
	"carbide-api/pkg/objects"
	"database/sql"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Responds with a JSON array of all releases in the database
//
// Success Code: 200 OK
func releaseGet(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string) {
	var all_releases []objects.Release
	all_releases, err := DB.GetAllReleasesforProduct(db, product_name)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	all_releases_json, err := json.Marshal(all_releases)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(all_releases_json)
	if err != nil {
		log.Error(err)
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
		httpJSONError(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}
	parent_product, err := DB.GetProduct(db, product_name)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	new_release.ProductId = &parent_product.Id
	err = DB.AddRelease(db, new_release)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	created_release, err := DB.GetRelease(db, new_release)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"release": *created_release.Name,
	}).Info("New release has been successfully created")
	created_release_json, err := json.Marshal(created_release)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(created_release_json)
	if err != nil {
		log.Error(err)
	}
	return
}

// Responds with the JSON representation of a release
//
// Success Code: 200 OK
func releaseGetByName(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string, release_name string) {
	var retrieved_release objects.Release
	retrieved_release.Name = &release_name
	parent_product, err := DB.GetProduct(db, product_name)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	retrieved_release.ProductId = &parent_product.Id
	retrieved_release, err = DB.GetRelease(db, retrieved_release)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	retrieved_release_json, err := json.Marshal(retrieved_release)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(retrieved_release_json)
	if err != nil {
		log.Error(err)
	}
	return
}

// Accepts a JSON payload of the updated release and responds with the new JSON object after it's been successfully updated in the database
//
// Success Code: 200 OK
func releasePutByName(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string, release_name string) {
	var recieved_release objects.Release
	err := json.NewDecoder(r.Body).Decode(&recieved_release)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}
	recieved_release.Name = &release_name
	parent_product, err := DB.GetProduct(db, product_name)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	recieved_release.ProductId = &parent_product.Id
	err = DB.UpdateRelease(db, recieved_release)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	updated_release, err := DB.GetRelease(db, recieved_release)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"release": *updated_release.Name,
	}).Info("Release has been successfully updated")
	updated_release_json, err := json.Marshal(updated_release)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(updated_release_json)
	if err != nil {
		log.Error(err)
	}
	return
}

// Deletes the release and responds with an empty payload
//
// Success Code: 204 No Content
func releaseDeleteByName(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string, release_name string) {
	var release_to_delete objects.Release
	release_to_delete.Name = &release_name
	parent_product, err := DB.GetProduct(db, product_name)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	release_to_delete.ProductId = &parent_product.Id
	err = DB.DeleteRelease(db, release_to_delete)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"release": release_name,
	}).Info("Release has been successfully deleted")
	w.WriteHeader(http.StatusNoContent)
	return
}
