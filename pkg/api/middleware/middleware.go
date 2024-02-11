package middleware

import (
	"net/http"
)

func CORS(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		enableCors(w, r)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func JWTAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		checkAuth(w, r)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
