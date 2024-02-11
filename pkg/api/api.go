package api

import (
	"carbide-images-api/pkg/api/utils"
	"database/sql"
	"fmt"
	"net/http"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

func InitRouter(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("POST /user", createUserHandler(db))
	mux.Handle("DELETE /user", deleteUserHandler(db))
	mux.Handle("POST /login", loginHandler(db))
	mux.Handle("GET /product", getAllProductsHandler(db))
	mux.Handle("POST /product", createProductHandler(db))
	mux.Handle("GET /product/{name}", getProductHandler(db))
	mux.Handle("PUT /product/{name}", updateProductHandler(db))
	mux.Handle("DELETE /product/{name}", deleteProductHandler(db))
	mux.Handle("GET /release", getAllReleasesHandler(db))
	mux.Handle("POST /release", createReleaseHandler(db))
	mux.Handle("GET /release/{name}", getReleaseHandler(db))
	mux.Handle("DELETE /release/{name}", deleteReleaseHandler(db))
	mux.Handle("GET /image", getAllImagesHandler(db))
	mux.Handle("POST /image", createImageHandler(db))
	mux.Handle("GET /image/{id}", getImageHandler(db))
	mux.Handle("PUT /image/{id}", updateImageHandler(db))
	mux.Handle("DELETE /image/{id}", deleteImageHandler(db))
	mux.Handle("GET /releaseImageMapping", getAllReleaseImageMappingsHandler(db))
	mux.Handle("POST /releaseImageMapping", createReleaseImageMappingHandler(db))
	mux.Handle("DELETE /releaseImageMapping", deleteReleaseImageMappingHandler(db))

	return mux
}

func sampleHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
	}

	return http.HandlerFunc(fn)
}

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
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func serveUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case http.MethodOptions:
		return
	default:
		utils.HttpJSONError(w, fmt.Sprintf("Expected method POST, DELETE, or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
		return
	}
}

func serveLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case http.MethodOptions:
		return
	default:
		utils.HttpJSONError(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
		return
	}
}

func serveProduct(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if !userIsAuthenticated(w, r) {
		log.Info("user is unauthorized\n")
		utils.RespondWithJSON(w, "user is unauthorized")
		return
	}
	var productName string
	productName, r.URL.Path = shiftPath(r.URL.Path)
	if r.URL.Path != "/" {
		var head string
		head, r.URL.Path = shiftPath(r.URL.Path)
		switch head {
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
		return
	}
	if productName == "" {
		switch r.Method {
		case http.MethodOptions:
			return
		default:
			http.Error(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	} else {
		switch r.Method {
		case http.MethodOptions:
			return
		default:
			http.Error(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}

func serveRelease(w http.ResponseWriter, r *http.Request, db *sql.DB, productName string) {
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
		case http.MethodOptions:
			return
		default:
			http.Error(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	} else {
		switch r.Method {
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
		case http.MethodOptions:
			return
		default:
			utils.HttpJSONError(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	} else {
		switch r.Method {
		case http.MethodOptions:
			return
		default:
			utils.HttpJSONError(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}

func serveReleaseImageMapping(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if !userIsAuthenticated(w, r) {
		log.Info("user is unauthorized\n")
		utils.RespondWithJSON(w, "user is unauthorized")
		return
	}
	switch r.Method {
	case http.MethodOptions:
		return
	default:
		utils.HttpJSONError(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
		return
	}
}
