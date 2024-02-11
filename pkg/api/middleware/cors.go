package middleware

import "net/http"

func enableCors(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	(w).Header().Set("Access-Control-Allow-Origin", origin)
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Charset, Accept-Language, Access-Control-Request-Method, Access-Control-Request-Headers, Authorization, Content-Length, Content-Type, Cookie, Date, Forwarded, Origin, User-Agent")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	// TODO: Separate allowed methods (and maybe headers) by endpoint
	(w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
}
