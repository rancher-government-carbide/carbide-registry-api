package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// shiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

type Serve struct {
	DB *sql.DB
}

func (h Serve) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// return 404 Not Found for any URL with a trailing slash (except "/" itself).
	if r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/") {
		http.NotFound(w, r)
		return
	}
	GlobalMiddleware(w, r)
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	switch head {
	case "product":
		serveProduct(w, r, h.DB)
		return
	case "image":
		serveImage(w, r, h.DB)
		return
	case "releaseImageMapping":
		serveReleaseImageMapping(w, r, h.DB)
		return
	case "user":
		serveUser(w, r, h.DB)
		return
	case "login":
		serveLogin(w, r, h.DB)
		return
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func serveUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case http.MethodPost:
		userPost(w, r, db)
		return
	case http.MethodDelete:
		userDelete(w, r, db)
		return
	case http.MethodOptions:
		return
	default:
		httpJSONError(w, fmt.Sprintf("Expected method POST, DELETE, or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
		return
	}
}

func serveLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if userIsAuthenticated(w, r) {
		log.Info("user is already logged in\n")
		w.Write([]byte("user is already logged in"))
		return
	}
	switch r.Method {
	case http.MethodPost:
		loginPost(w, r, db)
		return
	case http.MethodOptions:
		return
	default:
		httpJSONError(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
		return
	}
}

func serveProduct(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if !userIsAuthenticated(w, r) {
		log.Info("user is unauthorized\n")
		w.Write([]byte("user is unauthorized"))
		return
	}
	var productName string
	productName, r.URL.Path = shiftPath(r.URL.Path)
	if r.URL.Path != "/" {
		var head string
		head, r.URL.Path = shiftPath(r.URL.Path)
		switch head {
		case "release":
			serveRelease(w, r, db, productName)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
		return
	}
	if productName == "" {
		switch r.Method {
		case http.MethodGet:
			productGet(w, r, db)
			return
		case http.MethodPost:
			productPost(w, r, db)
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
			productGetByName(w, r, db, productName)
			return
		case http.MethodPut:
			productPutByName(w, r, db, productName)
			return
		case http.MethodDelete:
			productDeleteByName(w, r, db, productName)
			return
		case http.MethodOptions:
			return
		default:
			http.Error(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}

func serveRelease(w http.ResponseWriter, r *http.Request, db *sql.DB, productName string) {
	if !userIsAuthenticated(w, r) {
		log.Info("user is unauthorized\n")
		w.Write([]byte("user is unauthorized"))
		return
	}
	var release_name string
	release_name, r.URL.Path = shiftPath(r.URL.Path)
	if r.URL.Path != "/" {
		var head string
		head, r.URL.Path = shiftPath(r.URL.Path)
		switch head {
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
		return
	}
	if release_name == "" {
		switch r.Method {
		case http.MethodGet:
			releaseGet(w, r, db, productName)
			return
		case http.MethodPost:
			releasePost(w, r, db, productName)
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
			releaseGetByName(w, r, db, productName, release_name)
			return
		case http.MethodPut:
			releasePutByName(w, r, db, productName, release_name)
			return
		case http.MethodDelete:
			releaseDeleteByName(w, r, db, productName, release_name)
			return
		case http.MethodOptions:
			return
		default:
			http.Error(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}

func serveImage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if !userIsAuthenticated(w, r) {
		log.Info("user is unauthorized\n")
		w.Write([]byte("user is unauthorized"))
		return
	}
	var image_id_string string
	image_id_string, r.URL.Path = shiftPath(r.URL.Path)
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
			httpJSONError(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
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
			httpJSONError(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}

func serveReleaseImageMapping(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if !userIsAuthenticated(w, r) {
		log.Info("user is unauthorized\n")
		w.Write([]byte("user is unauthorized"))
		return
	}
	switch r.Method {
	case http.MethodGet:
		releaseImageMappingGet(w, r, db)
		return
	case http.MethodPost:
		releaseImageMappingPost(w, r, db)
		return
	case http.MethodDelete:
		releaseImageMappingDelete(w, r, db)
		return
	case http.MethodOptions:
		return
	default:
		httpJSONError(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
		return
	}
}
