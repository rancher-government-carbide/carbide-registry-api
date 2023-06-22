package api

import (
	"database/sql"
	"net/http"
	"path"
	"strings"
)

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func ShiftPath(p string) (head, tail string) {
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

	Middleware(w, r)
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	switch head {
	case "product":
		serveProduct(w, r, h.DB)
		return
	case "image":
		serveImage(w, r, h.DB)
		return
	case "release_image_mapping":
		serveReleaseImageMapping(w, r, h.DB)
		return
		// 	case "user":
		// 		serveUser(w, r, h.DB)
		// 		return
		// 	case "login":
		// 		serveLogin(w, r, h.DB)
		// 		return
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}
