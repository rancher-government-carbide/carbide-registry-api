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
func serveUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case http.MethodPost:
		userPost(w, r, db)
		return
	case http.MethodDelete:
		userDelete(w, r, db)
		return
	case http.MethodOptions:
		return
	default:
		http.Error(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
		return
	}
}

func userPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var newuser objects.User

	err := json.NewDecoder(r.Body).Decode(&newuser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Print(err)
		return
	}

	// add user to database
	err = objects.AddUser(db, newuser)
	if err != nil {
		log.Print(err)
		return
	}

	// reassign stored values in user object
	newuser, err = objects.GetUser(db, newuser.Username)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("User %s has been successfully created", newuser.Username)

	// respond with jwt
	token, err := generateJWT(newuser)
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Print(token, "\n")

	// provide token as cookie to frontend
	ck := http.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &ck)
	respondSuccess(w)

	return
}

func userDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// parse username from jwt token
	userid, err := verifyJWT(r)
	if err != nil {
		log.Print(err)
		http.Error(w, "Missing or expired cookie", 401)
		return
	}
	var user objects.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Print(err)
		return
	}
	// check credentials
	if err := objects.VerifyUser(db, user); err != nil {
		log.Print(err)
		return
	}
	// remove row with corresponding userid
	if err := objects.DeleteUser(db, userid); err != nil {
		log.Print(err)
		log.Printf("Failed to delete user with id %d", userid)
		w.Write([]byte(fmt.Sprintf("Failed to delete user with id %d", userid)))
		return
	}
	respondSuccess(w)

	return
}
