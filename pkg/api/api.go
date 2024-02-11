package api

import (
	"carbide-images-api/pkg/api/middleware"
	"carbide-images-api/pkg/api/utils"
	"database/sql"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/justinas/alice"
)

func InitRouter(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	globalMiddleware := alice.New(middleware.CORS)
	authMiddleware := alice.New(middleware.CORS, middleware.JWTAuth)
	mux.Handle("POST /user", globalMiddleware.Then(createUserHandler(db)))
	mux.Handle("DELETE /user", globalMiddleware.Then(deleteUserHandler(db)))
	mux.Handle("POST /login", authMiddleware.Then(loginHandler(db)))
	mux.Handle("GET /product", authMiddleware.Then(getAllProductsHandler(db)))
	mux.Handle("POST /product", authMiddleware.Then(createProductHandler(db)))
	mux.Handle("GET /product/{name}", authMiddleware.Then(getProductHandler(db)))
	mux.Handle("PUT /product/{name}", authMiddleware.Then(updateProductHandler(db)))
	mux.Handle("DELETE /product/{name}", authMiddleware.Then(deleteProductHandler(db)))
	mux.Handle("GET /release", authMiddleware.Then(getAllReleasesHandler(db)))
	mux.Handle("POST /release", authMiddleware.Then(createReleaseHandler(db)))
	mux.Handle("GET /release/{name}", authMiddleware.Then(getReleaseHandler(db)))
	mux.Handle("DELETE /release/{name}", authMiddleware.Then(deleteReleaseHandler(db)))
	mux.Handle("GET /image", authMiddleware.Then(getAllImagesHandler(db)))
	mux.Handle("POST /image", authMiddleware.Then(createImageHandler(db)))
	mux.Handle("GET /image/{id}", authMiddleware.Then(getImageHandler(db)))
	mux.Handle("PUT /image/{id}", authMiddleware.Then(updateImageHandler(db)))
	mux.Handle("DELETE /image/{id}", authMiddleware.Then(deleteImageHandler(db)))
	mux.Handle("GET /releaseImageMapping", authMiddleware.Then(getAllReleaseImageMappingsHandler(db)))
	mux.Handle("POST /releaseImageMapping", authMiddleware.Then(createReleaseImageMappingHandler(db)))
	mux.Handle("DELETE /releaseImageMapping", authMiddleware.Then(deleteReleaseImageMappingHandler(db)))

	return mux
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
	switch r.Method {
	case http.MethodOptions:
		return
	default:
		utils.HttpJSONError(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
		return
	}
}
