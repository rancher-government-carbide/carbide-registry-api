package api

import (
	"carbide-api/cmd/api/objects"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
			http.Error(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	} else {

		image_id_64, err := strconv.ParseInt(image_id_string, 10, 32)
		if err != nil {
			log.Print(err)
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
			http.Error(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}

func imageGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	images, err := objects.GetAllImages(db)
	if err != nil {
		log.Print(err)
		return
	}
	images_json, err := json.Marshal(images)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(images_json)
	if err != nil {
		log.Print(err)
	}
	return

}

func imagePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var image objects.Image
	err := json.NewDecoder(r.Body).Decode(&image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = objects.AddImage(db, image)
	if err != nil {
		log.Print(err)
		return
	}

	image, err = objects.GetImage(db, image.Id)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("Image %s has been successfully created", *image.ImageName)

	return
}

func imageGet1(w http.ResponseWriter, r *http.Request, db *sql.DB, image_id int32) {

	var image objects.Image
	image, err := objects.GetImage(db, image_id)
	if err != nil {
		log.Print(err)
		return
	}
	image_json, err := json.Marshal(image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(image_json)
	if err != nil {
		log.Print(err)
	}
	return
}

func imagePut1(w http.ResponseWriter, r *http.Request, db *sql.DB, image_id int32) {

	var updated_image objects.Image
	err := json.NewDecoder(r.Body).Decode(&updated_image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// ImageName field cannot be overwritten with json payload
	updated_image.Id = image_id
	err = objects.UpdateImage(db, updated_image)
	if err != nil {
		log.Print(err)
		return
	}

	log.Printf("Image %d has been successfully updated", updated_image.Id)

	return
}

func imageDelete1(w http.ResponseWriter, r *http.Request, db *sql.DB, image_id int32) {

	err := objects.DeleteImage(db, image_id)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("Image %d has been successfully deleted", image_id)

	return
}
