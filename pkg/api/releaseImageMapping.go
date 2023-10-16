package api

import (
	DB "carbide-api/pkg/database"
	"carbide-api/pkg/objects"
	"database/sql"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Responds with a JSON array of all images in the database
//
// Success Code: 200 OK
func releaseImageMappingGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	allReleaseImageMappings, err := DB.GetAllReleaseImgMappings(db)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	allReleaseImageMappingsJSON, err := json.Marshal(allReleaseImageMappings)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(allReleaseImageMappingsJSON)
	if err != nil {
		log.Error(err)
	}
	return
}

// Accepts a JSON payload of a new image and responds with the new JSON object after it's been successfully created in the database
//
// Success Code: 201 Created
func releaseImageMappingPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var receivedReleaseImageMapping objects.ReleaseImageMapping
	err := json.NewDecoder(r.Body).Decode(&receivedReleaseImageMapping)
	if err != nil || receivedReleaseImageMapping.ReleaseId == nil || receivedReleaseImageMapping.ImageId == nil {
		httpJSONError(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}
	err = DB.AddReleaseImgMapping(db, receivedReleaseImageMapping)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	newReleaseImageMapping, err := DB.GetReleaseImageMapping(db, receivedReleaseImageMapping)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"releaseImageMappingId": newReleaseImageMapping.Id,
		"releaseId":             *newReleaseImageMapping.ReleaseId,
		"imageId":               *newReleaseImageMapping.ImageId,
	}).Info("Release_image_mapping has been successfully created")
	newReleaseImageMappingJSON, err := json.Marshal(newReleaseImageMapping)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(newReleaseImageMappingJSON)
	if err != nil {
		log.Error(err)
	}
	return
}

// Deletes the releaseImageMapping and responds with an empty payload
//
// Success Code: 204 No Content
func releaseImageMappingDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var receivedReleaseImageMapping objects.ReleaseImageMapping
	err := json.NewDecoder(r.Body).Decode(&receivedReleaseImageMapping)
	if err != nil || receivedReleaseImageMapping.ReleaseId == nil || receivedReleaseImageMapping.ImageId == nil {
		httpJSONError(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}
	err = DB.DeleteReleaseImgMapping(db, receivedReleaseImageMapping)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"releaseImageMappingId": receivedReleaseImageMapping.Id,
		"releaseId":             *receivedReleaseImageMapping.ReleaseId,
		"imageId":               *receivedReleaseImageMapping.ImageId,
	}).Info("Release_image_mapping has been successfully deleted")
	w.WriteHeader(http.StatusNoContent)
	return
}

// Responds with the JSON representation of an releaseImageMapping
//
// Success Code: 200 OK
// func releaseImageMappingGet1(w http.ResponseWriter, r *http.Request, db *sql.DB, releaseImageMappingId int32) {
// 	var retrievedReleaseImageMapping objects.ReleaseImageMapping
// 	retrievedReleaseImageMapping, err := objects.GetReleaseImageMappingbyId(db, releaseImageMappingId)
// 	if err != nil {
// 		httpJSONError(w, err.Error(), http.StatusInternalServerError)
// 		log.Error(err)
// 		return
// 	}
// 	retrievedReleaseImageMappingJSON, err := json.Marshal(retrievedReleaseImageMapping)
// 	if err != nil {
// 		httpJSONError(w, err.Error(), http.StatusInternalServerError)
// 		log.Error(err)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	w.Header().Set("Content-Type", "application/json")
// 	_, err = w.Write(retrievedReleaseImageMappingJSON)
// 	if err != nil {
// 		log.Error(err)
// 	}
// 	return
// }

// // Accepts a JSON payload of the updated image and responds with the new JSON object after it's been successfully updated in the database
// //
// // Success Code: 200 OK
// func releaseImageMappingPut1(w http.ResponseWriter, r *http.Request, db *sql.DB, releaseImageMappingId int32) {
// 	var updatedImage objects.Image
// 	err := json.NewDecoder(r.Body).Decode(&updatedImage)
// 	if err != nil {
// 		httpJSONError(w, err.Error(), http.StatusBadRequest)
// 		log.Error(err)
// 		return
// 	}
// 	// image id cannot be overwritten with json payload
// 	updatedImage.Id = imageId
// 	err = objects.UpdateImage(db, updatedImage)
// 	if err != nil {
// 		httpJSONError(w, err.Error(), http.StatusInternalServerError)
// 		log.Error(err)
// 		return
// 	}
// 	updatedImage, err = objects.GetImagebyId(db, updatedImage.Id)
// 	if err != nil {
// 		log.Error(err)
// 		return
// 	}
// 	log.Info("Image %s has been successfully updated", *updatedImage.ImageName)
// 	updatedImageJSON, err := json.Marshal(updatedImage)
// 	if err != nil {
// 		httpJSONError(w, err.Error(), http.StatusInternalServerError)
// 		log.Error(err)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	w.Header().Set("Content-Type", "application/json")
// 	_, err = w.Write(updatedImageJSON)
// 	if err != nil {
// 		log.Error(err)
// 	}
// 	return
// }
