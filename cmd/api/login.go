package api

import (
	"carbide-api/cmd/api/objects"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Accepts: POST
// POST requests expect user credentials - responds with a jwt or error
func serveLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	Login_middleware(w, r)
	switch r.Method {
	case http.MethodPost:
		loginPost(w, r, db)
		return
	case http.MethodOptions:
		return
	default:
		http.Error(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
		return
	}
}

func loginPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// read login payload into user object
	var login objects.User
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// checks user+pass against stored database values
	if err := objects.VerifyUser(db, login); err != nil {
		log.Print(err)
		return
	}
	// fill user object with necessary information (userid)
	login, err = objects.GetUser(db, login.Username)
	if err != nil {
		log.Print(err)
		return
	}
	// creates jwt for user
	token, err := generateJWT(login)
	if err != nil {
		log.Print(err)
		return
	}
	// provide token as cookie to frontend
	ck := http.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	// responds with message to store cookie in browser for future requests
	http.SetCookie(w, &ck)

	log.Printf("%s logged in successfully", login.Username)
	respondSuccess(w)

	return
}
