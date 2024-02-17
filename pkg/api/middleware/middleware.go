package middleware

import (
	"net/http"
)

func ChainHandlers(handlers ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
    return func(finalHandler http.Handler) http.Handler {
		for _, handler := range handlers {
            finalHandler = handler(finalHandler)
        }
        return finalHandler
    }
}

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
