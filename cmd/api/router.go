package api

import (
	"fmt"
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

type Api struct {
	LoginHandler *LoginHandler
	UserHandler  *UserHandler
}

func (h *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	middleware(w, r)
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	switch head {
	case "user":
		h.UserHandler.ServeHTTP(w, r)
		return
	case "login":
		h.LoginHandler.ServeHTTP(w, r)
		return
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

// Accepts: POST, DELETE
// POST accepts new user payloads - responds with jwt or error
// DELETE removes the account of the authenticated user
func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		userPost(w, r, h.DB)
		return
	case http.MethodDelete:
		userDelete(w, r, h.DB)
		return
	case http.MethodOptions:
		return
	default:
		http.Error(w, fmt.Sprintf("Expected method POST, OPTIONS or DELETE, got %v", r.Method), http.StatusMethodNotAllowed)
		fmt.Printf("Expected method POST, OPTIONS or DELETE, got %v", r.Method)
		return
	}

}

// Accepts: POST
// POST requests expect user credentials - responds with a jwt or error
func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	login_middleware(w, r)
	switch r.Method {
	case http.MethodPost:
		loginPost(w, r, h.DB)
		return
	case http.MethodOptions:
		return
	default:
		http.Error(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
		return
	}
}
