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

func releaseGet(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string) {

	var releases []objects.Release
	releases, err := objects.GetAllReleasesforProduct(db, product_name)
	if err != nil {
		log.Print(err)
		return
	}
	releases_json, err := json.Marshal(releases)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(releases_json)
	if err != nil {
		log.Print(err)
	}
	return

}

func releasePost(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string) {

	var release objects.Release
	err := json.NewDecoder(r.Body).Decode(&release)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = objects.AddRelease(db, release)
	if err != nil {
		log.Print(err)
		return
	}
	release, err = objects.GetRelease(db, release)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("New release %s has been successfully created", *release.Name)

	return
}

func releaseGet1(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string, release_name string) {

	var release objects.Release
	*release.Name = release_name
	release, err := objects.GetRelease(db, release)
	if err != nil {
		log.Print(err)
		return
	}
	release_json, err := json.Marshal(release)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(release_json)
	if err != nil {
		log.Print(err)
	}
	return
}

func releasePut1(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string, release_name string) {

	var release objects.Release
	err := json.NewDecoder(r.Body).Decode(&release)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	*release.Name = release_name
	err = objects.UpdateRelease(db, release)
	if err != nil {
		log.Print(err)
		return
	}
	release, err = objects.GetRelease(db, release)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("Release %s has been successfully updated", *release.Name)

	return
}

func releaseDelete1(w http.ResponseWriter, r *http.Request, db *sql.DB, product_name string, release_name string) {

	var release_to_delete objects.Release
	*release_to_delete.Name = release_name
	// get product id and assign to release.ProductId
	err := objects.DeleteRelease(db, release_to_delete)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("Release %s has been successfully deleted", release_name)
	return
}
