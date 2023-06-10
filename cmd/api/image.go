package api

import (
	"carbide-api/cmd/api/objects"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func serveImage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var image_name string
	image_name, r.URL.Path = ShiftPath(r.URL.Path)
	if image_name == "" {
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
		switch r.Method {
		case http.MethodGet:
			imageGet1(w, r, db, image_name)
			return
		case http.MethodPut:
			imagePut1(w, r, db, image_name)
			return
		case http.MethodDelete:
			imageDelete1(w, r, db, image_name)
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

	image, err = objects.GetImage(db, image.ImageName)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("Image %s has been successfully created", image.ImageName)

	return
}

func imageGet1(w http.ResponseWriter, r *http.Request, db *sql.DB, image_name string) {

	var image objects.Image
	image, err := objects.GetImage(db, image_name)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("Image %s has been successfully created", image.ImageName)

	return
}

func imagePut1(w http.ResponseWriter, r *http.Request, db *sql.DB, image string) {

	var updated_image objects.Image
	err := json.NewDecoder(r.Body).Decode(&updated_image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// ImageName field cannot be overwritten with json payload
	updated_image.ImageName = image
	err = objects.UpdateImage(db, updated_image)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("Image %s has been successfully updated", updated_image.ImageName)

	return
}

func imageDelete1(w http.ResponseWriter, r *http.Request, db *sql.DB, image_name string) {

	err := objects.DeleteImage(db, image_name)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("Image %s has been successfully created", image_name)

	return
}
