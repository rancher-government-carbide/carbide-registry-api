package api

import (
	DB "carbide-api/pkg/database"
	"carbide-api/pkg/objects"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func serveReleaseImageMapping(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case http.MethodGet:
		release_image_mappingGet(w, r, db)
		return
	case http.MethodPost:
		release_image_mappingPost(w, r, db)
		return
	case http.MethodDelete:
		release_image_mappingDelete(w, r, db)
		return
	case http.MethodOptions:
		return
	default:
		http_json_error(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
		return
	}
}

// Responds with a JSON array of all images in the database
//
// Success Code: 200 OK
func release_image_mappingGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	all_release_image_mappings, err := DB.GetAllReleaseImgMappings(db)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	all_release_image_mappings_json, err := json.Marshal(all_release_image_mappings)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(all_release_image_mappings_json)
	if err != nil {
		log.Error(err)
	}
	return
}

// Accepts a JSON payload of a new image and responds with the new JSON object after it's been successfully created in the database
//
// Success Code: 201 Created
func release_image_mappingPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var received_release_image_mapping objects.Release_Image_Mapping
	err := json.NewDecoder(r.Body).Decode(&received_release_image_mapping)
	if err != nil || received_release_image_mapping.ReleaseId == nil || received_release_image_mapping.ImageId == nil {
		http_json_error(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}
	err = DB.AddReleaseImgMapping(db, received_release_image_mapping)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	new_release_image_mapping, err := DB.GetReleaseImageMapping(db, received_release_image_mapping)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"release_image_mapping_id": new_release_image_mapping.Id,
		"release_id":               *new_release_image_mapping.ReleaseId,
		"image_id":                 *new_release_image_mapping.ImageId,
	}).Info("Release_image_mapping has been successfully created")
	new_release_image_mapping_json, err := json.Marshal(new_release_image_mapping)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(new_release_image_mapping_json)
	if err != nil {
		log.Error(err)
	}
	return
}

// Deletes the release_image_mapping and responds with an empty payload
//
// Success Code: 204 No Content
func release_image_mappingDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var received_release_image_mapping objects.Release_Image_Mapping
	err := json.NewDecoder(r.Body).Decode(&received_release_image_mapping)
	if err != nil || received_release_image_mapping.ReleaseId == nil || received_release_image_mapping.ImageId == nil {
		http_json_error(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}
	err = DB.DeleteReleaseImgMapping(db, received_release_image_mapping)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"release_image_mapping_id": received_release_image_mapping.Id,
		"release_id":               *received_release_image_mapping.ReleaseId,
		"image_id":                 *received_release_image_mapping.ImageId,
	}).Info("Release_image_mapping has been successfully deleted")
	w.WriteHeader(http.StatusNoContent)
	return
}

// Responds with the JSON representation of an release_image_mapping
//
// Success Code: 200 OK
// func release_image_mappingGet1(w http.ResponseWriter, r *http.Request, db *sql.DB, release_image_mapping_id int32) {
// 	var retrieved_release_image_mapping objects.Release_Image_Mapping
// 	retrieved_release_image_mapping, err := objects.GetReleaseImageMappingbyId(db, release_image_mapping_id)
// 	if err != nil {
// 		http_json_error(w, err.Error(), http.StatusInternalServerError)
// 		log.Error(err)
// 		return
// 	}
// 	retrieved_release_image_mapping_json, err := json.Marshal(retrieved_release_image_mapping)
// 	if err != nil {
// 		http_json_error(w, err.Error(), http.StatusInternalServerError)
// 		log.Error(err)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	w.Header().Set("Content-Type", "application/json")
// 	_, err = w.Write(retrieved_release_image_mapping_json)
// 	if err != nil {
// 		log.Error(err)
// 	}
// 	return
// }

// // Accepts a JSON payload of the updated image and responds with the new JSON object after it's been successfully updated in the database
// //
// // Success Code: 200 OK
// func release_image_mappingPut1(w http.ResponseWriter, r *http.Request, db *sql.DB, release_image_mapping_id int32) {
// 	var updated_image objects.Image
// 	err := json.NewDecoder(r.Body).Decode(&updated_image)
// 	if err != nil {
// 		http_json_error(w, err.Error(), http.StatusBadRequest)
// 		log.Error(err)
// 		return
// 	}
// 	// image id cannot be overwritten with json payload
// 	updated_image.Id = image_id
// 	err = objects.UpdateImage(db, updated_image)
// 	if err != nil {
// 		http_json_error(w, err.Error(), http.StatusInternalServerError)
// 		log.Error(err)
// 		return
// 	}
// 	updated_image, err = objects.GetImagebyId(db, updated_image.Id)
// 	if err != nil {
// 		log.Error(err)
// 		return
// 	}
// 	log.Info("Image %s has been successfully updated", *updated_image.ImageName)
// 	updated_image_json, err := json.Marshal(updated_image)
// 	if err != nil {
// 		http_json_error(w, err.Error(), http.StatusInternalServerError)
// 		log.Error(err)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	w.Header().Set("Content-Type", "application/json")
// 	_, err = w.Write(updated_image_json)
// 	if err != nil {
// 		log.Error(err)
// 	}
// 	return
// }
