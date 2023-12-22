package api

import (
	DB "carbide-images-api/pkg/database"
	"carbide-images-api/pkg/objects"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func userPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var newUser objects.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}
	// add user to database
	if err := DB.AddUser(db, newUser); err != nil {
		log.Print(err)
		return
	}
	// reassign stored values in user object
	newUser, err = DB.GetUser(db, newUser.Username)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("User %s has been successfully created", newUser.Username)
	// respond with jwt
	token, err := generateJWT(newUser)
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
	var userToDelete objects.User
	err = json.NewDecoder(r.Body).Decode(&userToDelete)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}
	// check credentials
	if err := DB.VerifyUser(db, userToDelete); err != nil {
		log.Print(err)
		return
	}
	// remove row with corresponding userid
	if err := DB.DeleteUser(db, userid); err != nil {
		log.Print(err)
		log.Printf("Failed to delete user with id %d", userid)
		w.Write([]byte(fmt.Sprintf("Failed to delete user with id %d", userid)))
		return
	}
	respondSuccess(w)
	return
}

func loginPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var login objects.User
	// payload should contain username and password of user (other fields are ignored)
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		httpJSONError(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}
	// checks user+pass against stored database values
	if err := DB.VerifyUser(db, login); err != nil {
		log.Print(err)
		return
	}
	// fill user object with necessary information (userid)
	login, err = DB.GetUser(db, login.Username)
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
}

// accepts POST requests with new user payloads - responds with jwt or error
// func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodPost:
// 		userPost()
// 	case http.MethodDelete:
// 		userDelete()
// 	case http.MethodOptions:
// 		return
// 	default:
// 		http.Error(w, fmt.Sprintf("Expected method POST, OPTIONS or DELETE, got %v", r.Method), http.StatusMethodNotAllowed)
// 		fmt.Printf("Expected method POST, OPTIONS or DELETE, got %v", r.Method)
// 		return
// 	}
//
// }
//
// // accepts POST requests with user credentials - responds with a jwt or error
// func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	_, err := verifyJWT(r)
// 	if err == nil {
// 		log.Print("User is already logged in\n")
// 		w.Write([]byte(fmt.Sprintf("User is already logged in")))
// 		return
// 	}
// 	switch r.Method {
// 	case http.MethodPost:
// 		return
// 	case http.MethodOptions:
// 		return
// 	default:
// 		http.Error(w, fmt.Sprintf("Expected method POST or OPTIONS, got %v", r.Method), http.StatusMethodNotAllowed)
// 		return
// 	}
// }
