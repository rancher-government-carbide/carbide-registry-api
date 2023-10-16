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
func imageGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	images, err := DB.GetAllImages(db)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	images_json, err := json.Marshal(images)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(images_json)
	if err != nil {
		log.Error(err)
	}
	return
}

// Accepts a JSON payload of a new image and responds with the new JSON object after it's been successfully created in the database
//
// Success Code: 201 Created
func imagePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var new_image objects.Image
	err := json.NewDecoder(r.Body).Decode(&new_image)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}
	if new_image.ImageName == nil {
		httpJSONError(w, "missing image name", http.StatusBadRequest)
		log.Error(err)
		return
	}
	err = DB.AddImage(db, new_image)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	created_image, err := DB.GetImagebyName(db, *new_image.ImageName)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"image": *created_image.ImageName,
	}).Info("Image has been successfully created")
	created_image_json, err := json.Marshal(created_image)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(created_image_json)
	if err != nil {
		log.Error(err)
	}
	return
}

// Responds with the JSON representation of an image (includes associated releases)
//
// Success Code: 200 OK
func imageGet1(w http.ResponseWriter, r *http.Request, db *sql.DB, image_id int32) {
	var image objects.Image
	image, err := DB.GetImagebyId(db, image_id)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	image_json, err := json.Marshal(image)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(image_json)
	if err != nil {
		log.Error(err)
	}
	return
}

// Accepts a JSON payload of the updated image and responds with the new JSON object after it's been successfully updated in the database
//
// Success Code: 200 OK
func imagePut1(w http.ResponseWriter, r *http.Request, db *sql.DB, image_id int32) {
	var updated_image objects.Image
	err := json.NewDecoder(r.Body).Decode(&updated_image)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}
	// image id cannot be overwritten with json payload
	updated_image.Id = image_id
	err = DB.UpdateImage(db, updated_image)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	updated_image, err = DB.GetImagebyId(db, updated_image.Id)
	if err != nil {
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"image": *updated_image.ImageName,
	}).Info("Image has been successfully updated")
	updated_image_json, err := json.Marshal(updated_image)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(updated_image_json)
	if err != nil {
		log.Error(err)
	}
	return
}

// Deletes the image and responds with an empty payload
//
// Success Code: 204 No Content
func imageDelete1(w http.ResponseWriter, r *http.Request, db *sql.DB, image_id int32) {
	err := DB.DeleteImage(db, image_id)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"image": image_id,
	}).Info("Image has been successfully deleted")
	w.WriteHeader(http.StatusNoContent)
	return
}
