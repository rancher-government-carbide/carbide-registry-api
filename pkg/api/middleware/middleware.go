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
		if !Authorized(w, r) {
			// log.Info("user is unauthorized\n")
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
