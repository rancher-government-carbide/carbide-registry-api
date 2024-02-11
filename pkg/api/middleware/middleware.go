package middleware

import (
	"net/http"
)

func Global(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		enableCors(w, r)
		next(w, r)
	}
	return http.HandlerFunc(fn)
}

func JWTAuth(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		checkAuth(w, r)
		next(w, r)
	}
	return http.HandlerFunc(fn)
}
