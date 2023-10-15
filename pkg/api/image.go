package api

import (
	DB "carbide-api/pkg/database"
	"carbide-api/pkg/objects"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func serveImage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var image_id_string string
	image_id_string, r.URL.Path = ShiftPath(r.URL.Path)
	if image_id_string == "" {
		switch r.Method {
		case http.MethodGet:
			imageGet(w, r, db)
			return
		case http.MethodPost:
			imagePost(w, r, db)
			return
		case http.MethodOptions:
			return
		default:
			http_json_error(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	} else {

		image_id_64, err := strconv.ParseInt(image_id_string, 10, 32)
		if err != nil {
			log.Error(err)
			return
		}
		image_id := int32(image_id_64)

		switch r.Method {
		case http.MethodGet:
			imageGet1(w, r, db, image_id)
			return
		case http.MethodPut:
			imagePut1(w, r, db, image_id)
			return
		case http.MethodDelete:
			imageDelete1(w, r, db, image_id)
			return
		case http.MethodOptions:
			return
		default:
			http_json_error(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}

// Responds with a JSON array of all images in the database
//
// Success Code: 200 OK
func imageGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	images, err := DB.GetAllImages(db)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	images_json, err := json.Marshal(images)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
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
		http_json_error(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}
	if new_image.ImageName == nil {
		http_json_error(w, "missing image name", http.StatusBadRequest)
		log.Error(err)
		return
	}
	err = DB.AddImage(db, new_image)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	created_image, err := DB.GetImagebyName(db, *new_image.ImageName)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	log.Printf("Image %s has been successfully created", *created_image.ImageName)
	created_image_json, err := json.Marshal(created_image)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
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
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	image_json, err := json.Marshal(image)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
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
		http_json_error(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}
	// image id cannot be overwritten with json payload
	updated_image.Id = image_id
	err = DB.UpdateImage(db, updated_image)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	updated_image, err = DB.GetImagebyId(db, updated_image.Id)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Image %s has been successfully updated", *updated_image.ImageName)
	updated_image_json, err := json.Marshal(updated_image)
	if err != nil {
		http_json_error(w, err.Error(), http.StatusInternalServerError)
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
		http_json_error(w, err.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}
	log.Info("Image %d has been successfully deleted", image_id)
	w.WriteHeader(http.StatusNoContent)
	return
}
