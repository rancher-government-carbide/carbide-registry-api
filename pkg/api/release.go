package api

import (
	DB "carbide-images-api/pkg/database"
	"carbide-images-api/pkg/objects"
	"database/sql"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Responds with a JSON array of all releases in the database
//
// Success Code: 200 OK
func releaseGet(w http.ResponseWriter, r *http.Request, db *sql.DB, productName string) {
	var allReleases []objects.Release
	allReleases, err := DB.GetAllReleasesforProduct(db, productName)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	allReleasesJSON, err := json.Marshal(allReleases)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(allReleasesJSON)
	if err != nil {
		log.Error(err)
	}
	return
}

// Accepts a JSON payload of a new release and responds with the new JSON object after it's been successfully created in the database
//
// Success Code: 201 Created
func releasePost(w http.ResponseWriter, r *http.Request, db *sql.DB, productName string) {
	var newRelease objects.Release
	err := json.NewDecoder(r.Body).Decode(&newRelease)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}
	parentProduct, err := DB.GetProduct(db, productName)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	newRelease.ProductId = &parentProduct.Id
	err = DB.AddRelease(db, newRelease)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	createdRelease, err := DB.GetRelease(db, newRelease)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"release": *createdRelease.Name,
	}).Info("New release has been successfully created")
	createdReleaseJSON, err := json.Marshal(createdRelease)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(createdReleaseJSON)
	if err != nil {
		log.Error(err)
	}
	return
}

// Responds with the JSON representation of a release
//
// Success Code: 200 OK
func releaseGetByName(w http.ResponseWriter, r *http.Request, db *sql.DB, productName string, releaseName string) {
	var retrievedRelease objects.Release
	retrievedRelease.Name = &releaseName
	parentProduct, err := DB.GetProduct(db, productName)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	retrievedRelease.ProductId = &parentProduct.Id
	retrievedRelease, err = DB.GetRelease(db, retrievedRelease)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	retrievedReleaseJSON, err := json.Marshal(retrievedRelease)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(retrievedReleaseJSON)
	if err != nil {
		log.Error(err)
	}
	return
}

// Accepts a JSON payload of the updated release and responds with the new JSON object after it's been successfully updated in the database
//
// Success Code: 200 OK
func releasePutByName(w http.ResponseWriter, r *http.Request, db *sql.DB, productName string, releaseName string) {
	var receivedRelease objects.Release
	err := json.NewDecoder(r.Body).Decode(&receivedRelease)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}
	receivedRelease.Name = &releaseName
	parentProduct, err := DB.GetProduct(db, productName)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	receivedRelease.ProductId = &parentProduct.Id
	err = DB.UpdateRelease(db, receivedRelease)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	updatedRelease, err := DB.GetRelease(db, receivedRelease)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"release": *updatedRelease.Name,
	}).Info("Release has been successfully updated")
	updatedReleaseJSON, err := json.Marshal(updatedRelease)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(updatedReleaseJSON)
	if err != nil {
		log.Error(err)
	}
	return
}

// Deletes the release and responds with an empty payload
//
// Success Code: 204 No Content
func releaseDeleteByName(w http.ResponseWriter, r *http.Request, db *sql.DB, productName string, releaseName string) {
	var releaseToDelete objects.Release
	releaseToDelete.Name = &releaseName
	parentProduct, err := DB.GetProduct(db, productName)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	releaseToDelete.ProductId = &parentProduct.Id
	err = DB.DeleteRelease(db, releaseToDelete)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"release": releaseName,
	}).Info("Release has been successfully deleted")
	w.WriteHeader(http.StatusNoContent)
	return
}
