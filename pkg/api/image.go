package api

import (
	DB "carbide-images-api/pkg/database"
	"carbide-images-api/pkg/objects"
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
	imagesJSON, err := json.Marshal(images)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(imagesJSON)
	if err != nil {
		log.Error(err)
	}
}

// Accepts a JSON payload of a new image and responds with the new JSON object after it's been successfully created in the database
//
// Success Code: 201 Created
func imagePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var newImage objects.Image
	err := json.NewDecoder(r.Body).Decode(&newImage)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
	}
	if newImage.ImageName == nil {
		httpJSONError(w, "missing image name", http.StatusBadRequest)
		log.Error(err)
		return
	}
	err = DB.AddImage(db, newImage)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	createdImage, err := DB.GetImagebyName(db, *newImage.ImageName)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"image": *createdImage.ImageName,
	}).Info("Image has been successfully created")
	createdImageJSON, err := json.Marshal(createdImage)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(createdImageJSON)
	if err != nil {
		log.Error(err)
	}
}

// Responds with the JSON representation of an image (includes associated releases)
//
// Success Code: 200 OK
func imageGet1(w http.ResponseWriter, r *http.Request, db *sql.DB, imageId int32) {
	var image objects.Image
	image, err := DB.GetImagebyId(db, imageId)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	imageJSON, err := json.Marshal(image)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(imageJSON)
	if err != nil {
		log.Error(err)
	}
}

// Accepts a JSON payload of the updated image and responds with the new JSON object after it's been successfully updated in the database
//
// Success Code: 200 OK
func imagePut1(w http.ResponseWriter, r *http.Request, db *sql.DB, imageId int32) {
	var updatedImage objects.Image
	err := json.NewDecoder(r.Body).Decode(&updatedImage)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}
	// image id cannot be overwritten with json payload
	updatedImage.Id = imageId
	err = DB.UpdateImage(db, updatedImage)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	updatedImage, err = DB.GetImagebyId(db, updatedImage.Id)
	if err != nil {
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"image": *updatedImage.ImageName,
	}).Info("Image has been successfully updated")
	updatedImageJSON, err := json.Marshal(updatedImage)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(updatedImageJSON)
	if err != nil {
		log.Error(err)
	}
}

// Deletes the image and responds with an empty payload
//
// Success Code: 204 No Content
func imageDelete1(w http.ResponseWriter, r *http.Request, db *sql.DB, imageId int32) {
	err := DB.DeleteImage(db, imageId)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	log.WithFields(log.Fields{
		"image": imageId,
	}).Info("Image has been successfully deleted")
	w.WriteHeader(http.StatusNoContent)
}
